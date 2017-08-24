package main

import (
  "github.com/notegio/0xrelay/types"
  "io/ioutil"
  "encoding/json"
  // "encoding/hex"
  "fmt"
  // "reflect"
)

func main() {
  order := types.Order{}
  if orderData, err := ioutil.ReadFile("formatted_transaction.json"); err == nil {
    if err := json.Unmarshal(orderData, &order); err != nil {
      println(err.Error())
      return
    }
  } else {
    println(err.Error())
    return
  }
  println(order.Signature.Verify(order.Maker))
  ob := order.Bytes()
  fmt.Printf("'%v'\n", order.Signature.V)
  newOrder := types.OrderFromBytes(ob)
  println(newOrder.Signature.Verify(newOrder.Maker))
  // println(reflect.DeepEqual())
}
