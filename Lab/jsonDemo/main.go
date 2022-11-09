package main

//基本的序列化****************************************************
import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name   string  //`json:"mingzi"`
	Age    int64   //`json:"-"`
	Weight float64 //`json:"zhongliang,omitempty"`
}

func main() {
	p1 := Person{
		Name: "wangkun",
		Age:  23,
	}
	// struct -> json string
	b, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
	// json string -> struct
	var p2 Person
	err = json.Unmarshal(b, &p2)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("p2:%#v\n", p2)
}

//`key1:"value1" key2:"value2"` 规则
//使用json tag指定字段名：可以通过给结构体字段添加tag来指定json序列化生成的字段名
//使用json tag忽略某个字段： `json:"-"`
//忽略空值字段：当 struct 中的字段没有值时， json.Marshal() 序列化的时候不会忽略这些字段，而是默认输出字段的类型零值，如果想要在序列序列化时忽略这些没有值的字段时，可以在对应字段添加omitempty tag。
//忽略嵌套的结构体空值字段：嵌套加omitempty
//不修改结构体忽略空值字段
//优雅处理字符串格式的数字
