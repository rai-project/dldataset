package vision

import (
	"bytes"
	"io"

	context "golang.org/x/net/context"
)

type base struct {
	ctx            context.Context
	baseWorkingDir string
}

func (base) Category() string {
	return "vision"
}

type LabeledImage struct {
	label string
	data  []byte
}

func (l LabeledImage) Label() string {
	return l.label
}

func (l LabeledImage) Data() (io.Reader, error) {
	return bytes.NewBuffer(l.data), nil
}
