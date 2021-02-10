package dstruct

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

/****** 辅助测试定义 ******/

type testStruct struct {
	Str string `json:"str"`
	Int int    `json:"int"`
}

type kv map[string]interface{}

func (k *kv) UnmarshalJSON(data []byte) error {
	kv := make(map[string]interface{}, 0)
	err := json.Unmarshal(data, &kv)
	if err != nil {
		return err
	}
	*k = kv
	return nil
}

/****** 测试用例 ******/

// 综合测试
func TestUnmarshalJSON(t *testing.T) {

	testStrJSONBytes := `{"int":1,"str":"str","bl":true, "arr_str":["1","2"],"arr_int":[1,2], "map":{"a":1,"b":64}, ` +
		`"test_struct":{"str":"string", "int":1},"interface":{"str":"string", "bl":true}}`
	ds := &DStruct{}
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"int":         TypeInt,
		"str":         TypeString,
		"bl":          TypeBool,
		"arr_str":     TypeArrayStr,
		"arr_int":     TypeArrayInt,
		"map":         reflect.TypeOf(map[string]int64{}),
		"test_struct": reflect.TypeOf(testStruct{}),
		"interface":   reflect.TypeOf(i),
	})
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Error("unmarshal: " + err.Error())
		return
	}

	require.Equal(t, ds.kv["int"], 1, "test int.")
	require.Equal(t, ds.kv["str"], "str", "test str.")
	require.Equal(t, ds.kv["bl"], true, "test bool.")
	require.Equal(t, ds.kv["arr_str"], []string{"1", "2"}, "test []string.")
	require.Equal(t, ds.kv["arr_int"], []int{1, 2}, "test []int.")
	require.Equal(t, ds.kv["map"], map[string]int64{"a": 1, "b": 64}, "test map[string]int64.")
	require.Equal(t, ds.kv["test_struct"], testStruct{Str: "string", Int: 1}, "test struct.")
	require.Equal(t, ds.kv["interface"], map[string]interface{}{"str": "string", "bl": true}, "test interface{}.")

}

func TestMapUnmarshalJSON(t *testing.T) {
	dataMap := map[string]map[int]map[string]testStruct{
		"l1-1": map[int]map[string]testStruct{
			1: map[string]testStruct{
				"l1-3": testStruct{Str: "l1-3 string", Int: 1},
			},
		},
		"l2-1": map[int]map[string]testStruct{
			2: map[string]testStruct{
				"l2-1": testStruct{Str: "l2-3 string", Int: 1},
			},
		},
		"l3-1": map[int]map[string]testStruct{
			2: nil,
		},
	}
	jsonWrap := map[string]interface{}{"test_field": dataMap}
	testStrJSONBytes, err := json.Marshal(jsonWrap)
	if err != nil {
		t.Error(err)
		return
	}

	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"test_field": reflect.TypeOf(dataMap),
	})
	err = json.Unmarshal(testStrJSONBytes, ds)
	if err != nil {
		t.Error("unmarshal: " + err.Error())
		return
	}

	rawData := make(map[string]map[int]map[string]testStruct, 0)
	ok, err := ds.Value("test_field", &rawData)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !ok {
		t.Errorf("not found test field")
		return
	}

	require.Equal(t, dataMap, rawData, "test many level map.")
}

func TestCustomUnmarshalJSON(t *testing.T) {

	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"test_unmarshaljson": reflect.TypeOf(&kv{}),
	})

	testStrJSONBytes := []byte(`{"test_unmarshaljson":{"k1":"v1", "k2":"k2"}}`)
	err := json.Unmarshal(testStrJSONBytes, ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	dataMap := kv{
		"k1": "v1",
		"k2": "k2",
	}
	rawData := &kv{}
	// 这里需要取两次地址
	ok, err := ds.Value("test_unmarshaljson", &rawData)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !ok {
		t.Errorf("not found test field")
		return
	}

	require.Equal(t, dataMap, *rawData, "test user custom UnmarshalJSON map.")
}

func TestPtrUnmarshalJSON(t *testing.T) {
	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"ptr": reflect.TypeOf(new(testStruct)),
	})

	testStrJSONBytes := []byte(`{"ptr":{"str":"string", "int":1}}`)
	err := json.Unmarshal(testStrJSONBytes, ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	rawData := &testStruct{
		Str: "string",
		Int: 1,
	}
	require.Equal(t, ds.kv["ptr"], rawData, "test ptr.")

}

func TestNone(t *testing.T) {

	testStrJSONBytes := `{}`
	ds := &DStruct{}
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"int64":         reflect.TypeOf(int64(1)),
		"string":        reflect.TypeOf(string("1")),
		"bl":            reflect.TypeOf(bool(true)),
		"arr_str":       reflect.TypeOf([]string{}),
		"map_str_int64": reflect.TypeOf(map[string]int64{}),
		"test_struct":   reflect.TypeOf(testStruct{}),
		"interface":     reflect.TypeOf(i),
	})
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	require.Equal(t, ds.kv["int64"], int64(0), "test int default.")
	require.Equal(t, ds.kv["string"], "", "test string default.")
	require.Equal(t, ds.kv["bl"], false, "test bool default.")
	require.Equal(t, ds.kv["arr_str"], ([]string)(nil), "test []string default.")
	require.Equal(t, ds.kv["test_struct"], testStruct{}, "test struct default.")
	require.Equal(t, ds.kv["interface"], nil, "test interface default.")

}

func TestInterface(t *testing.T) {

	testStrJSONBytes := []byte(`{"str":"str","bl":true}`)
	ds := &DStruct{}
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"str": reflect.TypeOf(i),
		"bl":  reflect.TypeOf(i),
	})
	err := json.Unmarshal(testStrJSONBytes, ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	require.Equal(t, ds.kv["str"], "str", "test string interface{}.")

}

func TestMapInterface(t *testing.T) {

	testStrJSONBytes := `{"map":{"str":"str","bl":true}, "arr":["1", "2"],` +
		` "arr_map":[{"str":"str"}],"arr_arr":[["str1","str2"]]}`
	ds := &DStruct{}
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"map":     reflect.TypeOf(i),
		"arr":     reflect.TypeOf(i),
		"arr_map": reflect.TypeOf(i),
		"arr_arr": reflect.TypeOf(i),
	})
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	require.Equal(t, ds.kv["map"], map[string]interface{}{"str": "str", "bl": true}, "test map interface{}.")
	require.Equal(t, ds.kv["arr"], []interface{}{"1", "2"}, "test map interface{}.")
	require.Equal(t, ds.kv["arr_map"], []interface{}{map[string]interface{}{"str": "str"}}, "test arr_map interface{}.")
	require.Equal(t, ds.kv["arr_arr"], []interface{}{[]interface{}{"str1", "str2"}}, "test arr_arr interface{}.")

}

func TestJSONNumber(t *testing.T) {
	testStrJSONBytes := `{"int":10000, "arr":[1]}`
	ds := &DStruct{}
	var i interface{}
	ds.SetFields(map[string]reflect.Type{
		"int": reflect.TypeOf(i),
		"arr": reflect.TypeOf(i),
	})
	ds.JSONNumber()
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	require.Equal(t, ds.kv["int"], (json.Number)("10000"), "test json number interface{}.")
	require.Equal(t, ds.kv["arr"], []interface{}{(json.Number)("1")}, "test json number  arr interface{}.")

	var intVal interface{}
	ok, err := ds.Value("int", &intVal)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !ok {
		t.Errorf("not found test int field")
		return
	}

	intJSONNumberVal, ok := intVal.(json.Number)
	if !ok {
		t.Errorf("test int field. not json.Unmber")
		return
	}
	require.Equal(t, intJSONNumberVal, (json.Number)("10000"), "test json number interface{}.")

}

func TestAliasType(t *testing.T) {
	type MyInt int64
	type MyStr string
	testStrJSONBytes := `{"int":10000, "arr":[1], "str":"str"}`
	ds := &DStruct{}
	ds.SetFields(map[string]reflect.Type{
		"int": reflect.TypeOf(MyInt(1)),
		"arr": reflect.TypeOf([]MyInt{}),
		"str": reflect.TypeOf(string("")),
	})
	ds.JSONNumber()
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	require.Equal(t, ds.kv["int"], (MyInt)(10000), "test alias int .")
	require.Equal(t, ds.kv["arr"], []MyInt{1}, "test alias int array  arr.")
	require.Equal(t, ds.kv["str"], "str", "test alias string .")

	var intVal MyInt
	ok, err := ds.Value("int", &intVal)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !ok {
		t.Errorf("not found test int field")
		return
	}

	require.Equal(t, intVal, (MyInt)(10000), "test json number interface{}.")

}
