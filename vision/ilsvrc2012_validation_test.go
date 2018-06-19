package vision

import (
	"testing"

	context "context"
	"github.com/rai-project/dldataset"
	"github.com/stretchr/testify/assert"
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

	err = ilsvrc.Load(ctx)
	assert.NoError(t, err)

	lst, err := ilsvrc.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, lst)
	assert.Equal(t, len(lst), 50000)
	assert.Equal(t, "ILSVRC2012_val_00016503.JPEG", lst[0])
	assert.Equal(t, "ILSVRC2012_val_00035805.JPEG", lst[1])

	for ii := 0; ii < len(lst); ii++ {
		data, err := ilsvrc.Next(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.IsType(t, &iLSVRC2012ValidationRecordIOLabeledData{}, data)
	}
}
