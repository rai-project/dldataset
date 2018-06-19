package dldataset

import (
	"io"

	context "context"
)

// LabeledData ...
type LabeledData interface {
	Label() string
	Data() (interface{}, error)
}

// Dataset ...
type Dataset interface {
	New(ctx context.Context) (Dataset, error)
	Category() string
	Name() string
	CanonicalName() string
	Download(ctx context.Context) error
	List(ctx context.Context) ([]string, error)
	Load(ctx context.Context) error
	Get(ctx context.Context, name string) (LabeledData, error)
	Next(ctx context.Context) (LabeledData, error)
	io.Closer
}
