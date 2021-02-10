package dstruct

import (
	"encoding/json"
	"reflect"
	"strings"

	"dstruct/jsoniter"
)

func (d *DStruct) UnmarshalJSON(data []byte) error {

	iter := jsoniter.BorrowIterator(data, d.jsonNumber)
	defer jsoniter.ReturnIterator(iter)
	d.jsonDecode(iter)
	return iter.Error

}

func (d *DStruct) JSONNumber() {
	d.jsonNumber = true
}

func (d *DStruct) jsonDecode(iter *jsoniter.Iterator) {

	d.init()
	c := iter.NextToken()
	if c == 'n' {
		iter.SkipThreeBytes('u', 'l', 'l')
		return
	}

	if c != '{' {
		iter.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
		return
	}
	c = iter.NextToken()
	if c == '}' {
		return
	}

	iter.UnreadByte()
	for c = ','; c == ','; c = iter.NextToken() {
		key := iter.ReadString()
		if iter.Error != nil {
			return
		}

		c = iter.NextToken()
		if c != ':' {
			iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			return
		}
		if typ, ok := d.fields[key]; ok {
			if typ == nil {
				var val interface{}
				defaultDecode.Interface(&val, iter)
				d.kv[key] = val

			} else {
				decode := defaultDecode.Decode(typ)
				if decode == nil {
					iter.ReportError("ReadObject", "field("+key+") type("+typ.Name()+") not found decode ")
					return
				}

				val := reflect.New(typ).Elem()
				decode(val, iter)
				d.kv[key] = val.Interface()
			}

		} else {
			iter.Skip()
		}
		if iter.Error != nil {
			return
		}
	}

	if c != '}' {
		iter.ReportError("struct Decode", `expect }, but found `+string([]byte{c}))
	}

	return
}

type decoder func(val reflect.Value, iter *jsoniter.Iterator)

var ummarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

type jsonDecode struct {
}

var (
	defaultDecode = &jsonDecode{}
)

func (jd *jsonDecode) Decode(typ reflect.Type) decoder {

	// custom json.Unmarshal
	if typ.Implements(ummarshalerType) {
		return defaultDecode.Unmarshal
	}

	switch typ.Kind().String() {

	case "string":
		return jd.String
	//     "invalid"
	case "bool":
		return jd.Bool
	case "int":
		return jd.Int
	case "int8":
		return jd.Int8
	case "int16":
		return jd.Int16
	case "int32":
		return jd.Int32
	case "int64":
		return jd.Int64
	case "uint":
		return jd.Uint
	case "uint8":
		return jd.Uint8
	case "uint16":
		return jd.Uint16
	case "uint32":
		return jd.Uint32
	case "uint64":
		return jd.Uint64
	//"uintptr"
	case "float32":
		return jd.Float32
	case "float64":
		return jd.Float64
	//"complex64"
	//"complex128"
	case "array":
		return jd.Array
	//"chan"
	//"func"
	case "map":
		return jd.Map
	case "ptr":
		return jd.Pointer
	case "slice":
		return jd.Array
	case "struct":
		return jd.Struct
		//"unsafe.Pointer"
	}

	return nil
}

func (jd *jsonDecode) String(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadString()
	val.Set(reflect.ValueOf(v))
}

func (jd *jsonDecode) Int64(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadInt64()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Int32(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadInt32()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Int16(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadInt16()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Int8(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadInt8()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Int(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadInt()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Uint64(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadUint64()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Uint32(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadUint32()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Uint16(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadUint16()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Uint8(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadUint8()
	val.Set(reflect.ValueOf(v))

}

func (jd *jsonDecode) Uint(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadUint()
	val.Set(reflect.ValueOf(v))
}

func (jd *jsonDecode) Bool(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadBool()
	val.Set(reflect.ValueOf(v))
}

func (jd *jsonDecode) Float32(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadFloat32()
	val.Set(reflect.ValueOf(v))
}

func (jd *jsonDecode) Float64(val reflect.Value, iter *jsoniter.Iterator) {
	v := iter.ReadFloat64()
	val.Set(reflect.ValueOf(v))
}

func (jd *jsonDecode) Array(val reflect.Value, iter *jsoniter.Iterator) {
	rawTyp := val.Type()
	elemTyp := rawTyp.Elem()

	c := iter.NextToken()
	if c == 'n' {
		iter.SkipThreeBytes('u', 'l', 'l')
		return
	}
	if c != '[' {
		iter.ReportError("decode slice", "expect [ or n, but found "+string([]byte{c}))
		return
	}
	c = iter.NextToken()
	if c == ']' {
		return
	}
	iter.UnreadByte()

	if !iter.IncrementDepth() {
		return
	}
	defer iter.DecrementDepth()

	elemDecode := jd.Decode(elemTyp)
	if elemDecode == nil {
		iter.ReportError("decode slice elem", "not support type "+elemTyp.Kind().String())
		return
	}

	for c = ','; c == ','; c = iter.NextToken() {
		elemVal := reflect.New(elemTyp).Elem()
		elemDecode(elemVal, iter)
		if iter.Error != nil {
			return
		}
		val.Set(reflect.Append(val, elemVal))
	}
	if c != ']' {
		iter.ReportError("decode slice", "expect ], but found "+string([]byte{c}))
		return
	}

}

func (jd *jsonDecode) Map(val reflect.Value, iter *jsoniter.Iterator) {

	valTyp := val.Type().Elem()
	keyTyp := val.Type().Key()

	if jd.readStartObject(iter) {
		return
	}

	val.Set(reflect.MakeMap(val.Type()))

	keyDecode := jd.DecodeMapKey(keyTyp)
	// 非字符串的时候，会解析失败
	if keyDecode == nil {
		iter.ReportError("decode map elem", "not support key type "+keyTyp.Kind().String())
		return
	}

	valDecode := jd.Decode(valTyp)
	if valDecode == nil {
		iter.ReportError("decode map elem", "not support value type "+valTyp.Kind().String())
		return
	}

	iter.UnreadByte()

	var c byte
	for c = ','; c == ','; c = iter.NextToken() {
		elemKey := reflect.New(keyTyp).Elem()
		keyDecode(elemKey, iter)
		if iter.Error != nil {
			return
		}
		c = iter.NextToken()
		if c != ':' {
			iter.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
			return
		}
		elemVal := reflect.New(valTyp).Elem()
		valDecode(elemVal, iter)
		if iter.Error != nil {
			return
		}
		val.SetMapIndex(elemKey, elemVal)
	}
	if c != '}' {
		iter.ReportError("ReadMapCB", `expect }, but found `+string([]byte{c}))
	}

}

func (jd *jsonDecode) Pointer(val reflect.Value, iter *jsoniter.Iterator) {

	var elemVal reflect.Value
	if val.IsNil() {
		elemVal = reflect.New(val.Type().Elem()).Elem()
	} else {
		elemVal = val.Elem()
	}

	decode := jd.Decode(elemVal.Type())
	decode(elemVal, iter)
	val.Set(elemVal.Addr())

}

type structField struct {
	field   string
	name    string
	pkgName string
	index   int
	typ     reflect.Type
}

func describeStruct(structType reflect.Type) map[string]structField {
	fieldTypeMap := make(map[string]structField, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag, hastag := field.Tag.Lookup("json")
		if !hastag && !field.Anonymous {
			continue
		}
		if tag == "-" || field.Name == "_" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if field.Anonymous && (tag == "" || tagParts[0] == "") {
			if field.Type.Kind() == reflect.Struct {
				subFieldMap := describeStruct(field.Type)
				for k, t := range subFieldMap {
					fieldTypeMap[k] = t
				}
				continue
			} else if field.Type.Kind() == reflect.Ptr {
				if field.Type.Elem().Kind() == reflect.Struct {
					subFieldMap := describeStruct(field.Type.Elem())
					for k, t := range subFieldMap {
						fieldTypeMap[k] = t
					}
					continue
				}
			}
		}

		fieldTypeMap[tag] = structField{
			field:   tag,
			name:    field.Name,
			pkgName: field.PkgPath,
			index:   i,
			typ:     field.Type,
		}
	}

	return fieldTypeMap
}

func (jd *jsonDecode) Struct(val reflect.Value, iter *jsoniter.Iterator) {

	if jd.readStartObject(iter) {
		return
	}

	fieldTypeMap := describeStruct(val.Type())

	iter.UnreadByte()
	var c byte
	for c = ','; c == ','; c = iter.NextToken() {
		key := iter.ReadString()
		if iter.Error != nil {
			return
		}

		c = iter.NextToken()
		if c != ':' {
			iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			return
		}
		if field, ok := fieldTypeMap[key]; ok {
			typ := field.typ
			decode := jd.Decode(typ)
			if decode == nil {
				iter.ReportError("ReadObject", "field("+field.pkgName+") type("+typ.Name()+") not found decode ")
				return
			}
			fieldVal := reflect.New(typ).Elem()
			decode(fieldVal, iter)
			val.Field(field.index).Set(fieldVal)

		} else {
			iter.Skip()
		}
		if iter.Error != nil {
			return
		}
	}

	if c != '}' {
		iter.ReportError("struct Decode", `expect }, but found `+string([]byte{c}))
	}
	return

}

func (jd *jsonDecode) Interface(val interface{}, iter *jsoniter.Iterator) {
	jsonVal := iter.Read()
	valOf := reflect.ValueOf(val).Elem()
	valOf.Set(reflect.ValueOf(jsonVal))
}

// TODO: optimziation. not use **ptr
func (jd *jsonDecode) Unmarshal(val reflect.Value, iter *jsoniter.Iterator) {

	elemVal := reflect.New(val.Type().Elem())

	unmarshaler := elemVal.Interface().(json.Unmarshaler)
	iter.NextToken()
	iter.UnreadByte() // skip spaces
	bytes := iter.SkipAndReturnBytes()
	if iter.Error != nil {
		return
	}

	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		iter.ReportError("unmarshalerDecoder", err.Error())
	}

	val.Set(elemVal)
}

// only map key
func (jd *jsonDecode) DecodeMapKey(typ reflect.Type) decoder {

	decode := jd.Decode(typ)
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		f := func(val reflect.Value, iter *jsoniter.Iterator) {
			c := iter.NextToken()
			if c != '"' {
				iter.ReportError("ReadMapCB", `expect ", but found `+string([]byte{c}))
				return
			}
			decode(val, iter)
			c = iter.NextToken()
			if c != '"' {
				iter.ReportError("ReadMapCB", `expect ", but found `+string([]byte{c}))
				return
			}
		}
		return f

	}

	return decode

}

// 处理byte stream中关于object对象， 校验是否合法object， over 表示是否需要后需处理
func (jd *jsonDecode) readStartObject(iter *jsoniter.Iterator) (over bool) {

	if !iter.IncrementDepth() {
		return true
	}
	defer iter.DecrementDepth()

	c := iter.NextToken()
	if c == 'n' {
		iter.SkipThreeBytes('u', 'l', 'l')
		return true
	}

	if c != '{' {
		iter.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
		return true
	}
	c = iter.NextToken()
	if c == '}' {
		return true
	}

	return false
}
