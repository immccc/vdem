package node

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Listener struct {
	writer  http.ResponseWriter
	request *http.Request

	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (listener *Listener) Start() {
	conn, _ := upgrader.Upgrade(listener.writer, listener.request, nil)
	listener.conn = conn

	for {
		_, msg, err := listener.conn.ReadMessage()
		if err != nil {
			log.Printf("Err receiving value: %v", err)
		}

		log.Printf("MSG Received %v", msg)

	}

}
