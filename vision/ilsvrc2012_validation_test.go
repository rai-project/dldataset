package vision

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

// TestILSVRC2012Validation ...
func TestILSVRC2012Validation(t *testing.T) {

	ctx := context.Background()

	ilsvrc, err := dldataset.Get("vision", "ilsvrc2012_validation_224")
	assert.NoError(t, err)
	assert.NotEmpty(t, ilsvrc)

	defer ilsvrc.Close()

	err = ilsvrc.Download(ctx)
	assert.NoError(t, err)

	lst, err := ilsvrc.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, lst)
	assert.Equal(t, len(lst), 50000)
	assert.Equal(t, "ILSVRC2012_val_00016503.JPEG", lst[0])
	assert.Equal(t, "ILSVRC2012_val_00035805.JPEG", lst[1])

	pp.Println(ilsvrc.Get(ctx, "ILSVRC2012_val_00016503.JPEG"))
}
