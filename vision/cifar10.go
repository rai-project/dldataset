package vision

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/image/types"
	"github.com/rai-project/utils"
	context "golang.org/x/net/context"
)

var cifar10 *CIFAR10

type CIFAR10 struct {
	base
	url                 string
	fileName            string
	extractedFolderName string
	md5sum              string
	trainFileNameList   map[string]string
	testFileNameList    map[string]string
	labelFileName       string
	labels              []string
	labelByteSize       int
	pixelByteSize       int
	imageDimensions     []int
	data                map[string]CIFAR10LabeledImage
	isDownloaded        bool
}

type CIFAR10LabeledImage struct {
	label string
	data  *types.RGBImage
}

func (l CIFAR10LabeledImage) Label() string {
	return l.label
}

func (l CIFAR10LabeledImage) Data() (interface{}, error) {
	return l.data, nil
}

func (*CIFAR10) Name() string {
	return "CIFAR10"
}

func (d *CIFAR10) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *CIFAR10) New(ctx context.Context) (dldataset.Dataset, error) {
	return cifar10, nil
}

func (d *CIFAR10) Download(ctx context.Context) error {
	if d.isDownloaded {
		return nil
	}
	workingDir := d.workingDir()
	downloadedFileName := filepath.Join(workingDir, d.fileName)
	downloadedFileName, err := downloadmanager.DownloadFile(d.url, downloadedFileName, downloadmanager.Context(ctx))
	if err != nil {
		return err
	}
	ok, err := utils.MD5Sum.CheckFile(downloadedFileName, d.md5sum)
	if err != nil {
		return errors.Wrapf(err, "unable to perform md5sum on %s", downloadedFileName)
	}
	if !ok {
		return errors.Wrapf(err, "the md5 sum for %s did not match expected %s", downloadedFileName, d.md5sum)
	}
	if err := downloadmanager.Unarchive(workingDir, downloadedFileName); err != nil {
		return err
	}
	if err := d.move(ctx); err != nil {
		return err
	}
	archiveOutputDir := filepath.Join(workingDir, d.extractedFolderName)
	defer os.RemoveAll(archiveOutputDir)

	return nil
}

func (d *CIFAR10) move(ctx context.Context) error {
	workingDir := d.workingDir()
	archiveOutputDir := filepath.Join(workingDir, d.extractedFolderName)
	filesHashes := map[string]string{}
	for fileName, md5 := range d.trainFileNameList {
		filesHashes[fileName] = md5
	}
	for fileName, md5 := range d.testFileNameList {
		filesHashes[fileName] = md5
	}
	for fileName, md5 := range filesHashes {
		filePath := filepath.Join(archiveOutputDir, fileName)
		if !com.IsFile(filePath) {
			return errors.Errorf("the file %s for %s was not found in the extracted directory", fileName, d.CanonicalName())
		}
		newPath := filepath.Join(workingDir, fileName)
		if err := os.Rename(filePath, newPath); err != nil {
			return errors.Wrapf(err, "cannot move the file %s to %s", filePath, newPath)
		}
		ok, err := utils.MD5Sum.CheckFile(newPath, md5)
		if err != nil {
			return err
		}
		if !ok {
			return errors.Wrapf(err, "the md5 sum for %s did not match expected %s", newPath, md5)
		}
	}
	labelFilePath := filepath.Join(archiveOutputDir, d.labelFileName)
	if !com.IsFile(labelFilePath) {
		return errors.Errorf("the file %s for %s was not found in the extracted directory", labelFilePath, d.CanonicalName())
	}
	newLabelFilePath := filepath.Join(workingDir, d.labelFileName)
	if err := os.Rename(labelFilePath, newLabelFilePath); err != nil {
		return errors.Wrapf(err, "cannot move the file %s to %s", labelFilePath, newLabelFilePath)
	}
	return nil
}

func (d *CIFAR10) List(ctx context.Context) ([]string, error) {
	if err := d.read(ctx); err != nil {
		return nil, err
	}
	keys := []string{}
	for key := range d.data {
		keys = append(keys, key)
	}
	return keys, nil
}

func (d *CIFAR10) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	if err := d.read(ctx); err != nil {
		return nil, err
	}
	data, ok := d.data[name]
	if !ok {
		return nil, errors.Errorf("unable to find %s in the %s dataset", name, d.CanonicalName())
	}
	return data, nil
}

func (d *CIFAR10) read(ctx context.Context) error {
	if err := d.readLabels(ctx); err != nil {
		return err
	}
	if err := d.readData(ctx); err != nil {
		return err
	}
	return nil
}

func (d *CIFAR10) readData(ctx context.Context) error {
	if len(d.data) != 0 {
		return nil
	}

	workingDir := d.workingDir()
	data := map[string]CIFAR10LabeledImage{}

	read := func(offset int, class, fileName string) (int, error) {
		idx := offset
		filePath := filepath.Join(workingDir, fileName)
		f, err := os.Open(filePath)
		if err != nil {
			return idx, errors.Wrapf(err, "failed to open %s while performing md5 checksum", filePath)
		}
		defer f.Close()

		for {
			entry, err := d.readEntry(ctx, f)
			if err == io.EOF {
				return 0, nil
			}
			if err != nil {
				return idx, errors.Wrapf(err, "failed reading entry for %s", filePath)
			}
			data[class+"/"+strconv.Itoa(idx)] = *entry
			idx++
		}
		return idx, nil
	}
	idx := 0
	for fileName := range d.trainFileNameList {
		newIdx, err := read(idx, "train", fileName)
		if err != nil {
			return err
		}
		idx = newIdx
	}
	idx = 0
	for fileName := range d.trainFileNameList {
		newIdx, err := read(idx, "test", fileName)
		if err != nil {
			return err
		}
		idx = newIdx
	}

	d.data = data

	return nil
}

func (d *CIFAR10) readEntry(ctx context.Context, reader io.Reader) (*CIFAR10LabeledImage, error) {
	var labelIdx int8
	labelByteSize := int64(d.labelByteSize)
	labelBytesReader := io.LimitReader(reader, labelByteSize)
	err := binary.Read(labelBytesReader, binary.LittleEndian, &labelIdx)
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("unable to read label")
	}
	if int(labelIdx) >= len(d.labels) {
		return nil, errors.Errorf("the label %v is out of range of %v", labelIdx, len(d.labels))
	}

	pixelByteSize := int64(d.pixelByteSize)
	pixelBytesReader := io.LimitReader(reader, pixelByteSize)

	img := types.NewRGBImage(image.Rect(0, 0, d.imageDimensions[0], d.imageDimensions[1]))

	err = binary.Read(pixelBytesReader, binary.LittleEndian, img.Pix)
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("unable to read label")
	}

	return &CIFAR10LabeledImage{
		label: d.labels[labelIdx],
		data:  img,
	}, nil
}

func (d *CIFAR10) readLabels(ctx context.Context) error {
	if len(d.labels) != 0 {
		return nil
	}

	workingDir := d.workingDir()
	labelFilePath := filepath.Join(workingDir, d.labelFileName)
	if !com.IsFile(labelFilePath) {
		return errors.Errorf("the label file %s was not found", labelFilePath)
	}

	var labels []string
	f, err := os.Open(labelFilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot read %s", labelFilePath)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}
	d.labels = labels
	return nil
}

func (d *CIFAR10) Close() error {
	return nil
}

func (d *CIFAR10) workingDir() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	return filepath.Join(d.baseWorkingDir, category, name)
}

func init() {
	config.AfterInit(func() {
		cifar10 = &CIFAR10{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			url:                 "https://www.cs.toronto.edu/~kriz/cifar-10-binary.tar.gz",
			fileName:            "cifar-10-binary.tar.gz",
			extractedFolderName: "cifar-10-batches-bin",
			md5sum:              "c32a1d4ab5d03f1284b67883e8d87530",
			trainFileNameList: map[string]string{
				"data_batch_1.bin": "5dd7e06a14cb22eb9f671a540d1b7c25",
				"data_batch_2.bin": "5ea93a67294ea407fff1d09f752e9692",
				"data_batch_3.bin": "942cd6a4bcdd0dd3c604fbe906cb4421",
				"data_batch_4.bin": "ae636b3ba5c66a11e91e8cb52e771fcb",
				"data_batch_5.bin": "53f37980c15c3d472c316c40844f3f0d",
			},
			testFileNameList: map[string]string{
				"test_batch.bin": "803d5f7f4d78ea53de84dbe85f74fb6d",
			},
			labelFileName:   "batches.meta.txt",
			imageDimensions: []int{32, 32, 3},
			labelByteSize:   1,
			pixelByteSize:   3072,
			isDownloaded:    false,
		}
		dldataset.Register(cifar10)
	})
}
