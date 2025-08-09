package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"immccc/vdem/vote"
	"log"
	"maps"
	http "net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Node acts as a Nostr relay. Holds instances of connected clients and it's in charge of handling
// and storing status of the decentralized network.
type Node struct {
	Config                             NodeConfig
	PeersByPubKey                      map[string]peer.Peer

	eventsWaitingForConfirmationById   map[string]messaging.Event
	eventsWaitingForRegisterSenderById map[string][]messaging.Event

	pollsById 						   map[string]vote.Poll
}

var whitelisted_events_for_unregistered_clients = map[uint16]bool{
	messaging.ConnectionAttemptKind:               true,
	messaging.OtherPeersOnNetworkNotificationKind: true, // TODO OUT!! Add a mechanism to validate a message comes from  same network
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

	if _, alreadyExistsPeer := node.PeersByPubKey[pr.PubKey]; alreadyExistsPeer {
		return nil
	}

	node.PeersByPubKey[pr.PubKey] = *pr
	log.Printf("node: %s, peer: %s, status: ADDED", node.Config.PubKey, pr.PubKey)
	log.Printf("node: %s, peers: %d", node.Config.PubKey, len(node.PeersByPubKey))

	// TODO rename "connect" param
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

func (node *Node) Send(event *messaging.Event, peers... *peer.Peer) error {

	event.Sign(node.Config.PrivateKey)
	if len(node.PeersByPubKey) == 0 {
		return errors.New("no peers to send request")
	}

	if len(peers) == 0 {
		peers = node.getPeersAsArrayOfPointers()
	}

	for _, peer := range peers {
		peer.SendMessage(messaging.BuildEventMessage(event))
		log.Printf("node: %v, event: %v, kind: %d, status: SENT, peer: %v", node.Config.PubKey, event.Id, event.Kind, peer.PubKey)
	}

	peerPubKeys := slices.Collect(maps.Keys(node.PeersByPubKey))
	log.Printf("node: %v, event: %v, kind: %d, status: SENT, peers: %v", node.Config.PubKey, event.Id, event.Kind, peerPubKeys)

	node.addEventAwaitingConfirmation(event)
	return nil
}

func (node *Node) addEventAwaitingConfirmation(event *messaging.Event) {
	if node.eventsWaitingForConfirmationById == nil {
		node.eventsWaitingForConfirmationById = make(map[string]messaging.Event)
	}

	node.eventsWaitingForConfirmationById[event.Id] = *event

	log.Printf("node: %s, event: %s, kind: %d, status: PENDING", node.Config.PubKey, event.Id, event.Kind)
}

func (node *Node) handleEventPeerExistence(event *messaging.Event) bool {
	_, existPeer := node.PeersByPubKey[event.PubKey]
	if existPeer {
		return true
	}

	if whitelisted_events_for_unregistered_clients[event.Kind] {
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

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected error on event read: ", err)
				break
			} else if websocket.IsUnexpectedCloseError(err) {
				log.Printf("node: %s, message: CLOSED_CONNECTION, error: %s", node.Config.PubKey, err)
				break
			}
			// TODO Disconnect peer in case it's disconnected
			continue
		}

		var message []any
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Printf("node: %s, message: %s, error: %v", node.Config.PubKey, rawMessage, err)
			continue
		}

		messageType, messageTypeIsString := message[0].(string)
		if !messageTypeIsString {
			log.Printf("node: %s, messageType: %v, error: not a string", node.Config.PubKey, messageType)
			continue
		}

		if _, existsKey := ActionsPerMessage[messageType]; !existsKey {
			log.Printf("node: %s, messageType: %v, error: action not implemented for this kind of message", node.Config.PubKey, messageType)
			continue
		}
		ActionsPerMessage[messageType](node, message[1:], r)
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
		return err
	} else {
		log.Printf("node: %s, event: %s, status: COMPLETED", node.Config.PubKey, event.Id)
	}

	peerDest := node.PeersByPubKey[event.PubKey]
	peerDest.SendMessage(messaging.BuildOkMessage(event.Id, true))

	log.Printf("node: %s, event: %s, peer: %s, status: CONFIRMED", node.Config.PubKey, event.Id, event.PubKey)

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

func (node *Node) getPeersAsArrayOfPointers() []*peer.Peer {
	peers := make([]*peer.Peer, 0, len(node.PeersByPubKey))
	for pubKey := range node.PeersByPubKey {
		peer := node.PeersByPubKey[pubKey]
		peers = append(peers, &peer)
	}
	return peers
}