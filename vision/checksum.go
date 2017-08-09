package vision

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/pkg/errors"
)

type md5sumTy struct{}

var md5sum = md5sumTy{}

func (md5sumTy) Check(reader io.Reader, expected string) (bool, error) {
	hash := md5.New()
	_, err := io.Copy(hash, reader)
	if err != nil {
		return false, errors.Wrap(err, "failed to copy reader to md5 hash")
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	return actual == expected, nil
}

func (md5sumTy) CheckFile(path string, expected string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, errors.Wrapf(err, "failed to open %s while performing md5 checksum", path)
	}
	defer f.Close()
	ok, err := md5sum.Check(f, md5)
	if err != nil {
		return false, errors.Wrapf(err, "unable to perform md5sum on %s", path)
	}
	return ok, nil
}
