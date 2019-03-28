package dldataset

import (
	"io"

	context "context"

	"github.com/rai-project/dlframework"
)

// LabeledData ...
type LabeledData interface {
	Label() string
	Feature() *dlframework.Feature
	Features() dlframework.Features
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
