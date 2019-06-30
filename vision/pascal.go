package vision

import (
	context "context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dldataset/reader"
	"github.com/rai-project/dldataset/reader/tfrecord"
	"github.com/rai-project/dldataset/vision/support/object_detection"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/image/types"
	protobuf "github.com/ubccr/terf/protobuf"
)

// PascalLabeledImage ...
type PascalLabeledImage struct {
	width     int64
	height    int64
	fileName  string
	sourceID  string
	sha256    string
	difficult []int64
	truncated []int64
	pose      []byte
	features  []*dlframework.Feature
	data      *types.RGBImage
}

// PascalValidationTFRecord ...
type PascalValidationTFRecord struct {
	base
	name           string
	baseURL        string
	recordFileName string
	md5sum         string
	labelMap       *object_detection.StringIntLabelMap
	recordReader   *reader.TFRecordReader
}

var (
	Pascal2007ValidationTFRecord *PascalValidationTFRecord
	Pascal2012ValidationTFRecord *PascalValidationTFRecord
)

func NewPascalLabeledImageFromRecord(rec *protobuf.Example) *PascalLabeledImage {
	height := tfrecord.FeatureInt64(rec, "image/height")
	width := tfrecord.FeatureInt64(rec, "image/width")
	fileName := tfrecord.FeatureString(rec, "image/filename")
	sourceID := tfrecord.FeatureString(rec, "image/source_id")
	sha256 := tfrecord.FeatureString(rec, "image/key/sha256")
	imgFormat := tfrecord.FeatureString(rec, "image/format")
	img, err := getImageRecord(tfrecord.FeatureBytes(rec, "image/encoded"), imgFormat)
	if err != nil {
		panic(err)
	}
	bboxXmin := tfrecord.FeatureFloat32Slice(rec, "image/object/bbox/xmin")
	bboxXmax := tfrecord.FeatureFloat32Slice(rec, "image/object/bbox/xmax")
	bboxYmin := tfrecord.FeatureFloat32Slice(rec, "image/object/bbox/ymin")
	bboxYmax := tfrecord.FeatureFloat32Slice(rec, "image/object/bbox/ymax")
	classText := tfrecord.FeatureStringSlice(rec, "image/object/class/text")
	classesLabels := tfrecord.FeatureInt64Slice(rec, "image/object/class/label")
	difficult := tfrecord.FeatureInt64Slice(rec, "image/object/difficult")
	truncated := tfrecord.FeatureInt64Slice(rec, "image/object/truncated")
	pose := tfrecord.FeatureBytes(rec, "image/object/view")

	numBBoxes := len(bboxXmax)
	features := make([]*dlframework.Feature, numBBoxes)
	for ii := 0; ii < numBBoxes; ii++ {
		features[ii] = feature.New(
			feature.BoundingBoxType(),
			feature.BoundingBoxXmin(bboxXmin[ii]),
			feature.BoundingBoxXmax(bboxXmax[ii]),
			feature.BoundingBoxYmin(bboxYmin[ii]),
			feature.BoundingBoxYmax(bboxYmax[ii]),
			feature.BoundingBoxIndex(int32(classesLabels[ii])),
			feature.BoundingBoxLabel(classText[ii]),
			feature.AppendMetadata("difficult", difficult[ii]),
			feature.AppendMetadata("truncated", truncated[ii]),
			feature.AppendMetadata("pose", pose[ii]),
		)
	}

	return &PascalLabeledImage{
		width:     width,
		height:    height,
		fileName:  fileName,
		sourceID:  sourceID,
		sha256:    sha256,
		difficult: difficult,
		truncated: truncated,
		pose:      pose,
		features:  features,
		data:      img,
	}
}

// Label ...
func (l *PascalLabeledImage) Label() string {
	return "<undefined>"
}

// Data ...
func (l *PascalLabeledImage) Data() (interface{}, error) {
	return l.data, nil
}

// Feature ...
func (d *PascalLabeledImage) Feature() *dlframework.Feature {
	return d.features[0]
}

// Features ...
func (d *PascalLabeledImage) Features() dlframework.Features {
	return d.features
}

// Name ...
func (d *PascalValidationTFRecord) Name() string {
	return d.name
}

// CanonicalName ...
func (d *PascalValidationTFRecord) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *PascalValidationTFRecord) workingDir() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	return filepath.Join(d.baseWorkingDir, category, name)
}

// Download ...
func (d *PascalValidationTFRecord) Download(ctx context.Context) error {
	workingDir := d.workingDir()
	fileName := d.recordFileName
	downloadedFileName := filepath.Join(workingDir, fileName)
	if com.IsFile(downloadedFileName) {
		return nil
	}
	downloadedFileName, _, err := downloadmanager.DownloadFile(
		urlJoin(d.baseURL, fileName),
		downloadedFileName,
		downloadmanager.Context(ctx),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to download %v", fileName)
	}
	return nil
}

// New ...
func (d *PascalValidationTFRecord) New(ctx context.Context) (dldataset.Dataset, error) {
	return nil, nil
}

// Get ...
func (d *PascalValidationTFRecord) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	return nil, errors.New("get is not implemented for " + d.CanonicalName())
}

// List ...
func (d *PascalValidationTFRecord) List(ctx context.Context) ([]string, error) {
	return nil, errors.New("list is not implemented for " + d.CanonicalName())
}

func (d *PascalValidationTFRecord) loadRecord(ctx context.Context) error {
	workingDir := d.workingDir()
	recordFileName := filepath.Join(workingDir, d.recordFileName)
	if !com.IsFile(recordFileName) {
		return errors.Errorf("unable to find the record file in %v make sure to download the dataset first", recordFileName)
	}

	recordIOReader, err := reader.NewTFRecordReader(recordFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to load record from %v", recordFileName)
	}
	d.recordReader = recordIOReader
	return nil
}

// Load ...
func (d *PascalValidationTFRecord) Load(ctx context.Context) error {
	return d.loadRecord(ctx)
}

// Next ...
func (d *PascalValidationTFRecord) Next(ctx context.Context) (dldataset.LabeledData, error) {
	rec, err := d.recordReader.NextRecord(ctx)
	if err != nil {
		return nil, err
	}

	return NewPascalLabeledImageFromRecord(rec), nil
}

// Close ...
func (d *PascalValidationTFRecord) Close() error {
	if d.recordReader != nil {
		d.recordReader.Close()
	}
	return nil
}

func init() {
	config.AfterInit(func() {

		const baseURLPrefix = "https://s3.amazonaws.com/store.carml.org/datasets"

		labelMap, err := object_detection.Get("pascal_label_map.pbtxt")
		if err != nil {
			panic(fmt.Sprintf("failed to get pascal_label_map.pbtxt due to %v", err))
		}

		baseWorkingDir := filepath.Join(dldataset.Config.WorkingDirectory, "dldataset")
		Pascal2007ValidationTFRecord = &PascalValidationTFRecord{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: baseWorkingDir,
			},
			name:           "Pascal2007",
			labelMap:       labelMap,
			baseURL:        baseURLPrefix + "/pascal2007",
			recordFileName: "validation.tfrecord",
			md5sum:         "e646ecf0bf838fa39d34e58d87c3e914",
		}

		Pascal2012ValidationTFRecord = &PascalValidationTFRecord{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: baseWorkingDir,
			},
			name:           "Pascal2012",
			labelMap:       labelMap,
			baseURL:        baseURLPrefix + "/pascal2012",
			recordFileName: "validation.tfrecord",
			md5sum:         "9a59d26492103b8635ba0c916d68535a",
		}

		dldataset.Register(Pascal2007ValidationTFRecord)
		dldataset.Register(Pascal2012ValidationTFRecord)
	})
}
