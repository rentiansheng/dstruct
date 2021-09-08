package dstruct

import (
	"reflect"
)

/**
 * @Description:
 * @Date: 2021/9/7 19:58
 */


// validation 进行数据校验方法
type validation interface {
	Validate() error
}

var validateType = reflect.TypeOf((*validation)(nil)).Elem()

