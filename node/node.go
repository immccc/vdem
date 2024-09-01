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

	"github.com/gorilla/websocket"
)

// A Node acts as a Nostr relay. Holds instances of connected clients and it's in charge of handling
// and storing status of the decentralized network.
type Node struct {
	Config        NodeConfig
	PeersByPubKey map[string]peer.Peer

	eventsWaitingForConfirmationById   map[string]messaging.Event
	eventsWaitingForRegisterSenderById map[string][]messaging.Event
}

var upgrader = websocket.Upgrader{}

func (node *Node) Start(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	log.Printf("node: %s, peers: %d", node.Config.PubKey, len(node.PeersByPubKey))

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
			log.Fatalf("node: %s, err: %v", node.Config.PubKey, err)
		}
	}()

	for !serverStarted {
		//TODO https if applicable
		http.Get(fmt.Sprintf("http://localhost:%v/health", node.Config.ServerPort))
		time.Sleep(time.Second)
	}
}

func (node *Node) AddPeer(pr *peer.Peer, connect bool) error {
	if node.PeersByPubKey == nil {
		node.PeersByPubKey = make(map[string]peer.Peer)
	}

	node.PeersByPubKey[pr.PubKey] = *pr
	log.Printf("node: %s, peer: %s, status: ADDED", node.Config.PubKey, pr.PubKey)
	log.Printf("node: %s, peers: %d", node.Config.PubKey, len(node.PeersByPubKey))

	if !connect {
		return nil
	}

	event := messaging.BuildConnectionAttemptEvent(node.Config.PubKey, "", node.Config.ServerPort)
	event.Sign(node.Config.PrivateKey)
	err := node.Send(&event)

	if err != nil {
		return err
	}

	log.Printf("node: %s, peer: %s, status: CONNECTION_PENDING", node.Config.PubKey, pr.PubKey)

	return nil
}

func (node *Node) Send(event *messaging.Event) error {

	event.Sign(node.Config.PrivateKey)
	if len(node.PeersByPubKey) == 0 {
		return errors.New("no peers to send request")
	}

	for _, peer := range node.PeersByPubKey {
		peer.SendMessage(messaging.BuildEventMessage(event))
	}

	node.addEventAwaitingConfirmation(event)

	return nil
}

func (node *Node) addEventAwaitingConfirmation(event *messaging.Event) {
	if node.eventsWaitingForConfirmationById == nil {
		node.eventsWaitingForConfirmationById = make(map[string]messaging.Event)
	}

	node.eventsWaitingForConfirmationById[event.Id] = *event

	log.Printf("event: %s, kind: %d, status: PENDING", event.Id, event.Kind)
}

func (node *Node) handleEventPeerExistence(event *messaging.Event) bool {
	_, existPeer := node.PeersByPubKey[event.PubKey]
	if existPeer {
		return true
	}

	if node.eventsWaitingForRegisterSenderById == nil {
		node.eventsWaitingForRegisterSenderById = make(map[string][]messaging.Event)
	}

	if node.eventsWaitingForRegisterSenderById[event.PubKey] == nil {
		node.eventsWaitingForRegisterSenderById[event.PubKey] = make([]messaging.Event, 1)
	}

	node.eventsWaitingForRegisterSenderById[event.PubKey] = append(node.eventsWaitingForRegisterSenderById[event.PubKey], *event)

	return false
}

func (node *Node) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade connection error: ", err)
		return
	}

	defer conn.Close()

	for {
		_, rawMessage, err := conn.ReadMessage()

		if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			log.Println("Error on event read: ", err)
			break
		}

		var message []any
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Printf("node: %s, message: %s, %v", node.Config.PubKey, rawMessage, err)
			continue
		}

		eventType, eventTypeIsString := message[0].(string)
		if !eventTypeIsString {
			log.Printf("node: %s, eventType: %v, error: not a string", node.Config.PubKey, eventType)
			continue
		}

		if _, existsKey := ActionsPerMessage[eventType]; !existsKey {
			log.Printf("node: %s, eventType: %v, error: action not implemented for this kind of message", node.Config.PubKey, eventType)
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

	//Enqueue event if source is not yet registered
	existsPeer := node.handleEventPeerExistence(&event)
	if !existsPeer {
		log.Printf("node: %s, event: %s, peer: %s, status: FROM_UNKNOWN_PEER", node.Config.PubKey, event.Id, event.PubKey)
		return nil
	}

	//Confirm received event
	err = node.ConfirmEvent(&event)
	if err != nil {
		return err
	}

	return nil
}

func (node *Node) ConfirmEvent(event *messaging.Event) error {
	if !node.Config.ForceAcknowledge {
		return errors.New("acknowledge of an event upon criteria is not implemented yet")
	}

	log.Printf("node: %s, event: %s, status: ACCEPTED", node.Config.PubKey, event.Id) // TODO Debug level to logs!
	eventAction, err := node.getActionsPerEvent(event)
	if err != nil {
		return err
	}

	err = eventAction(node, event)
	if err != nil {
		log.Printf("node: %s, event: %s, status: ERROR, err: %v", node.Config.PubKey, event.Id, err)
	} else {
		log.Printf("node: %s, event: %s, status: COMPLETED", node.Config.PubKey, event.Id)
	}

	peerDest := node.PeersByPubKey[event.PubKey]
	peerDest.SendMessage(messaging.BuildOkMessage(event.Id, true))

	// TODO Store non replaceable events as NIP-01 says
	return nil
}

func (node *Node) ChangeEventAcceptance(eventId string, accepted bool) error {
	_, eventFound := node.eventsWaitingForConfirmationById[eventId]
	if !eventFound {
		return errors.New("event is not registered amongst waiting ones")
	}

	if !accepted {
		log.Printf("node: %s, event: %s, status: DENIED", node.Config.PubKey, eventId)
	} else {
		log.Printf("node: %s, event: %s, status: ACCEPTED", node.Config.PubKey, eventId)
	}

	delete(node.eventsWaitingForConfirmationById, eventId)

	return nil
}

func (node *Node) Close() {
	for _, pr := range node.PeersByPubKey {
		pr.Close()
	}
}
