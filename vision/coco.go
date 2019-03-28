package vision

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/rai-project/dldataset/reader/tfrecord"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	protobuf "github.com/ubccr/terf/protobuf"
)

// CIFAR100LabeledImage ...
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

func getImageRecord(data []byte, format string) (*types.RGBImage, error) {
	img, err := image.Read(bytes.NewBuffer(data), image.Context(nil))
	if err != nil {
		return nil, err
	}

	rgbImage, ok := img.(*types.RGBImage)
	if !ok {
		return nil, errors.Errorf("expecting an rgb image")
	}

	return rgbImage, nil
}

func nextCocoFromRecord(rec *protobuf.Example) *CocoLabeledImage {
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
