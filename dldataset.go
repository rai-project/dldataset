package dldataset

import (
	"io"

	context "golang.org/x/net/context"
)

type LabeledData interface {
	Label() string
	Data() (interface{}, error)
}

type Dataset interface {
	New(ctx context.Context) (Dataset, error)
	Category() string
	Name() string
	CanonicalName() string
	Download(ctx context.Context) error
	List(ctx context.Context) ([]string, error)
	Get(ctx context.Context, name string) (LabeledData, error)
	io.Closer
}
