// interface用来做类型转换和类的定义

package main

import (
  "fmt"
)

type Interface1 interface {
  GetName() string
}
type Interface2 interface {
  GetName() string
  GetAge() int
}
type Struct1 {
  name string
}
type Struct2 {
  name string
  age int
}

func (s *Struct1) GetName() string {
  return s.name
}
func (s *Struct2) GetName() string {
  return s.name
}
func (s *Struct2) GetAge() int {
  return s.age
}

func InterfaceConvert() {
  var i1 Interface1 = &Struct1{name: "struct1"}
  var i2 Interface2 = &Struct1{name: "struct2", age: 20}

  // struct1
  fmt.Println(i1.GetName())
  // struct2
  fmt.Println(i2.GetName())
  // 20
  fmt.Println(i2.GetAge())

  // i2定义成Struct2类型也可以
  var i3 Interface1 = i2
  //struct2
  fmt.Println(i3.GetName())

  // s3定义为struct就需要interface转换为对应的类型, 且不能转换成Struct1
  // var s3 Struct2 = *i2.(*Struct2)
}
func main() {
    var anyVal interface{} = "hello"
    if s, ok := anyVal.(string); ok {
        fmt.Println("转换成功: ", s)
    } else {
        fmt.Println("转换失败")
    }
  
    //int 转换成interface可以成功，int32不行
    var valInt32 int32 = 3
    var val int = 3
    var val1 interface{} = valInt32
    var val2 interface{} = val
    if i, ok := val.(int); ok {
      fmt.Println("转换成功: ", i)
    } else {
        fmt.Println("转换失败")
    }
}
