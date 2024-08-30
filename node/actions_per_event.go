package node

import (
	"immccc/vdem/messaging"
	"net/http"
)

type actionFunction func(*Node, *messaging.Event, *http.Request)

var ActionsPerEvent = map[uint16]actionFunction{
	messaging.ConnectionAttemptKind: func(node *Node, event *messaging.Event, request *http.Request) {
		node.Connect(event, request)
	},
}
