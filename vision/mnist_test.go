package vision

import (
	"testing"

	context "context"

	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
)

// TestDownloadMNIST ...
func TestDownloadMNIST(t *testing.T) {
	ctx := context.Background()

	MNIST, err := dldataset.Get("vision", "mnist")
	assert.NoError(t, err)
	assert.NotEmpty(t, MNIST)

	defer MNIST.Close()

	err = MNIST.Download(ctx)
	assert.NoError(t, err)

	fileList, err := MNIST.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileList)

	lbl, err := MNIST.Get(ctx, "train/1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl)

	lbl, err = MNIST.Get(ctx, "test/1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl)

	// pp.Println(lbl)

}
