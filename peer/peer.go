package peer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Peer struct {
	Port int
	Host string
	conn *websocket.Conn
}

func (peer *Peer) ToURL() string {
	return fmt.Sprintf("ws://%v:%v", peer.Host, peer.Port)
}

func FromHost(url string) Peer {
	hostAndPort := strings.Split(url, ":")

	host := hostAndPort[0]
	portStr := hostAndPort[1]
	port, _ := strconv.Atoi(portStr)

	return Peer{
		Host: host,
		Port: port,
	}
}

func (peer *Peer) getConnection() *websocket.Conn {
	if peer.conn == nil {
		conn, resp, err := websocket.DefaultDialer.Dial(peer.ToURL(), nil)
		if err != nil {
			log.Printf("Cannot open connection to peer %v . Error is %v", peer.ToURL(), err)
			log.Printf("Resp: %v", resp)
		}
		peer.conn = conn
	}

	return peer.conn

}

func (peer *Peer) SendMessage(msg []byte) {
	conn := peer.getConnection()
	conn.WriteMessage(websocket.TextMessage, msg)
}

func (peer *Peer) Close() {
	if peer.conn == nil {
		return
	}

	peer.conn.Close()
}
