package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

func main() {
    sock, err := websocket.Dial("ws://localhost:5050", "", "http://localhost/")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

	for {
		var raw string
		err := websocket.Message.Receive(sock, &raw)
		if err != nil {
            fmt.Println(err.Error())
            return
		}

        fmt.Println(raw)
	}
}
