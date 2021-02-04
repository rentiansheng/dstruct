package dstruct

import (
	"fmt"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

type test struct {
	Str string `json:"str"`
	Int int    `json:"int"`
}

func TestUnmarshalJSON(t *testing.T) {

	/*aa := saa{}

	bb := `{"dd":[1,2], "aa":"1"}`
	jsoniter.Unmarshal([]byte(bb), &aa)*/
	bb := `{"dd":1, "aa":"str", "bl":true, "arr":["1","2"], "map":{"a":1,"b":64}, "test_struct":{"str":"string", "int":1}}`
	aa := &DStruct{}
	var i interface{}
	aa.SetFields(map[string]reflect.Type{
		"dd":  reflect.TypeOf(int64(1)),
		"aa":  reflect.TypeOf(string("1")),
		"bl":  reflect.TypeOf(bool(true)),
		"arr": reflect.TypeOf([]string{}),
		"map": reflect.TypeOf(map[string]int64{}),
		//"test_struct": reflect.TypeOf(test{}),
		"test_struct": reflect.TypeOf(i),
	})
	jsoniter.Unmarshal([]byte(bb), &aa)
	fmt.Println(aa.String())

}

func TestUnmarshalStrJSON(t *testing.T) {

	/*aa := saa{}

	bb := `{"dd":[1,2], "aa":"1"}`
	jsoniter.Unmarshal([]byte(bb), &aa)*/
	bb := `{"dd":1, "aa":"str", "bl":true}`
	aa := &DStruct{}
	aa.SetFields(map[string]reflect.Type{"aa": reflect.TypeOf(string("1"))})
	jsoniter.Unmarshal([]byte(bb), &aa)

	fmt.Println(aa.String())

}

func TestUnmarshalJSON1(t *testing.T) {

	aa := make(map[string]interface{}, 0) //saa{}

	bb := `{"dd":1, "aa":"1"}`
	jsoniter.Unmarshal([]byte(bb), &aa)

}

type saa struct {
	DD  []int64 `json:"dd"`
	AA  *string `json:"aa"`
	DD1 []int64 `json:"dd1"`
	AA1 string  `json:"aa1"`
	DD2 []int64 `json:"dd2"`
	AA2 string  `json:"aa2"`
	DD3 []int64 `json:"dd3"`
	AA3 string  `json:"aa3"`
	DD4 []int64 `json:"dd4"`
	AA4 string  `json:"aa4"`
	DD5 []int64 `json:"dd5"`
	AA5 string  `json:"aa5"`
}

func (s saa) DynamicStruct() {

}
