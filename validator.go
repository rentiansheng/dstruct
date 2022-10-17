package dstruct

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/10/17
    @desc:

***************************/

var validate = validator.New()

// validateStruct 不会递归调用
func validateStruct(value reflect.Value) error {
	if value.Interface() == nil || value.IsZero() {
		return nil
	}

	switch value.Kind() {
	case reflect.Ptr:
		return validate.Struct(value.Elem().Interface())
	case reflect.Struct:
		return validate.Struct(value.Interface())
	/*case reflect.Slice, reflect.Array:

	for i := 0; i < value.Len(); i++ {
		if err := validateStruct(value.Index(i)); err != nil {
			return err
		}
	}
	return nil*/
	default:
		return nil
	}

}
