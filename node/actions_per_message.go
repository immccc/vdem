package node

import (
	"immccc/vdem/messaging"
	"net/http"
)

type actionFunction func(*Node, *messaging.Message, *http.Request)

var ActionsPerMessage = map[uint16]actionFunction{
	messaging.ConnectionAttemptKind: func(node *Node, msg *messaging.Message, request *http.Request) {
		node.Connect(msg, request)
	},
}
