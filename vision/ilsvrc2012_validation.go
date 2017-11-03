package vision

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dldataset/reader"
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
	recordReader      *reader.RecordIOReader
	fileOffsetMapping map[string]recordIoOffset
}

type iLSVRC2012ValidationRecordIOLabeledData struct {
	*reader.Record
}

type recordIoOffset struct {
	start int
	end   int
}

func (d *iLSVRC2012ValidationRecordIOLabeledData) Label() string {
	return synset[int(d.LabelIndex)]
}

func (d *iLSVRC2012ValidationRecordIOLabeledData) Data() (interface{}, error) {
	return d.Image, nil
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

func (d *ILSVRC2012ValidationRecordIO) loadRecord(ctx context.Context) error {
	workingDir := d.workingDir()
	recordFileName := filepath.Join(workingDir, d.recordFileName)
	if !com.IsFile(recordFileName) {
		return errors.Errorf("unable to find the record file in %v make sure to download the dataset first", recordFileName)
	}

	recordIOReader, err := reader.NewRecordIOReader(recordFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to load record from %v", recordFileName)
	}
	d.recordReader = recordIOReader
	return nil
}

func (d *ILSVRC2012ValidationRecordIO) Load(ctx context.Context) error {
	return d.loadRecord(ctx)
}

func (d *ILSVRC2012ValidationRecordIO) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	return nil, errors.New("get is not implemented for " + d.CanonicalName())
}

func (d *ILSVRC2012ValidationRecordIO) Next(ctx context.Context) (dldataset.LabeledData, error) {
	rec, err := d.recordReader.Next(ctx)
	if err != nil {
		return nil, err
	}

	return &iLSVRC2012ValidationRecordIOLabeledData{
		Record: rec,
	}, nil
}

func (d *ILSVRC2012ValidationRecordIO) Close() error {
	if d.recordReader != nil {
		d.recordReader.Close()
	}
	return nil
}

func init() {
	config.AfterInit(func() {

		const fileListPath = "/vision/support/ilsvrc2012_validation_file_list.txt"

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
