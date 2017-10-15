package vision

import "github.com/rai-project/image/types"

// ILSVRC2012ValidationLabeledImage ...
type ILSVRC2012ValidationLabeledImage struct {
	label string
	data  *types.RGBImage
}

// Label ...
func (l ILSVRC2012ValidationLabeledImage) Label() string {
	return l.label
}

// Data ...
func (l ILSVRC2012ValidationLabeledImage) Data() (interface{}, error) {
	return l.data, nil
}
