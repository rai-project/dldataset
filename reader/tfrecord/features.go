package tfrecord

import (
	"github.com/ubccr/terf"
	protobuf "github.com/ubccr/terf/protobuf"
)

// FeatureBool ...
func FeatureBool(rec *protobuf.Example, key string) bool {
	return FeatureInt(rec, key) == 1
}

// FeatureInt64 ...
func FeatureInt64(rec *protobuf.Example, key string) int64 {
	f, ok := rec.Features.Feature[key]
	if !ok {
		return 0
	}

	val, ok := f.Kind.(*protobuf.Feature_Int64List)
	if !ok {
		return 0
	}

	return val.Int64List.Value[0]
}

// FeatureInt ...
func FeatureInt(rec *protobuf.Example, key string) int {
	return int(FeatureInt64(rec, key))
}

// FeatureInt32 ...
func FeatureInt32(rec *protobuf.Example, key string) int32 {
	return int32(FeatureInt64(rec, key))
}

// FeatureFloat64 ...
func FeatureFloat64(rec *protobuf.Example, key string) float64 {
	return float64(FeatureFloat32(rec, key))
}

// FeatureFloat32 ...
func FeatureFloat32(rec *protobuf.Example, key string) float32 {
	f, ok := rec.Features.Feature[key]
	if !ok {
		return 0
	}

	val, ok := f.Kind.(*protobuf.Feature_FloatList)
	if !ok {
		return 0
	}

	return val.FloatList.Value[0]
}

// FeatureBytes ...
func FeatureBytes(rec *protobuf.Example, key string) []byte {
	return terf.ExampleFeatureBytes(rec, key)
}

// FeatureString ...
func FeatureString(rec *protobuf.Example, key string) string {
	return string(FeatureBytes(rec, key))
}

// FeatureBytesSlice ...
func FeatureBytesSlice(rec *protobuf.Example, key string) [][]byte {
	// TODO: return error if key is not found?
	f, ok := rec.Features.Feature[key]
	if !ok {
		return nil
	}

	val, ok := f.Kind.(*protobuf.Feature_BytesList)
	if !ok {
		return nil
	}
	return val.BytesList.Value
}

// FeatureStringSlice ...
func FeatureStringSlice(rec *protobuf.Example, key string) []string {
	slice := FeatureBytesSlice(rec, key)
	if slice == nil {
		return nil
	}

	res := make([]string, len(slice))
	for ii, val := range slice {
		res[ii] = string(val)
	}

	return res
}

// FeatureInt64Slice ...
func FeatureInt64Slice(rec *protobuf.Example, key string) []int64 {

	f, ok := rec.Features.Feature[key]
	if !ok {
		return nil
	}

	val, ok := f.Kind.(*protobuf.Feature_Int64List)
	if !ok {
		return nil
	}

	return val.Int64List.Value
}

// FeatureIntSlice ...
func FeatureIntSlice(rec *protobuf.Example, key string) []int {
	slice := FeatureInt64Slice(rec, key)
	if slice == nil {
		return nil
	}

	res := make([]int, len(slice))
	for ii, val := range slice {
		res[ii] = int(val)
	}

	return res
}

// FeatureInt32Slice ...
func FeatureInt32Slice(rec *protobuf.Example, key string) []int32 {
	slice := FeatureInt64Slice(rec, key)
	if slice == nil {
		return nil
	}

	res := make([]int32, len(slice))
	for ii, val := range slice {
		res[ii] = int32(val)
	}

	return res
}

// FeatureFloat64Slice ...
func FeatureFloat64Slice(rec *protobuf.Example, key string) []float64 {
	slice := FeatureFloat32Slice(rec, key)
	if slice == nil {
		return nil
	}

	res := make([]float64, len(slice))
	for ii, val := range slice {
		res[ii] = float64(val)
	}

	return res
}

// FeatureFloat32Slice ...
func FeatureFloat32Slice(rec *protobuf.Example, key string) []float32 {
	f, ok := rec.Features.Feature[key]
	if !ok {
		return nil
	}

	val, ok := f.Kind.(*protobuf.Feature_FloatList)
	if !ok {
		return nil
	}

	return val.FloatList.Value
}
