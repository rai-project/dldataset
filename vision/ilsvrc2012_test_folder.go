package vision

// import (
// 	"net/http"
// 	"os"
// 	"path"
// 	"path/filepath"
// 	"strings"

// 	context "context"

// 	"github.com/pkg/errors"
// 	"github.com/rai-project/config"
// 	"github.com/rai-project/dldataset"
// 	"github.com/rai-project/downloadmanager"
// 	"github.com/rai-project/image"
// 	"github.com/rai-project/image/types"
// )

// var iLSVRC2012TestFolder *ILSVRC2012TestFolder

// // ILSVRC2012TestFolder ...
// type ILSVRC2012TestFolder struct {
// 	base
// 	baseURL   string
// 	filePaths []string
// 	fileURLs  map[string]string
// 	data      map[string]ILSVRC2012TestLabeledImage
// }

// // New ...
// func (d *ILSVRC2012TestFolder) New(ctx context.Context) (dldataset.Dataset, error) {
// 	return iLSVRC2012TestFolder, nil
// }

// func (d *ILSVRC2012TestFolder) Load(ctx context.Context) error {
// 	return nil
// }

// func (d *ILSVRC2012TestFolder) workingDir() string {
// 	category := strings.ToLower(d.Category())
// 	name := strings.ToLower(d.Name())
// 	return filepath.Join(d.baseWorkingDir, category, name)
// }

// // Name ...
// func (d *ILSVRC2012TestFolder) Name() string {
// 	return "ilsvrc2012_test_folder"
// }

// // CanonicalName ...
// func (d *ILSVRC2012TestFolder) CanonicalName() string {
// 	category := strings.ToLower(d.Category())
// 	name := strings.ToLower(d.Name())
// 	key := path.Join(category, name)
// 	return key
// }

// // Download ...
// func (d *ILSVRC2012TestFolder) Download(ctx context.Context) error {
// 	return nil
// }

// // List ...
// func (d *ILSVRC2012TestFolder) List(ctx context.Context) ([]string, error) {
// 	return d.filePaths, nil
// }

// // GetWithoutDownloadManager ...
// func (d *ILSVRC2012TestFolder) GetWithoutDownloadManager(ctx context.Context, name string) (dldataset.LabeledData, error) {
// 	fileURL, ok := d.fileURLs[name]
// 	if !ok {
// 		return nil, errors.Errorf("the file path %v for the dataset %v was not found", name, d.CanonicalName())
// 	}
// 	req, err := http.Get(fileURL)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to perform http get request to %v", fileURL)
// 	}
// 	defer req.Body.Close()

// 	img, err := image.Read(req.Body)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to read image from %v", fileURL)
// 	}

// 	if _, ok := img.(*types.RGBImage); !ok {
// 		return nil, errors.Wrapf(err, "failed to read rgb image from %v", fileURL)
// 	}

// 	label := path.Dir(name)

// 	return &ILSVRC2012TestLabeledImage{
// 		data:  img.(*types.RGBImage),
// 		label: label,
// 	}, nil
// }

// // Get ...
// func (d *ILSVRC2012TestFolder) Get(ctx context.Context, name string) (dldataset.LabeledData, error) {
// 	fileURL, ok := d.fileURLs[name]
// 	if !ok {
// 		return nil, errors.Errorf("the file path %v for the dataset %v was not found", name, d.CanonicalName())
// 	}

// 	workingDir := d.workingDir()
// 	downloadedFileName := filepath.Join(workingDir, name)
// 	downloadedFileName, err := downloadmanager.DownloadFile(
// 		fileURL,
// 		downloadedFileName,
// 		downloadmanager.Context(ctx),
// 		downloadmanager.Cache(true),
// 		downloadmanager.CheckMD5Sum(false),
// 	)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to download %v", fileURL)
// 	}

// 	f, err := os.Open(downloadedFileName)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to open %v", downloadedFileName)
// 	}

// 	defer f.Close()

// 	img, err := image.Read(f, image.Context(ctx))
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to read image from %v", fileURL)
// 	}

// 	if _, ok := img.(*types.RGBImage); !ok {
// 		return nil, errors.Wrapf(err, "failed to read rgb image from %v", fileURL)
// 	}

// 	label := path.Dir(name)

// 	return &ILSVRC2012TestLabeledImage{
// 		data:  img.(*types.RGBImage),
// 		label: label,
// 	}, nil
// }

// func (d *ILSVRC2012TestFolder) Next(ctx context.Context) (dldataset.LabeledData, error) {
// 	return nil, errors.New("next iterator is not implemented for " + d.CanonicalName())
// }

// // Close ...
// func (d *ILSVRC2012TestFolder) Close() error {
// 	return nil
// }

// func init() {
// 	const fileListPath = "/vision/support/ilsvrc2012_test_file_list.txt"
// 	const baseURL = "http://store.carml.org.s3.amazonaws.com/datasets/ilsvrc2012_test/"
// 	config.AfterInit(func() {

// 		filePaths := strings.Split(_escFSMustString(false, fileListPath), "\n")

// 		fileURLs := map[string]string{}
// 		for _, p := range filePaths {
// 			fileURLs[p] = baseURL + p
// 		}

// 		iLSVRC2012TestFolder = &ILSVRC2012TestFolder{
// 			base: base{
// 				ctx:            context.Background(),
// 				baseWorkingDir: filepath.Join(dldataset.Config.WorkingDirectory, "dldataset"),
// 			},
// 			baseURL:   baseURL,
// 			fileURLs:  fileURLs,
// 			filePaths: filePaths,
// 		}
// 		dldataset.Register(iLSVRC2012TestFolder)
// 	})
// }
