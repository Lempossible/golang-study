// this is a test of using reflect

package main

import (
  "fmt"
  "reflect"
)

func main() {
  var num = 42
  var val = reflect.ValueOf(num)
  // output is "value of num is 42". reflect.ValueOf
  fmt.Println("value of num is ", val)
  // output is "type of num is int". val.Type
  fmt.Println("type of num is ", val.Type())

  var t = val.Type()
  switch t.Kind() {
  case reflect.Int:
    fmt.Println("type is int")
  case reflect.String:
    fmt.Println("type is string")
  case reflect.Array:
    fmt.Println("type is array")
  default:
      fmt.Println("unknown type")
  }
}
