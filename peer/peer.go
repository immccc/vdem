package peer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// A Peer acts as a Nostr client. Clients are part of nodes and communicates with them.
type Peer struct {
	PubKey string
	Port   int
	Host   string
	conn   *websocket.Conn
}

func (peer *Peer) ToURL() string {
	return fmt.Sprintf("ws://%v:%v", peer.Host, peer.Port)
}

func New(pubKey string, url string) Peer {
	hostAndPort := strings.Split(url, ":")

	host := hostAndPort[0]
	if len(host) == 0 {
		host = "localhost"
	}

	portStr := hostAndPort[1]
	port, _ := strconv.Atoi(portStr)

	return Peer{
		PubKey: pubKey,
		Host:   host,
		Port:   port,
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
	conn.WriteMessage(websocket.TextMessage, msg[:])
}

func (peer *Peer) Close() {
	if peer.conn == nil {
		return
	}

	peer.conn.Close()
}
