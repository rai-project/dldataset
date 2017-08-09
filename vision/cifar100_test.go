package vision

import (
	"testing"

	context "golang.org/x/net/context"

	"github.com/k0kubun/pp"
	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
)

func TestDownloadCIFAR100(t *testing.T) {
	ctx := context.Background()

	cifar10, err := dldataset.Get("vision", "cifar100")
	assert.NoError(t, err)
	assert.NotEmpty(t, cifar10)

	err = cifar10.Download(ctx)
	assert.NoError(t, err)

	lbl, err := cifar10.Get(ctx, "1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl)

	pp.Println(lbl)

}
