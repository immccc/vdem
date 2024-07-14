package node

import (
	"fmt"
	http "net/http"
	"os"
	"log"
)

type Node struct {
	Config NodeConfig
	Peers []Peer
}

func (node *Node) Start() {
	if len(node.Peers) == 0 {
		log.Println("Node has started without peers, won't be listening for udpates!")
	}
		

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%v", node.Config.ServerPort), nil)

	if err != nil {
		fmt.Printf("Error on server: %s\n", err)
		os.Exit(1)
	}
}

func (node *Node) AddPeer(peer Peer) {
	if node.Peers == nil {
		node.Peers = make([]Peer, 1)
	}

	node.Peers = append(node.Peers, peer)

}

func serveWs(w http.ResponseWriter, r *http.Request) {

}
