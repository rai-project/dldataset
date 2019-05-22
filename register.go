package dldataset

import (
	"errors"
	"path"
	"strings"

	"github.com/k0kubun/pp"

	"golang.org/x/sync/syncmap"
)

var datasets syncmap.Map

// Get ...
func Get(category, name string) (Dataset, error) {
	category = strings.ToLower(category)
	name = strings.ToLower(name)
	key := path.Join(category, name)
	val, ok := datasets.Load(key)
	if !ok {
		log.WithField("category", category).
			WithField("name", name).
			Warn("cannot find dataset")
		return nil, errors.New("cannot find dataset")
	}
	dataset, ok := val.(Dataset)
	if !ok {
		log.WithField("category", category).
			WithField("name", name).
			Warn("invalid dataset")
		return nil, errors.New("invalid dataset")
	}
	return dataset, nil
}

// Register ...
func Register(d Dataset) {
	pp.Println(d.CanonicalName())
	datasets.Store(d.CanonicalName(), d)
}

// Datasets ...
func Datasets() []string {
	names := []string{}
	datasets.Range(func(key, _ interface{}) bool {
		if name, ok := key.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return names
}
