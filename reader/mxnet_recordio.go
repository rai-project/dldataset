package reader

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"unsafe"

	context "context"

	"github.com/pkg/errors"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
)

const (
	kMagic = uint32(0xced7230a)
)

type RecordIOReader struct {
	r io.ReadCloser
}

type Record struct {
	ID         uint64
	LabelIndex float32
	Image      *types.RGBImage
}

func NewRecordIOReader(path string) (*RecordIOReader, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open %v", path)
	}
	return &RecordIOReader{
		r: r,
	}, nil
}

func (r *RecordIOReader) Next(ctx context.Context) (*Record, error) {
	f := r.r

	var magic uint32
	err := binary.Read(f, binary.LittleEndian, &magic)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read magic")
	}
	if magic != kMagic {
		return nil, errors.New("invalid magic number")
	}

	var cflagLength uint32
	err = binary.Read(f, binary.LittleEndian, &cflagLength)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read cflag / length")
	}

	cflag := decodeFlag(cflagLength)
	if cflag != 0 {
		return nil, errors.New("only cflag==0 is currently supported")
	}

	length := decodeLength(cflagLength)

	paddedLength := ((length + uint32(3)) >> uint32(2)) << uint32(2)

	var flag uint32
	err = binary.Read(f, binary.LittleEndian, &flag)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read image header flag")
	}

	var label float32
	err = binary.Read(f, binary.LittleEndian, &label)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read image header label")
	}

	var imageId0 uint64
	err = binary.Read(f, binary.LittleEndian, &imageId0)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read image header imageId0")
	}

	var imageId1 uint64
	err = binary.Read(f, binary.LittleEndian, &imageId1)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read image header imageId1")
	}

	headerSize := uint32(unsafe.Sizeof(flag) +
		unsafe.Sizeof(label) +
		unsafe.Sizeof(imageId0) +
		unsafe.Sizeof(imageId1))

	bts := make([]byte, length-headerSize)
	_, err = io.ReadFull(f, bts)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read image jpeg data")
	}

	padding := make([]byte, paddedLength-length)
	io.ReadFull(f, padding)

	img, err := image.Read(bytes.NewBuffer(bts), image.Context(nil))
	if err != nil {
		return nil, err
	}

	rgbImage, ok := img.(*types.RGBImage)
	if !ok {
		return nil, errors.Errorf("expecting an rgb image")
	}

	return &Record{
		ID:         imageId1,
		LabelIndex: label,
		Image:      rgbImage,
	}, nil

}

func (r *RecordIOReader) Close() error {
	return r.r.Close()
}

/*!
 * \brief decode the flag part of lrecord
 * \param rec the lrecord
 * \return the flag
 */
func decodeFlag(rec uint32) uint32 {
	return (rec >> uint32(29)) & uint32(7)
}

/*!
 * \brief decode the length part of lrecord
 * \param rec the lrecord
 * \return the length
 */
func decodeLength(rec uint32) uint32 {
	return rec & ((uint32(1) << uint32(29)) - 1)
}
