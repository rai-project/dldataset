package vision

import (
	"testing"

	"github.com/rai-project/dldataset"
	"github.com/rai-project/image/types"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestILSVRC2012Validation(t *testing.T) {

	ctx := context.Background()

	ilsvrc, err := dldataset.Get("vision", "ilsvrc2012_validation")
	assert.NoError(t, err)
	assert.NotEmpty(t, ilsvrc)

	defer ilsvrc.Close()

	err = ilsvrc.Download(ctx)
	assert.NoError(t, err)

	fileList, err := ilsvrc.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileList)

	for ii := 0; ii < 10; ii++ {
		lbl, err := ilsvrc.Get(ctx, fileList[ii])
		assert.NoError(t, err)
		assert.NotEmpty(t, lbl)

		data, err := lbl.Data()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.IsType(t, &types.RGBImage{}, data)
	}
}
