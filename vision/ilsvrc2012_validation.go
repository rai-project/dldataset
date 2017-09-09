package vision

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/image/types"
	context "golang.org/x/net/context"
)

var iLSVRC2012Validation *ILSVRC2012Validation

type ILSVRC2012Validation struct {
	base
	baseURL   string
	filePaths []string
	data      map[string]ILSVRC2012ValidationLabeledImage
}

type ILSVRC2012ValidationLabeledImage struct {
	label string
	data  *types.RGBImage
}

func (l ILSVRC2012ValidationLabeledImage) Label() string {
	return l.label
}

func (l ILSVRC2012ValidationLabeledImage) Data() (interface{}, error) {
	return l.data, nil
}

func (d *ILSVRC2012Validation) New(ctx context.Context) (dldataset.Dataset, error) {
	return iLSVRC2012Validation, nil
}

func (d *ILSVRC2012Validation) Name() string {
	return "ilsvrc2012_validation"
}

func (d *ILSVRC2012Validation) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

func (d *ILSVRC2012Validation) Download(ctx context.Context) error {
	return nil
}

func (d *ILSVRC2012Validation) List(ctx context.Context) ([]string, error) {
	return d.filePaths, nil
}

func (d *ILSVRC2012Validation) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	panic("TODO ILSVRC2012Validation/GET")
	return nil, nil
}

func (d *ILSVRC2012Validation) Close() error {
	return nil
}

func init() {
	config.AfterInit(func() {
		const fileListPath = "/vision/support/ilsvrc2012_validation_file_list.txt"
		filePaths := strings.Split(_escFSMustString(false, fileListPath), "\n")
		iLSVRC2012Validation = &ILSVRC2012Validation{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			baseURL:   "http://store.carml.org.s3.amazonaws.com/datasets/ilsvrc2012_validation/",
			filePaths: filePaths,
		}
		dldataset.Register(iLSVRC2012Validation)
	})
}
