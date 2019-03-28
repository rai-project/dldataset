package object_detection

import (
	proto "github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func Get(path string) (*StringIntLabelMap, error) {
	bts, err := Asset(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find %v", path)
	}
	s := new(StringIntLabelMap)
	err = proto.UnmarshalText(string(bts), s)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal %v", path)
	}

	return s, nil
}
