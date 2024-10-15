package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	websocket2 "prettyy-server-online/utils/websocket"
	"sync"
)

var mutex = &sync.Mutex{}

func main() {
	for {
		msg := <-websocket2.Message
		fmt.Println("message-loader:", msg)
		mutex.Lock()
		for client := range websocket2.Clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(websocket2.Clients, client)
			}
		}
		mutex.Unlock()
	}
}
