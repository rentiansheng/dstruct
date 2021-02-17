package dstruct

import (
	"reflect"
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
