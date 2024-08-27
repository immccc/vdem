package node

import (
	"immccc/vdem/event"
	"net/http"
)

type actionFunction func(*Node, *event.Event, *http.Request)

var ActionsPerEvent = map[uint16]actionFunction{
	event.ConnectionAttemptKind: func(node *Node, event *event.Event, request *http.Request) {
		node.Connect(event, request)
	},
}
