package reader

import (
	"bytes"
	context "context"
	goimage "image"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/ubccr/terf"
	protobuf "github.com/ubccr/terf/protobuf"
)

type TFRecordReader struct {
	r io.ReadCloser
	*terf.Reader
}

// NewTFRecordReader ...
func NewTFRecordReader(path string) (*TFRecordReader, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open %v", path)
	}
	return &TFRecordReader{
		r:      r,
		Reader: terf.NewReader(r),
	}, nil
}

func (r *TFRecordReader) NextRecord(ctx context.Context) (*protobuf.Example, error) {
	nxt, err := r.Reader.Next()
	if err != nil {
		return nil, err
	}
	return nxt, nil
}

func (r *TFRecordReader) Next(ctx context.Context) (*ImageRecord, error) {
	nxt, err := r.Reader.Next()
	if err != nil {
		return nil, err
	}

	imgRecord := new(terf.Image)
	err = imgRecord.UnmarshalExample(nxt)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal image")
	}

	if strings.ToLower(imgRecord.Format) == "cifar" {
		img := types.NewRGBImage(goimage.Rect(0, 0, imgRecord.Width, imgRecord.Height))
		imgPix := img.Pix
		inputPix := imgRecord.Raw
		channels := 3

		for h := 0; h < imgRecord.Height; h++ {
			for w := 0; w < imgRecord.Width; w++ {
				for c := 0; c < channels; c++ {
					imgPix[channels*(h*imgRecord.Width+w)+c] =
						// inputPix format = The first 1024 entries contain the red channel values,
						// the next 1024 the green, and the final 1024 the blue.
						inputPix[c*imgRecord.Height*imgRecord.Width+h*imgRecord.Width+w]
				}
			}
		}
		return &ImageRecord{
			ID:         uint64(imgRecord.ID),
			LabelIndex: float32(imgRecord.LabelID),
			Image:      img,
		}, nil
	}

	img, err := image.Read(bytes.NewBuffer(imgRecord.Raw), image.Context(nil))
	if err != nil {
		return nil, err
	}

	rgbImage, ok := img.(*types.RGBImage)
	if !ok {
		return nil, errors.Errorf("expecting an rgb image")
	}

	return &ImageRecord{
		ID:         uint64(imgRecord.ID),
		LabelIndex: float32(imgRecord.LabelID),
		Image:      rgbImage,
	}, nil
}

func (r *TFRecordReader) Close() error {
	return r.r.Close()
}
