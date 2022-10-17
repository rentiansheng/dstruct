package dstruct

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func (d *DStruct) UnmarshalBSON(data []byte) error {

	d.init()

	elems, err := bsoncore.Document(data).Elements()
	if err != nil {
		return err
	}

	for _, elem := range elems {
		field := elem.Key()
		rawTye, ok := d.fields[field]
		if !ok {
			continue
		}
		dc := bsoncodec.DecodeContext{Registry: bson.DefaultRegistry}
		vr := bsonrw.NewBSONValueReader(elem.Value().Type, elem.Value().Data)
		decoder, err := bson.NewDecoderWithContext(dc, vr)
		if err != nil {
			return err
		}
		val := reflect.New(rawTye)
		if err = decoder.Decode(val.Interface()); err != nil {
			return err
		}
		if val.IsNil() {
			d.kv[field] = val.Elem()
			continue
		}
		d.kv[field] = val.Elem().Interface()
		if err := validateStruct(val); err != nil {
			return fmt.Errorf("validator: ield(" + field + ") type(" + val.Type().Name() + ") " + err.Error())
		}
	}

	return nil
}

func (d *DStruct) MarshalBSON() ([]byte, error) {
	return bson.Marshal(d.kv)
}
