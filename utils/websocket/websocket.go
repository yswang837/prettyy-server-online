package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[*websocket.Conn]bool)
var Message = make(chan string, 10)
