package reader

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"testing"
	"unsafe"

	context "golang.org/x/net/context"

	goimage "image"
	"image/png"

	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	_ "github.com/rai-project/tracer/noop"
	"github.com/stretchr/testify/assert"
)

func nextRecord(t *testing.T, f io.Reader) {

	ctx := context.Background()
	var err error
	var magic uint32
	err = binary.Read(f, binary.LittleEndian, &magic)
	assert.NoError(t, err)

	assert.Equal(t, kMagic, magic)

	var cflagLength uint32
	err = binary.Read(f, binary.LittleEndian, &cflagLength)
	assert.NoError(t, err)

	assert.Equal(t, uint32(0), decodeFlag(cflagLength))
	assert.Equal(t, uint32(15839), decodeLength(cflagLength))

	length := decodeLength(cflagLength)
	var upperLength = ((length + uint32(3)) >> uint32(2)) << uint32(2)

	var flag uint32
	err = binary.Read(f, binary.LittleEndian, &flag)
	assert.NoError(t, err)

	assert.Equal(t, uint32(0), flag)

	var label float32
	err = binary.Read(f, binary.LittleEndian, &label)
	assert.NoError(t, err)

	assert.Equal(t, float32(577), label)

	var imageId0 uint64
	err = binary.Read(f, binary.LittleEndian, &imageId0)
	assert.NoError(t, err)

	assert.Equal(t, uint64(16503), imageId0)

	var imageId1 uint64
	err = binary.Read(f, binary.LittleEndian, &imageId1)
	assert.NoError(t, err)

	assert.Equal(t, uint64(0), imageId1)

	headerSize := uint32(unsafe.Sizeof(flag) + unsafe.Sizeof(label) + unsafe.Sizeof(imageId0) + unsafe.Sizeof(imageId1))
	bts := make([]byte, decodeLength(cflagLength)-headerSize)
	_, err = io.ReadFull(f, bts)
	assert.NoError(t, err)

	padding := make([]byte, upperLength-cflagLength)
	io.ReadFull(f, padding)

	img, err := image.Read(bytes.NewBuffer(bts), image.Context(ctx))
	assert.NoError(t, err)
	assert.NotNil(t, img)

	rgbImage := img.(*types.RGBImage)

	rgbaImage := toRGBAImage(rgbImage)

	imageFile, err := os.Create("image.png")
	assert.NoError(t, err)
	defer imageFile.Close()

	err = png.Encode(imageFile, rgbaImage)
	assert.NoError(t, err)
}

func TestMXNetRecordIO(t *testing.T) {
	const testfile = "/home/abduld/data/carml/dldataset/vision/ilsvrc2012_validation_224/imagenet1k-val-224.rec"

	f, err := os.Open(testfile)
	assert.NoError(t, err)

	defer f.Close()

	nextRecord(t, f)
}

func toRGBAImage(rgbImage *types.RGBImage) *goimage.RGBA {
	srcPixels := rgbImage.Pix

	srcHeight := rgbImage.Bounds().Dy()
	srcWidth := rgbImage.Bounds().Dx()

	inputImage := goimage.NewRGBA(goimage.Rect(0, 0, srcWidth, srcHeight))
	for ii := 0; ii < srcHeight; ii++ {
		inOffset := 3 * ii * srcWidth
		inputOffset := ii * inputImage.Stride
		for jj := 0; jj < srcWidth; jj++ {

			inputImage.Pix[inputOffset+0] = srcPixels[inOffset+0]
			inputImage.Pix[inputOffset+1] = srcPixels[inOffset+1]
			inputImage.Pix[inputOffset+2] = srcPixels[inOffset+2]
			inputImage.Pix[inputOffset+3] = 0xFF

			inOffset += 3
			inputOffset += 4
		}

	}

	return inputImage
}
