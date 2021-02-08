### golang 动态结构体


#### 介绍

 用户可以自定字段的具体类型, 预计支持json， bson 等方式



####  目前开发进度

json Unmarshal 主要功能开发完成，指针类型暂不支持

bson Unmarshal 暂未启动


##### json Unmarshal

使用  github.com/json-iterator/go 来 Iterator实现数据读取和接卸，后需计划使用unsafe.Pointer ， 减少reflect性能损耗

eg:

```
    
	testStrJSONBytes := `{"int":1,"str":"str","bl":true, "arr_str":["1","2"],"arr_int":[1,2], "map":{"a":1,"b":64}, ` +
		`"test_struct":{"str":"string", "int":1},"interface":{"str":"string", "bl":true}}`
        // define a dynamic structure 
	ds := &DStruct{}
	var i interface{}
        // define the fields of the dynamic structure
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
        // Unmarshal dynamic structure, you can directly use golang official
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
        t.Error(err.Error())
		return
	}
    ds.Int64("int")
    ds.String("str")
    ds.Value("interface", &i)

```



