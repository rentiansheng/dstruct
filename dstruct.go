package dstruct

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrType     error = errors.New("type mismatch")
	ErrNotFound error = errors.New("not found field")
)

type DMode int64

const (
	// 默认模式，进行严格的数据校验
	strictMode DMode = iota
	// 不进行数据格式校验，数据不一致的时候，返回需要默认值
	broadMode
)

type DStruct struct {
	mode DMode

	fields map[string]reflect.Type
	/**** 为了减少转换这里单独存储不同类型数据 ****/
	kv map[string]interface{}
}

func (d *DStruct) SetFields(fields map[string]reflect.Type) {
	if d.fields == nil {
		d.fields = make(map[string]reflect.Type)
	}
	for key, val := range fields {
		d.fields[key] = val
	}
}

func (d *DStruct) ResetFields(fields map[string]reflect.Type) {

	d.fields = make(map[string]reflect.Type)
	for key, val := range fields {
		d.fields[key] = val
	}
}

/*
func (d DStruct) FieldType(name string) reflect.Type {
	if DType, ok := d.fields[name]; ok {
		return DType
	}
	return reflect.Invalid
}

func (d DStruct) Int64(name string) (int64, error) {
	if err := d.checkType(name, reflect.Int64); err != nil {
		return 0, err
	}

	return d.int64Vals[name], nil

}

func (d DStruct) Str(name string) (string, error) {
	if err := d.checkType(name, reflect.String); err != nil {
		return "", err
	}

	return d.strVals[name], nil
}

func (d DStruct) Bool(name string) (bool, error) {
	if err := d.checkType(name, reflect.Bool); err != nil {
		return false, err
	}

	return d.boolVals[name], nil
}

// 检查字段类型
func (d DStruct) checkType(name string, DType reflect.Type) error {

	if d.mode == broadMode {
		return nil
	}
	realType, ok := d.fields[name]
	// 字段不存在
	if !ok {
		return ErrNotFound
	}
	// 类型不匹配
	if DType != realType {
		return ErrType
	}

	return nil
}

func (d DStruct) Types() map[string]int64 {
	fields := make(map[string]int64, len(d.fields))
	for key, t := range d.fields {
		fields[key] = int64(t)
	}
	return fields
}
*/
func (d DStruct) String() string {

	str, _ := json.Marshal(d.kv)
	return string(str)
}
