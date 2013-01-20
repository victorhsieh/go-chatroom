package main

import (
    "./chatroom"
    "fmt"
)
import zmq "github.com/alecthomas/gozmq"

func main() {
  context, _ := zmq.NewContext()

  sub, err := context.NewSocket(zmq.SUB)
  if err != nil {
      panic(err)
  }
  sub.Connect(chatroom.PUBLISHER_ADDRESS)
  sub.SetSockOptString(zmq.SUBSCRIBE, "#pttdev")

  for {
    msg, _ := sub.Recv(0)
    fmt.Println(string(msg))
  }
}
