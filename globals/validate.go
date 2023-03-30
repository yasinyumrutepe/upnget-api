package globals

import (
	"reflect"

	"github.com/go-playground/validator"
)

func ValidateSlice[T any](slice []T) (validates []T, errList []map[string]interface{}) {
	var validate = validator.New()
	for _, v := range slice {
		if err := validate.Struct(v); err == nil {
			validates = append(validates, v)
		} else {
			errList = append(errList, map[string]interface{}{
				"message": err.Error(),
				"value":   v,
			})
		}
	}
	return
}

func ValidateStruct[T any](str T) (errList []map[string]interface{}) {
	var validate = validator.New()
	err := validate.Struct(str)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			errList = append(errList, map[string]interface{}{
				"cause_of_error": v.Tag(),
				"failed_field":   v.Field(),
				"value":          v.Value(),
			})
		}
		return errList
	}
	return nil
}
func InSlice(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}
	return
}

// func ValidateMap(m map[string]interface{}, rules map[string]interface{}) (errList []map[string]interface{}) {
// 	var validate = validator.New()
// 	errs := validate.ValidateMap(m, rules)
// 	if len(errs) > 0 {
// 		for i, v := range errs {
// 			var tmpSlice []map[string]interface{}
// 			for _, err := range v.(validator.ValidationErrors) {
// 				tmpSlice = append(tmpSlice, map[string]interface{}{
// 					"cause_of_error": err.Tag(),
// 					"failed_field":   i,
// 					"value":          err.Value(),
// 				})
// 			}
// 			errList = append(errList, tmpSlice...)
// 		}
// 		return errList
// 	}
// 	return nil
// }
