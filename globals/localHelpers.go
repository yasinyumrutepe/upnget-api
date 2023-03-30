package globals

import "errors"

func FloatToUIntInSlice(arr []interface{}) ([]uint, error) {
	var result []uint
	for _, v := range arr {
		if _, ok := v.(float64); ok {
			result = append(result, uint(v.(float64)))
		} else {
			return nil, errors.New("data is empty")
		}
	}
	return result, nil
}
func Ptr[Y any](a Y) *Y {
	return &a
}
