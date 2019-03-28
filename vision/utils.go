package vision

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
)

func urlJoin(base string, n string) string {
	if strings.HasPrefix(base, "/") {
		base = strings.TrimPrefix(base, "/")
	}
	return strings.Join([]string{base, n}, "/")
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
