package vision

import (
	context "golang.org/x/net/context"
)

type base struct {
	ctx            context.Context
	baseWorkingDir string
}

func (base) Category() string {
	return "vision"
}
