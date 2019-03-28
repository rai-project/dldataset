package vision

import (
	context "context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/rai-project/dldataset/vision/support/object_detection"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/dldataset/reader"
	"github.com/rai-project/dldataset/reader/tfrecord"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/image/types"
	protobuf "github.com/ubccr/terf/protobuf"
)

// CocoLabeledImage ...
type CocoLabeledImage struct {
	width    int64
	height   int64
	fileName string
	sourceID string
	sha256   string
	area     []float32
	isCrowd  []int64
	features []*dlframework.Feature
	data     *types.RGBImage
}

// CocoValidationTFRecord ...
type CocoValidationTFRecord struct {
	base
	name             string
	baseURL          string
	recordFileName   string
	md5sum           string
	labelMap         object_detection.StringIntLabelMap
	completeLabelMap object_detection.StringIntLabelMap
	recordReader     *reader.TFRecordReader
}

var (
	coco2014ValidationTFRecord *CocoValidationTFRecord
	coco2017ValidationTFRecord *CocoValidationTFRecord
)

// Label ...
func (l *CocoLabeledImage) Label() string {
	return "<undefined>"
}

// Data ...
func (l *CocoLabeledImage) Data() (interface{}, error) {
	return l.data, nil
}

// Feature ...
func (d *CocoLabeledImage) Feature() *dlframework.Feature {
	return d.features[0]
}

// Features ...
func (d *CocoLabeledImage) Features() dlframework.Features {
	return d.features
}

// Close ...
func (d *CocoValidationTFRecord) Close() error {
	if d.recordReader != nil {
		d.recordReader.Close()
	}
	return nil
}

// Name ...
func (d *CocoValidationTFRecord) Name() string {
	return d.name
}

// CanonicalName ...
func (d *CocoValidationTFRecord) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *CocoValidationTFRecord) workingDir() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	return filepath.Join(d.baseWorkingDir, category, name)
}

// Download ...
func (d *CocoValidationTFRecord) Download(ctx context.Context) error {
	workingDir := d.workingDir()
	fileName := d.recordFileName
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
}

// New ...
func (d *CocoValidationTFRecord) New(ctx context.Context) (dldataset.Dataset, error) {
	return nil, nil
}

// Get ...
func (d *CocoValidationTFRecord) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	return nil, errors.New("get is not implemented for " + d.CanonicalName())
}

// List ...
func (d *CocoValidationTFRecord) List(ctx context.Context) ([]string, error) {
	return nil, errors.New("list is not implemented for " + d.CanonicalName())
}

func (d *CocoValidationTFRecord) loadRecord(ctx context.Context) error {
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
func (d *CocoValidationTFRecord) Load(ctx context.Context) error {
	return d.loadRecord(ctx)
}

// Next ...
func (d *CocoValidationTFRecord) Next(ctx context.Context) (dldataset.LabeledData, error) {
	rec, err := d.recordReader.NextRecord(ctx)
	if err != nil {
		return nil, err
	}

	return NewCocoLabeledImageFromRecord(rec), nil
}

func NewCocoLabeledImageFromRecord(rec *protobuf.Example) *CocoLabeledImage {
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
	class := tfrecord.FeatureStringSlice(rec, "image/object/class/text")
	isCrowd := tfrecord.FeatureInt64Slice(rec, "image/object/is_crowd")
	area := tfrecord.FeatureFloat32Slice(rec, "image/object/area")

	numBBoxes := len(bboxXmax)
	features := make([]*dlframework.Feature, numBBoxes)
	for ii := 0; ii < numBBoxes; ii++ {
		features[ii] = feature.New(
			feature.BoundingBoxType(),
			feature.BoundingBoxXmin(bboxXmin[ii]),
			feature.BoundingBoxXmax(bboxXmax[ii]),
			feature.BoundingBoxYmin(bboxYmin[ii]),
			feature.BoundingBoxYmax(bboxYmax[ii]),
			feature.BoundingBoxLabel(class[ii]),
			feature.AppendMetadata("isCrowd", isCrowd[ii]),
			feature.AppendMetadata("area", area[ii]),
		)
	}

	return &CocoLabeledImage{
		width:    width,
		height:   height,
		fileName: fileName,
		sourceID: sourceID,
		sha256:   sha256,
		area:     area,
		isCrowd:  isCrowd,
		features: features,
		data:     img,
	}
}

func init() {
	config.AfterInit(func() {

		const baseURLPrefix = "https://s3.amazonaws.com/store.carml.org/datasets"

		labelMap, err := object_detection.Get("mscoco_label_map.pbtxt")
		if err != nil {
			panic(fmt.Sprintf("failed to get mscoco_label_map.pbtxt due to %v", err))
		}

		completeLabelMap, err := object_detection.Get("mscoco_complete_label_map.pbtxt")
		if err != nil {
			panic(fmt.Sprintf("failed to get mscoco_complete_label_map.pbtxt due to %v", err))
		}

		baseWorkingDir := filepath.Join(dldataset.Config.WorkingDirectory, "dldataset")
		coco2014ValidationTFRecord = &CocoValidationTFRecord{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: baseWorkingDir,
			},
			name:             "coco2014",
			baseURL:          baseURLPrefix + "/coco2014",
			labelMap:         labelMap,
			completeLabelMap: completeLabelMap,
			recordFileName:   "coco_val.record-00000-of-00001",
			md5sum:           "b1f63512f72d3c84792a1f53ec40062a",
		}

		coco2017ValidationTFRecord = &CocoValidationTFRecord{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: baseWorkingDir,
			},
			name:             "coco2017",
			baseURL:          baseURLPrefix + "/coco2017",
			labelMap:         labelMap,
			completeLabelMap: completeLabelMap,
			recordFileName:   "coco_val.record-00000-of-00001",
			md5sum:           "b8a0cfed5ad569d4572b4ad8645acb5b",
		}

		dldataset.Register(coco2014ValidationTFRecord)
		dldataset.Register(coco2017ValidationTFRecord)
	})
}
