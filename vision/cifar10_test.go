package vision

import (
	"testing"

	context "golang.org/x/net/context"

	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
)

func TestDownloadCIFAR10(t *testing.T) {
	ctx := context.Background()

	cifar10, err := dldataset.Get("vision", "cifar10")
	assert.NoError(t, err)
	assert.NotEmpty(t, cifar10)

	cifar10.Download(ctx)
}
