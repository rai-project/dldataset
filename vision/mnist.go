package vision

import (
	"path"
	"strings"

	context "golang.org/x/net/context"

	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	mnistLoader "github.com/unixpickle/mnist"
)

type Mnist struct {
	base
	trainingData mnistLoader.DataSet
	testData     mnistLoader.DataSet
}

var mnist *Mnist

func (*Mnist) Name() string {
	return "Mnist"
}

func (d *Mnist) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *Mnist) New(ctx context.Context) (dldataset.Dataset, error) {
	return mnist, nil
}

func (d *Mnist) Download(ctx context.Context) error {
	return nil
}

func (d *Mnist) List(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (d *Mnist) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	return nil, nil
}

func (d *Mnist) Close() error {
	return nil
}

func init() {
	config.AfterInit(func() {
		mnist = &Mnist{
			base: base{
				ctx: context.Background(),
			},
			trainingData: mnistLoader.LoadTestingDataSet(),
			testData:     mnistLoader.LoadTestingDataSet(),
		}
		dldataset.Register(mnist)
	})
}
