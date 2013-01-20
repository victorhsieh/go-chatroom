package main

import (
    "./chatroom"
    "fmt"
)
import zmq "github.com/alecthomas/gozmq"

func admin_server(c chan int) {
    // TODO maybe create another zmq channel
    var seq int64 = 0
    for {
        num := <-c
        if num == 1 {
            seq += int64(num)
        } else {
            fmt.Println("[DEBUG] stat: total message:", seq)
        }
    }
}

func bind_to_channel(sock zmq.Socket) (channel chan []byte) {
    channel = make(chan []byte)
    go func(){
        for {
            msg, err := sock.Recv(0)
            if err != nil {
                fmt.Println("[ERROR] die at sock.Recv:", err.Error())
                break
            }
            channel <- msg
        }
    }()
    return
}

func main() {
  context, _ := zmq.NewContext()
  publisher, _ := context.NewSocket(zmq.PUB)
  defer publisher.Close()
  publisher.Bind(chatroom.PUBLISHER_ADDRESS)

  router, _ := context.NewSocket(zmq.PULL)
  defer router.Close()
  router.Bind(chatroom.POST_ADDRESS)

  stat_channel := make(chan int)
  go admin_server(stat_channel)

  receiver_channel := bind_to_channel(router)

  for {
    select {
    case msg := <-receiver_channel:
        go func() {
            stat_channel <- 1
            publisher.Send(msg, 0)
        }()
    }
  }
}
