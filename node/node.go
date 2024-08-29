package node

import (
	"encoding/json"
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

var upgrader = websocket.Upgrader{}

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
		event := messaging.BuildConnectionAttemptEvent(node.Config.PubKey, "", node.Config.ServerPort)
		event.Sign(node.Config.PrivateKey)
		node.Send(&event, &node.Peers[0])
	}
}

func (node *Node) AddPeer(pr peer.Peer) {
	if node.Peers == nil {
		node.Peers = make([]peer.Peer, 0)
	}

	node.Peers = append(node.Peers, pr)
}

func (node *Node) Send(event *messaging.Event, peer *peer.Peer) {
	//TODO Connections to last until events are answered!
	conn, resp, err := websocket.DefaultDialer.Dial(peer.ToURL(), nil)
	if err != nil {
		log.Printf("Cannot send event! %v to peer %v. Error is %v", event, peer.ToURL(), err)
		log.Printf("Resp: %v", resp)
	}
	defer conn.Close()

	event.Sign(node.Config.PrivateKey)

	conn.WriteMessage(websocket.TextMessage, messaging.BuildEventMessage(event))

}

func (node *Node) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade connection error: ", err)
		return
	}

	defer conn.Close()

	for {
		log.Println("Reading...")

		_, rawMessage, err := conn.ReadMessage()

		if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			log.Println("Error on event read: ", err)
			break
		}

		var message []string
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Printf("Error on unmarshalling message %s", rawMessage)
		}

		if string(message[0]) != messaging.EventMsgType {
			log.Printf("Message of type %v received, whereas only messages of type %v are supported for now", message[0], messaging.EventMsgType)
			continue
		}

		var event messaging.Event
		err = json.Unmarshal([]byte(message[1]), &event)
		if err != nil {
			log.Println("Error on parsing event from message: ", err)
		}

		if !event.Verify() {
			log.Println("ACHTUNG! Event cannot be verified!")
			continue
		}

		log.Println("Received event: ", event) // TODO Debug level to logs!
		ActionsPerEvent[event.Kind](node, &event, r)
	}
}

func (node *Node) Connect(event *messaging.Event, request *http.Request) {
	if node.Config.ForceConnectionRequests {
		newPeer := peer.FromHost(request.Host)
		node.AddPeer(newPeer)
		return
	}

	castPoll("y", node.Peers[:], 0.5)

}

func castPoll(s string, peer []peer.Peer, votesQuorum float64) {
	panic("unimplemented")
}
