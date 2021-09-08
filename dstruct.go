package dstruct

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrType     error = errors.New("type mismatch")
	ErrNotFound error = errors.New("not found field")
	ErrNotSet         = errors.New("eflect.Value.Addr of unaddressable value")
)

type DMode int64

const (
	// 默认模式，进行严格的数据校验
	strictMode DMode = iota
	// 不进行数据格式校验，数据不一致的时候，返回需要默认值
	broadMode
)



var (
	TypeString  = reflect.TypeOf("")
	TypeInvalid = reflect.TypeOf(nil)
	TypeBool    = reflect.TypeOf(bool(true))
	TypeInt     = reflect.TypeOf(int(0))
	TypeInt8    = reflect.TypeOf(int8(0))
	TypeInt16   = reflect.TypeOf(int16(0))
	TypeInt32   = reflect.TypeOf(int32(0))
	TypeInt64   = reflect.TypeOf(int64(0))
	TypeUint    = reflect.TypeOf(uint(0))
	TypeUint8   = reflect.TypeOf(uint8(0))
	TypeUint16  = reflect.TypeOf(uint16(0))
	TypeUint32  = reflect.TypeOf(uint32(0))
	TypeUint64  = reflect.TypeOf(uint64(0))
	//uintptr
	TypeFloat32 = reflect.TypeOf(float32(0))
	TypeFloat64 = reflect.TypeOf(float32(0))
	//complex64,
	//complex128,
	TypeArrayInt64 = reflect.TypeOf([]int64{})
	TypeArrayInt   = reflect.TypeOf([]int{})
	TypeArrayStr   = reflect.TypeOf([]string{})

	//chan,
	//func,
	TypeMapStr = reflect.TypeOf(map[string]string{})
	//TypePtr    = reflect.Type()
	//TypeSlice  = reflect.Type()
	//TypeStruct = reflect.Type()
	//unsafe.Pointer
)

type DStruct struct {
	mode DMode

	fields map[string]reflect.Type
	kv     map[string]interface{}

	// json Unmarshal 特有字段，数组类型是否使用jsonNumber
	jsonNumber bool
}

func (d *DStruct) init() {
	d.kv = make(map[string]interface{}, len(d.fields))
	for field, typ := range d.fields {
		if typ == nil {
			d.kv[field] = nil
		} else {
			d.kv[field] = reflect.New(typ).Elem().Interface()
		}
	}
}

func (d *DStruct) SetFields(fields map[string]reflect.Type) {
	if d.fields == nil {
		d.fields = make(map[string]reflect.Type)
	}
	for key, typ := range fields {
		d.fields[key] = typ
	}
}

func (d *DStruct) SetOneFields(field string, typ reflect.Type) {
	if d.fields == nil {
		d.fields = make(map[string]reflect.Type)
	}
	d.fields[field] = typ
}

func (d *DStruct) ResetFields(fields map[string]reflect.Type) {

	d.fields = make(map[string]reflect.Type)
	for key, val := range fields {
		d.fields[key] = val
	}
}

func (d DStruct) FieldType(name string) reflect.Type {
	if typ, ok := d.fields[name]; ok {
		return typ
	}
	return nil
}

func (d DStruct) Int64(name string) (int64, error) {
	if err := d.checkType(name, TypeInt64); err != nil {
		return 0, err
	}

	return d.kv[name].(int64), nil

}

func (d DStruct) Str(name string) (string, error) {
	if err := d.checkType(name, TypeString); err != nil {
		return "", err
	}

	return d.kv[name].(string), nil
}

func (d DStruct) Bool(name string) (bool, error) {
	if err := d.checkType(name, TypeBool); err != nil {
		return false, err
	}

	return d.kv[name].(bool), nil
}

func (d *DStruct) Value(name string, val interface{}) (bool, error) {
	iVal, ok := d.kv[name]
	if !ok {
		return false, nil
	}
	valueOf := reflect.ValueOf(val)
	iValueOf := reflect.ValueOf(iVal)

	// 处理多次reflect.Set引起类型错误的问题
	if valOf, ok := iVal.(reflect.Value); ok {
		iValueOf = valOf
	}

	if !valueOf.Elem().CanSet() || valueOf.Kind() != reflect.Ptr {
		return false, ErrNotSet
	}

	if err := d.checkType(name, valueOf.Elem().Type()); err != nil {
		return false, err
	}

	if valueOf.Elem().Type().Kind() == reflect.Interface || valueOf.Elem().Type() == iValueOf.Type() {
		valueOf.Elem().Set(iValueOf)
	}
	return true, nil
}

// 检查字段类型
func (d DStruct) checkType(name string, typ reflect.Type) error {

	if d.mode == broadMode {
		return nil
	}
	realType, ok := d.fields[name]
	// 字段不存在
	if !ok {
		return ErrNotFound
	}

	if typ.Kind() == reflect.Interface {
		return nil
	}

	// 类型不匹配
	if typ != realType {
		return ErrType
	}

	return nil
}

// Clone 克隆一个新的结构体，不包含值
func (d *DStruct) Clone() *DStruct {
	newD := &DStruct{}
	newD.mode = d.mode
	newD.fields = d.fields
	newD.jsonNumber = d.jsonNumber
	newD.kv = nil
	for field, typ := range d.fields {
		newD.fields[field] = typ
	}
	return newD
}

func (d DStruct) String() string {

	str, err := json.Marshal(d.kv)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

