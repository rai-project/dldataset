package vision

import (
	context "context"
)

type base struct {
	ctx            context.Context
	baseWorkingDir string
}

// Category ...
func (base) Category() string {
	return "vision"
}
