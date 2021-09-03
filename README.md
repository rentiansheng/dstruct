### golang 动态结构体


#### 介绍

 用户可以自定字段的具体类型, 预计支持json， bson 等方式



####  目前开发进度

json Unmarshal 主要功能开发完成，指针类型暂不支持

bson Unmarshal 暂未启动



#### 设计思路

struct 在我们使用的时候之所以方便，简洁就是在做反序列的时候根据结构的定义字段和类型。 因此在使用的可以不用在做类型判断和数据校验
如果可以做一个dynamic struct，这个struct字段是根据执行逻辑动态调整字段和字段类型即可。

实现定一个DStruct 结构， 给DStruct结构实现基于byte stream json 反序列， 利用golang可以自定义struct UnmarshalJSON的特性，
当做json Unmarshal调用DStruct 根据用户定义的配置项目反序列化。

#### json Unmarshal

使用  github.com/json-iterator/go 来 Iterator实现数据读取和数据反序列化，后需计划使用unsafe.Pointer ， 减少reflect性能损耗

eg:
```golang

 
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
		"interface":   reflect.TypeOf(i),
	})
	// Unmarshal dynamic structure, you can directly use golang official
	err := json.Unmarshal([]byte(testStrJSONBytes), ds)
	if err != nil {
		t.Error(err.Error())
		return
	}
        intVal, err := ds.Int64("int")
        if err != nil {
		fmt.Error(err.Error())
		return
	}
		fmt.Println("int test, err: ", err, "val: ", intVal)
	str, err := ds.Str("str")
	if err != nil {
		fmt.Error(err.Error())
		return
	}
	fmt.Println("str test, err: ", err, "val: ", str)
	exists, err := ds.Value("interface", &i)
	if err != nil {
		fmt.Error(err.Error())
		return
	}
	fmt.Println("exist field: ", exists, "val: ", i)
```



#### bson Unmarshal

eg:
```golang

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

```
