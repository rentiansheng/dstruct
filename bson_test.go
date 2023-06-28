package dstruct

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUnmarshalBSON(t *testing.T) {
	testMap := map[string]interface{}{
		"int":     int(1),
		"str":     "str",
		"bl":      true,
		"arr_str": []string{"str1", "str2"},
		"arr_int": []int64{1, 2},
		"map": map[string]int64{
			"int1": int64(1),
		},
		"test_struct": testStruct{Str: "string", Int: 1},
	}

	bsonBytes, err := bson.Marshal(testMap)
	if err != nil {
		t.Error("marshal: " + err.Error())
		return
	}
	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"int":         TypeInt,
		"str":         TypeString,
		"bl":          TypeBool,
		"arr_str":     TypeArrayStr,
		"arr_int":     TypeArrayInt,
		"map":         reflect.TypeOf(map[string]int64{}),
		"test_struct": reflect.TypeOf(testStruct{}),
	})
	err = bson.Unmarshal(bsonBytes, ds)
	if err != nil {
		t.Error("unmarshal: " + err.Error())
		return
	}

	require.Equal(t, ds.kv["int"], 1, "test int.")
	require.Equal(t, ds.kv["str"], "str", "test str.")
	require.Equal(t, ds.kv["bl"], true, "test bool.")
	require.Equal(t, ds.kv["arr_str"], []string{"str1", "str2"}, "test []string.")
	require.Equal(t, ds.kv["arr_int"], []int{1, 2}, "test []int.")
	require.Equal(t, ds.kv["map"], map[string]int64{"int1": 1}, "test map[string]int64.")
	require.Equal(t, ds.kv["test_struct"], testStruct{Str: "string", Int: 1}, "test struct.")

}

func TestMarshalBSON(t *testing.T) {
	testMap := map[string]interface{}{
		"int":     int(1),
		"str":     "str",
		"bl":      true,
		"arr_str": []string{"str1", "str2"},
		"arr_int": []int64{1, 2},
		"map": map[string]int64{
			"int1": int64(1),
		},
		"test_struct": testStruct{Str: "string", Int: 1},
	}
	bsonBytes, err := bson.Marshal(testMap)
	if err != nil {
		t.Error("marshal: " + err.Error())
		return
	}
	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"int":         TypeInt,
		"str":         TypeString,
		"bl":          TypeBool,
		"arr_str":     TypeArrayStr,
		"arr_int":     TypeArrayInt,
		"map":         reflect.TypeOf(map[string]int64{}),
		"test_struct": reflect.TypeOf(testStruct{}),
	})
	err = bson.Unmarshal(bsonBytes, ds)
	if err != nil {
		t.Error("unmarshal: " + err.Error())
		return
	}

	bsonBytes, err = bson.Marshal(ds)
	if err != nil {
		t.Error("marshal: " + err.Error())
		return
	}

	newDS := ds.Clone()
	err = bson.Unmarshal(bsonBytes, newDS)
	if err != nil {
		t.Error("unmarshal: " + err.Error())
		return
	}

	require.Equal(t, newDS.kv["int"], 1, "test int.")
	require.Equal(t, newDS.kv["str"], "str", "test str.")
	require.Equal(t, newDS.kv["bl"], true, "test bool.")
	require.Equal(t, newDS.kv["arr_str"], []string{"str1", "str2"}, "test []string.")
	require.Equal(t, newDS.kv["arr_int"], []int{1, 2}, "test []int.")
	require.Equal(t, newDS.kv["map"], map[string]int64{"int1": 1}, "test map[string]int64.")
	require.Equal(t, newDS.kv["test_struct"], testStruct{Str: "string", Int: 1}, "test struct.")
}

func TestValidatorUnmarshalBSON(t *testing.T) {
	suits := []struct {
		input  string
		hasErr bool
	}{
		{
			input: `{"int":4,"str":"str","test_struct":{"str":"123", "int":4}}`,
		},
		{
			input: `{"int":4,"str":"str","test_struct":{"str":"1234567890", "int":10}}`,
		},
		{
			input: `{"int":4,"str":"str","test_struct":{"str":"12345", "int":5}}`,
		},
		// str validator
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"", "int":4}}`,
			hasErr: true,
		},
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"123456789012121", "int":4}}`,
			hasErr: true,
		},
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"12", "int":4}}`,
			hasErr: true,
		},

		// int validator
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"123", "int":1}}`,
			hasErr: true,
		},
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"123", "int":11}}`,
			hasErr: true,
		},
		{
			input:  `{"int":4,"str":"str","test_struct":{"str":"123", "int":-1}}`,
			hasErr: true,
		},
	}

	ds := &DStruct{}
	ds.ValidateOn()
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"int":         TypeInt,
		"str":         TypeString,
		"bl":          TypeBool,
		"arr_str":     TypeArrayStr,
		"arr_int":     TypeArrayInt,
		"map":         reflect.TypeOf(map[string]int64{}),
		"test_struct": reflect.TypeOf(testValidateStruct{}),
		"interface":   reflect.TypeOf(i),
	})
	for idx, suit := range suits {
		tmpMap := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(suit.input), &tmpMap)
		if err != nil {
			t.Errorf("unmarshal index %d error %s ", idx, err.Error())
			continue
		}
		testBsonBytes, err := bson.Marshal(tmpMap)
		if err != nil {
			t.Errorf("unmarshal index %d error %s ", idx, err.Error())
			continue
		}

		err = bson.Unmarshal(testBsonBytes, ds)
		if suit.hasErr {
			if err == nil {
				t.Errorf("unmarshal index %d validator not work ", idx)
			} else if !strings.HasPrefix(err.Error(), "validator: field(test_struct) type(testValidateStruct) Key:") {
				t.Errorf("unmarshal index %d error %s ", idx, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("unmarshal index %d error %s ", idx, err.Error())
			}
		}

	}
}
