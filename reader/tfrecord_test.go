package reader

import (
	context "context"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/GeertJohan/go-sourcepath"
	"github.com/stretchr/testify/assert"
)

var (
	fixturesPath = filepath.Join(sourcepath.MustAbsoluteDir(), "_fixtures")
)

func TestTFRecord(t *testing.T) {
	reader, err := NewTFRecordReader(filepath.Join(fixturesPath, "cifar10_validation.tfrecord"))
	assert.NoError(t, err)
	assert.NotEmpty(t, reader)

	defer reader.Close()

	rec, err := reader.Next(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, rec)

	assert.Equal(t, rec.ID, uint64(0))

	rec, err = reader.Next(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, rec)

	assert.Equal(t, rec.ID, uint64(1))
	assert.NotEmpty(t, rec.Image.Pix)

	out, _ := os.Create("_fixtures/test.png")
	defer out.Close()

	err = png.Encode(out, rec.Image.ToRGBAImage())
	assert.NoError(t, err)
}
