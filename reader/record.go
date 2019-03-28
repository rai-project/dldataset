package reader

import "github.com/rai-project/image/types"

type ImageRecord struct {
	ID         uint64
	LabelIndex float32
	Image      *types.RGBImage
}

type ImageSegmentationRecord struct {
	ID         uint64
	LabelIndex float32
	Image      *types.RGBImage
}