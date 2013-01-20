package main

import (
    "./chatroom"
    "flag"
    "fmt"
    "math/rand"
    "os"
    "time"
)
import zmq "github.com/alecthomas/gozmq"

func usage() {
    fmt.Println("Usage: ./bot nickname channel")
    flag.PrintDefaults()
    os.Exit(2)
}

func random_quote() string {
    quotes := [...]string{ "Hello world", "/me hacking", "/kb in2" }
    return quotes[rand.Intn(len(quotes))]
}

func main() {
    flag.Usage = usage
    flag.Parse()
    args := flag.Args()
    if len(args) != 2 {
        usage()
    }
    nickname, channel := args[0], args[1]

    context, _ := zmq.NewContext()
    requester, err := context.NewSocket(zmq.PUSH)
    if err != nil {
        panic(err)
    }
    requester.Connect(chatroom.POST_ADDRESS)

    for {
        msg := random_quote()
        packet := fmt.Sprintf("%s %s %s", channel, nickname, msg)
        requester.Send([]byte(packet), 0)
        time.Sleep(100 * time.Millisecond)
    }
}
