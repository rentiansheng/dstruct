### golang 动态结构体


#### 介绍

 用户可以自定字段的具体类型, 预计支持json， bson 等方式



####  目前开发进度

json Unmarshal 主要功能开发完成，指针类型暂不支持

bson Unmarshal 暂未启动


##### json Unmarshal

使用  github.com/json-iterator/go 来 Iterator实现数据读取和接卸，后需计划使用unsafe.Pointer ， 减少reflect性能损耗


