package vision

import (
	"testing"

	context "golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
)

func TestDownloadCIFAR100(t *testing.T) {
	ctx := context.Background()

	cifar100, err := dldataset.Get("vision", "cifar100")
	assert.NoError(t, err)
	assert.NotEmpty(t, cifar100)

	defer cifar100.Close()

	err = cifar100.Download(ctx)
	assert.NoError(t, err)

	fileList, err := cifar100.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileList)

	lbl, err := cifar100.Get(ctx, "train/1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl)

	lbl, err = cifar100.Get(ctx, "test/1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl)

	// pp.Println(lbl)

}
