package vision

import (
	"image"
	"path"
	"strconv"
	"strings"

	context "golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/image/types"
	mnistLoader "github.com/unixpickle/mnist"
)

type MNIST struct {
	base
	trainingData mnistLoader.DataSet
	testData     mnistLoader.DataSet
}

var mnist *MNIST

type MNISTLabeledImage struct {
	label string
	data  *types.RGBImage
}

func (l MNISTLabeledImage) Label() string {
	return l.label
}

func (l MNISTLabeledImage) Data() (interface{}, error) {
	return l.data, nil
}

func (*MNIST) Name() string {
	return "MNIST"
}

func (d *MNIST) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *MNIST) New(ctx context.Context) (dldataset.Dataset, error) {
	return mnist, nil
}

func (d *MNIST) Download(ctx context.Context) error {
	return nil
}

func (d *MNIST) List(ctx context.Context) ([]string, error) {
	lst := []string{}
	for ii := range d.trainingData.Samples {
		lst = append(lst, "train/"+strconv.Itoa(ii))
	}
	for ii := range d.testData.Samples {
		lst = append(lst, "test/"+strconv.Itoa(ii))
	}
	return lst, nil
}

func (d *MNIST) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	var dataset mnistLoader.DataSet
	if strings.HasPrefix(name, "train/") {
		name = strings.TrimPrefix(name, "train/")
		dataset = d.trainingData
	} else if strings.HasPrefix(name, "test/") {
		name = strings.TrimPrefix(name, "test/")
		dataset = d.trainingData
	} else {
		return nil, errors.Errorf("cannot find %s in the mnist dataset", name)
	}
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil, errors.Errorf("expecting an integer, but got %s", name)
	}
	if idx >= len(dataset.Samples) {
		return nil, errors.Errorf("the index %d is out of range %d", idx, len(dataset.Samples))
	}

	elem := dataset.Samples[idx]

	img := types.NewRGBImage(image.Rect(0, 0, dataset.Width, dataset.Height))
	data := img.Pix

	for ii, intensity := range elem.Intensities {
		if intensity == 1 {
			data[3*ii+0] = byte(1)
			data[3*ii+1] = byte(1)
			data[3*ii+2] = byte(1)
		} else {
			data[3*ii+0] = byte(0)
			data[3*ii+1] = byte(0)
			data[3*ii+2] = byte(0)
		}
	}

	return &MNISTLabeledImage{
		data:  img,
		label: strconv.Itoa(elem.Label),
	}, nil
}

func (d *MNIST) Close() error {
	return nil
}

func init() {
	config.AfterInit(func() {
		mnist = &MNIST{
			base: base{
				ctx: context.Background(),
			},
			trainingData: mnistLoader.LoadTestingDataSet(),
			testData:     mnistLoader.LoadTestingDataSet(),
		}
		dldataset.Register(mnist)
	})
}
