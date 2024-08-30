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

	event.Sign(node.Config.PrivateKey)
	peer.SendMessage(messaging.BuildEventMessage(event))

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

		var message []any
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Printf("Error on unmarshalling message %s", rawMessage)
			continue
		}

		eventType, eventTypeIsString := message[0].(string)
		if !eventTypeIsString {
			log.Println("Message type is not a string: ", eventType)
			continue
		}

		if eventType != messaging.EventMsgType {
			log.Printf("Message of type %v received, whereas only messages of type %v are supported for now", message[0], messaging.EventMsgType)
			continue
		}

		eventAsMap, eventIsMap := message[1].(map[string]interface{})
		if !eventIsMap {
			log.Println("Event expected in msg is not a a map!")
			continue
		}

		eventAsStr, err := json.Marshal(eventAsMap)
		if err != nil {
			log.Println("Another failure on serialization. This time, map -> stringified json", err)
		}

		var event messaging.Event
		err = json.Unmarshal(eventAsStr, &event)
		if err != nil {
			log.Println("Error on stringified json -> Event ", err)
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

func (node *Node) Close() {
	for _, pr := range node.Peers {
		pr.Close()
	}
}

func castPoll(s string, peer []peer.Peer, votesQuorum float64) {
	panic("unimplemented")
}
