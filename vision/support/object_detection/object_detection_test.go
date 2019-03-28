package object_detection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPascalLabelMap(t *testing.T) {
	mp, err := Get("pascal_label_map.pbtxt")
	assert.NoError(t, err)
	assert.NotEmpty(t, mp)

	items := mp.GetItem()

	assert.NotEmpty(t, items)

	fst := items[0]

	assert.Equal(t, fst.Name, "aeroplane")
	assert.Equal(t, fst.Id, int32(1))
	assert.Equal(t, fst.DisplayName, "")
}
