package vision

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/dldataset"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	context "golang.org/x/net/context"
)

var iLSVRC2012ValidationFolder *ILSVRC2012ValidationFolder

// ILSVRC2012ValidationFolder ...
type ILSVRC2012ValidationFolder struct {
	base
	baseURL   string
	filePaths []string
	fileURLs  map[string]string
	data      map[string]ILSVRC2012ValidationLabeledImage
}

// New ...
func (d *ILSVRC2012ValidationFolder) New(ctx context.Context) (dldataset.Dataset, error) {
	return iLSVRC2012ValidationFolder, nil
}

func (d *ILSVRC2012ValidationFolder) workingDir() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	return filepath.Join(d.baseWorkingDir, category, name)
}

// Name ...
func (d *ILSVRC2012ValidationFolder) Name() string {
	return "ilsvrc2012_validation_folder"
}

// CanonicalName ...
func (d *ILSVRC2012ValidationFolder) CanonicalName() string {
	category := strings.ToLower(d.Category())
	name := strings.ToLower(d.Name())
	key := path.Join(category, name)
	return key
}

// Download ...
func (d *ILSVRC2012ValidationFolder) Download(ctx context.Context) error {
	return nil
}

// List ...
func (d *ILSVRC2012ValidationFolder) List(ctx context.Context) ([]string, error) {
	return d.filePaths, nil
}

// GetWithoutDownloadManager ...
func (d *ILSVRC2012ValidationFolder) GetWithoutDownloadManager(ctx context.Context, name string) (dldataset.LabeledData, error) {
	fileURL, ok := d.fileURLs[name]
	if !ok {
		return nil, errors.Errorf("the file path %v for the dataset %v was not found", name, d.CanonicalName())
	}
	req, err := http.Get(fileURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to perform http get request to %v", fileURL)
	}
	defer req.Body.Close()

	img, err := image.Read(req.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read image from %v", fileURL)
	}

	if _, ok := img.(*types.RGBImage); !ok {
		return nil, errors.Wrapf(err, "failed to read rgb image from %v", fileURL)
	}

	label := path.Dir(name)

	return &ILSVRC2012ValidationLabeledImage{
		data:  img.(*types.RGBImage),
		label: label,
	}, nil
}

// Get ...
func (d *ILSVRC2012ValidationFolder) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
	fileURL, ok := d.fileURLs[name]
	if !ok {
		return nil, errors.Errorf("the file path %v for the dataset %v was not found", name, d.CanonicalName())
	}

	workingDir := d.workingDir()
	downloadedFileName := filepath.Join(workingDir, name)
	downloadedFileName, err := downloadmanager.DownloadFile(
		fileURL,
		downloadedFileName,
		downloadmanager.Context(ctx),
		downloadmanager.Cache(true),
		downloadmanager.CheckMD5Sum(false),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download %v", fileURL)
	}

	f, err := os.Open(downloadedFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %v", downloadedFileName)
	}

	defer f.Close()

	img, err := image.Read(f, image.Context(ctx))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read image from %v", fileURL)
	}

	if _, ok := img.(*types.RGBImage); !ok {
		return nil, errors.Wrapf(err, "failed to read rgb image from %v", fileURL)
	}

	label := path.Dir(name)

	return &ILSVRC2012ValidationLabeledImage{
		data:  img.(*types.RGBImage),
		label: label,
	}, nil
}

// Close ...
func (d *ILSVRC2012ValidationFolder) Close() error {
	return nil
}

func init() {
	const fileListPath = "/vision/support/ilsvrc2012_validation_file_list.txt"
	const baseURL = "http://store.carml.org.s3.amazonaws.com/datasets/ilsvrc2012_validation/"
	config.AfterInit(func() {

		filePaths := strings.Split(_escFSMustString(false, fileListPath), "\n")

		fileURLs := map[string]string{}
		for _, p := range filePaths {
			fileURLs[p] = baseURL + p
		}

		iLSVRC2012ValidationFolder = &ILSVRC2012ValidationFolder{
			base: base{
				ctx:            context.Background(),
				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
			},
			baseURL:   baseURL,
			fileURLs:  fileURLs,
			filePaths: filePaths,
		}
		dldataset.Register(iLSVRC2012ValidationFolder)
	})
}
