package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"log"
	http "net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// A Node acts as a Nostr relay. Holds instances of connected clients and it's in charge of handling
// and storing status of the decentralized network.
type Node struct {
	Config NodeConfig
	Peers  map[string]peer.Peer
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
}

func (node *Node) AddPeer(pr *peer.Peer, connect bool) {
	if node.Peers == nil {
		node.Peers = make(map[string]peer.Peer)
	}

	node.Peers[pr.ToURL()] = *pr

	if !connect {
		return
	}

	event := messaging.BuildConnectionAttemptEvent(node.Config.PubKey, "", node.Config.ServerPort)
	event.Sign(node.Config.PrivateKey)
	node.Send(&event, pr)
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

		if _, existsKey := ActionsPerMessage[eventType]; !existsKey {
			log.Printf("Received msg of type %s is not implemented and can't be processed!", eventType)
			continue
		}
		ActionsPerMessage[eventType](node, message[1:], r)
	}
}

func (node *Node) ParseEvent(rawEvent *any, r *http.Request) error {
	eventAsMap, eventIsMap := (*rawEvent).(map[string]interface{})
	if !eventIsMap {
		return errors.New("event expected in msg is not a a map")

	}

	eventAsStr, err := json.Marshal(eventAsMap)
	if err != nil {
		return err
	}

	var event messaging.Event
	err = json.Unmarshal(eventAsStr, &event)
	if err != nil {
		return err
	}

	if !event.Verify() {
		return errors.New("cannot verify and trust message from event")

	}

	log.Println("Received event: ", event) // TODO Debug level to logs!
	ActionsPerEvent[event.Kind](node, &event, r)

	return nil
}

func (node *Node) Connect(event *messaging.Event, request *http.Request) {
	pr := peer.FromHost(request.Host)

	if node.Config.ForceConnectionRequests {
		node.AddPeer(&pr, false)
		return
	}

	// TODO Consensus is just one way to accept new connections, and a bit naive.
	// Certificates emitted by networks themselves are a better way to connect to, 
	// but it's a good way to start testing voting protocol.
	// TODO Also refactor to a separated function
	u := uuid.New()
	for _, prDest := range node.Peers {
		prDest.SendMessage(messaging.BuildReqMessage(u.String()))
	}

}

func (node *Node) Close() {
	for _, pr := range node.Peers {
		pr.Close()
	}
}
