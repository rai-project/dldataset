package vision

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PaddlePaddle/recordio"
	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/downloadmanager"
	"github.com/spf13/cast"
	context "golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

var (
	iLSVRC2012ValidationRecordIO    *ILSVRC2012ValidationRecordIO
	iLSVRC2012Validation224RecordIO *ILSVRC2012ValidationRecordIO
	iLSVRC2012Validation227RecordIO *ILSVRC2012ValidationRecordIO
)

// ILSVRC2012ValidationFolder ...
type ILSVRC2012ValidationRecordIO struct {
	base
	imageSize         int
	baseURL           string
	listFileName      string
	indexFileName     string
	recordFileName    string
	index             *recordio.Index
	recordScanner     *recordio.RangeScanner
	recordFile        *os.File
	fileOffsetMapping map[string]recordIoOffset
	data              map[string]ILSVRC2012ValidationLabeledImage
}

type recordIoOffset struct {
	start int
	end   int
}

func (d *ILSVRC2012ValidationRecordIO) New(ctx context.Context) (dldataset.Dataset, error) {
	return nil, nil
}
func (d *ILSVRC2012ValidationRecordIO) Name() string {
	if d.imageSize == 0 {
		return "ilsvrc2012_validation"
	}
	return "ilsvrc2012_validation_" + cast.ToString(d.imageSize)
}

func (d *ILSVRC2012ValidationRecordIO) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *ILSVRC2012ValidationRecordIO) workingDir() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	return filepath.Join(d.baseWorkingDir, category, name)
}

func (d *ILSVRC2012ValidationRecordIO) Download(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	files := []string{d.listFileName, d.indexFileName, d.recordFileName}
	workingDir := d.workingDir()
	for ii := range files {
		fileName := files[ii]
		grp.Go(func() error {
			downloadedFileName := filepath.Join(workingDir, fileName)
			if com.IsFile(downloadedFileName) {
				return nil
			}
			downloadedFileName, err := downloadmanager.DownloadFile(
				urlJoin(d.baseURL, fileName),
				downloadedFileName,
				downloadmanager.Context(ctx),
			)
			if err != nil {
				return errors.Wrapf(err, "failed to download %v", fileName)
			}
			return nil
		})
	}
	err := grp.Wait()
	if err != nil {
		return err
	}
	_, err = d.populate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func keysFileOffset(s map[string]recordIoOffset) []string {
	keys := make([]string, len(s))

	ii := 0
	for k := range s {
		keys[ii] = k
		ii++
	}
	return keys
}

func (d *ILSVRC2012ValidationRecordIO) populate(ctx context.Context) ([]string, error) {

	workingDir := d.workingDir()
	listFileName := filepath.Join(workingDir, d.listFileName)
	if !com.IsFile(listFileName) {
		return nil, errors.Errorf("unable to find the list file in %v make sure to download the dataset first", listFileName)
	}

	bts, err := ioutil.ReadFile(listFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %v", listFileName)
	}

	fileContent := strings.TrimSpace(string(bts))
	lines := strings.Split(fileContent, "\n")
	files := make([]string, len(lines))
	d.fileOffsetMapping = make(map[string]recordIoOffset)
	for ii, line := range lines {
		fields := strings.Fields(line)
		fileName := fields[len(fields)-1]
		d.fileOffsetMapping[fileName] = recordIoOffset{
			start: cast.ToInt(fields[0]),
			end:   cast.ToInt(fields[1]),
		}
		files[ii] = fileName
	}

	return files, nil
}

func (d *ILSVRC2012ValidationRecordIO) List(ctx context.Context) ([]string, error) {

	if len(d.fileOffsetMapping) != 0 {
		return d.populate(ctx)
	}

	return keysFileOffset(d.fileOffsetMapping), nil
}

func (d *ILSVRC2012ValidationRecordIO) loadIndex(ctx context.Context) error {
	workingDir := d.workingDir()
	indexFileName := filepath.Join(workingDir, d.indexFileName)
	if !com.IsFile(indexFileName) {
		return errors.Errorf("unable to find the index file in %v make sure to download the dataset first", indexFileName)
	}

	bts, err := ioutil.ReadFile(indexFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to read %v", indexFileName)
	}

	idx, err := recordio.LoadIndex(bytes.NewReader(bts))
	if err != nil {
		return errors.Wrapf(err, "failed to load index from %v", indexFileName)
	}
	d.index = idx
	return nil
}

func (d *ILSVRC2012ValidationRecordIO) loadRecord(ctx context.Context, offset recordIoOffset) error {
	workingDir := d.workingDir()
	recordFileName := filepath.Join(workingDir, d.recordFileName)
	if !com.IsFile(recordFileName) {
		return errors.Errorf("unable to find the record file in %v make sure to download the dataset first", recordFileName)
	}

	if d.recordFile == nil {
		f, err := os.Open(recordFileName)
		if err != nil {
			return errors.Wrapf(err, "failed to open %v", recordFileName)
		}
		d.recordFile = f
	}

	rng := recordio.NewRangeScanner(d.recordFile, d.index, offset.start, offset.end)
	if rng == nil {
		return errors.Errorf("failed to load record from %v", recordFileName)
	}
	d.recordScanner = rng
	return nil
}

func (d *ILSVRC2012ValidationRecordIO) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	fileOffsetMapping := d.fileOffsetMapping
	offset, ok := fileOffsetMapping[name]
	if !ok {
		return nil, errors.Errorf("file %v not found", name)
	}
	if d.index == nil {
		if err := d.loadIndex(ctx); err != nil {
			return nil, err
		}
	}
	if d.recordScanner == nil {
		if err := d.loadRecord(ctx, offset); err != nil {
			return nil, err
		}
	}
	d.recordScanner.Record()
	return nil, nil
}

func (d *ILSVRC2012ValidationRecordIO) Close() error {
	if d.recordFile != nil {
		d.recordFile.Close()
	}
	return nil
}

func init() {
	config.AfterInit(func() {

		const fileListPath = "/vision/support/ilsvrc2012_validation_file_list.txt"
		const baseURL = "http://store.carml.org.s3.amazonaws.com/datasets/imagenet1k-val-"

		iLSVRC2012ValidationRecordIO = &ILSVRC2012ValidationRecordIO{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			baseURL:        "https://s3.amazonaws.com/store.carml.org/datasets/ILSVRC2012_img_val_256",
			listFileName:   "imagenet1k-val.lst",
			indexFileName:  "imagenet1k-val.idx",
			recordFileName: "imagenet1k-val.rec",
		}

		iLSVRC2012Validation224RecordIO = &ILSVRC2012ValidationRecordIO{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			imageSize:      224,
			baseURL:        "https://s3.amazonaws.com/store.carml.org/datasets/ILSVRC2012_img_val_224",
			listFileName:   "imagenet1k-val.lst",
			indexFileName:  "imagenet1k-val.idx",
			recordFileName: "imagenet1k-val.rec",
		}

		iLSVRC2012Validation227RecordIO = &ILSVRC2012ValidationRecordIO{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			imageSize:      227,
			baseURL:        "https://s3.amazonaws.com/store.carml.org/datasets/ILSVRC2012_img_val_227",
			listFileName:   "imagenet1k-val.lst",
			indexFileName:  "imagenet1k-val.idx",
			recordFileName: "imagenet1k-val.rec",
		}
		dldataset.Register(iLSVRC2012ValidationRecordIO)
		dldataset.Register(iLSVRC2012Validation224RecordIO)
		dldataset.Register(iLSVRC2012Validation227RecordIO)
	})
}
