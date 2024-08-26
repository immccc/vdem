package node

import (
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"log"
	http "net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Node struct {
	Config NodeConfig
	Peers  []peer.Peer
}

func (node *Node) Start(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	if node.Peers == nil || len(node.Peers) == 0 {
		log.Println("Node has started without peers, wont be listening for udpates!")
	}

	mux := http.NewServeMux()
	serverStarted := false

	mux.HandleFunc("/health/", func(_ http.ResponseWriter, __ *http.Request) {
		serverStarted = true
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		node.serveWs(w, r)
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", node.Config.ServerPort), mux)

		if err != nil {
			log.Fatalf("Error on server: %v", err)
		}
	}()

	for !serverStarted {
		//TODO https if applicable
		http.Get(fmt.Sprintf("http://localhost:%v/health", node.Config.ServerPort))
		time.Sleep(time.Second)
	}

	if node.Peers != nil && len(node.Peers) > 0 {
		//TODO Add host in the config
		msg := messaging.BuildConnectionAttemptMsg(node.Config.PubKey, "", node.Config.ServerPort)
		msg.Sign(node.Config.PrivateKey)
		node.Send(&msg, &node.Peers[0])
	}
}

func (node *Node) AddPeer(pr peer.Peer) {
	if node.Peers == nil {
		node.Peers = make([]peer.Peer, 0)
	}

	node.Peers = append(node.Peers, pr)
}

func (node *Node) Send(msg *messaging.Message, peer *peer.Peer) {
	conn, resp, err := websocket.DefaultDialer.Dial(peer.ToURL(), nil)
	if err != nil {
		log.Printf("Cannot send message! %v to peer %v. Error is %v", msg, peer.ToURL(), err)
		log.Printf("Resp: %v", resp)
	}
	defer conn.Close()

	msg.Sign(node.Config.PrivateKey)
	conn.WriteJSON(msg)

}

func (node *Node) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade connection error: ", err)
		return
	}

	defer conn.Close()

	for {
		message := messaging.Message{}

		log.Println("Reading...")

		err := conn.ReadJSON(&message)

		if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			log.Println("Error on message read: ", err)
			break
		}

		if !message.Verify() {
			log.Println("ACHTUNG! Message cannot be verified!")
			continue
		}

		log.Println("Received message: ", message) // TODO Debug level to logs!
		ActionsPerMessage[message.Kind](node, &message, r)
	}
}

func (node *Node) Connect(msg *messaging.Message, request *http.Request) {
	if !node.Config.ForceConnectionRequests {
		log.Printf("{Host: %v, connection: DENIED, reason: Only connections are accepted when ForceConnectionRequests is True }", request.Host)
		return
	}

	newPeer := peer.FromHost(request.Host)
	node.AddPeer(newPeer)

}
