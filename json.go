package dstruct

import (
	"encoding/json"
	"reflect"
	"strings"
	"unsafe"

	"dstruct/jsoniter"
)

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(rtype unsafe.Pointer, m unsafe.Pointer, key, val unsafe.Pointer)

func (d *DStruct) UnmarshalJSON(data []byte) error {

	cfg := jsoniter.Config{
		EscapeHTML: true,
	}.Froze()

	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	d.jsonDecode(iter)
	return iter.Error

}

func (d *DStruct) init() {
	d.kv = make(map[string]interface{}, 0)
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

var marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

type jsonDecode struct {
}

var (
	defaultDecode = &jsonDecode{}
)

func (jd *jsonDecode) Decode(typ reflect.Type) decoder {

	// custom json.Unmarshal
	if typ.Implements(marshalerType) {
		return defaultDecode.Unmarshal
	}

	decodeMap := map[string]decoder{
		"string": jd.String,
		//     "invalid",
		"bool":   jd.Bool,
		"int":    jd.Int,
		"int8":   jd.Int8,
		"int16":  jd.Int16,
		"int32":  jd.Int32,
		"int64":  jd.Int64,
		"uint":   jd.Uint,
		"uint8":  jd.Uint8,
		"uint16": jd.Uint16,
		"uint32": jd.Uint32,
		"uint64": jd.Uint64,
		//"uintptr",
		"float32": jd.Float32,
		"float64": jd.Float64,
		//"complex64",
		//"complex128",
		"array": jd.Array,
		//"chan",
		//"func",
		"map": jd.Map,
		//"ptr",
		"slice":  jd.Array,
		"struct": jd.Struct,
		//"unsafe.Pointer",
	}

	return decodeMap[typ.Kind().String()]
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

	val.Set(reflect.MakeMap(val.Type()))

	if !iter.IncrementDepth() {
		return
	}
	defer iter.DecrementDepth()
	keyDecode := jd.Decode(keyTyp)
	if keyDecode == nil {
		iter.ReportError("decode map elem", "not support key type "+keyTyp.Kind().String())
		return
	}

	valDecode := jd.Decode(valTyp)
	if valDecode == nil {
		iter.ReportError("decode map elem", "not support value type "+valTyp.Kind().String())
		return
	}

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
	for c = iter.NextToken(); c == ','; c = iter.NextToken() {
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

	fieldTypeMap := describeStruct(val.Type())

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
	iter.ReadVal(val)
}

func (jd *jsonDecode) Unmarshal(val reflect.Value, iter *jsoniter.Iterator) {

	iter.ReadObject()

}
