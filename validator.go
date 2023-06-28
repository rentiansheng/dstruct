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

type validatorType func(value reflect.Value) error

var validate = validator.New()

// validateStruct 不会递归调用
func validateStruct(value reflect.Value) error {
	if value.Interface() == nil || value.IsZero() {
		return nil
	}

	data := value.Interface()
	if data == nil {
		return nil
	}
	if v, ok := data.(Validate); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	switch value.Kind() {
	case reflect.Ptr:
		data := value.Elem().Interface()
		if data == nil {
			return nil
		}
		return validate.Struct(data)
	case reflect.Struct:
		return validate.Struct(data)
	case reflect.Slice, reflect.Array:

		for i := 0; i < value.Len(); i++ {
			if err := validateStruct(value.Index(i)); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}

}
