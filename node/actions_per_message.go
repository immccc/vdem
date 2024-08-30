package node

import (
	"immccc/vdem/messaging"
	"log"
	"net/http"
)

type actionFunctionForMessage func(node *Node, message []any, params ...any)

var ActionsPerMessage = map[string]actionFunctionForMessage{
	messaging.EventMsgType: func(node *Node, message []any, params ...any) {
		r := params[0].(*http.Request)
		err := node.ParseEvent(&message[0], r)
		if err != nil {
			log.Println("Unable to parse event: ", err)
		}
	},
	messaging.OkType: func(node *Node, message []any, params ...any) {
		log.Println("OK TYPE!!!!!!!") // TODO
	},
	messaging.ReqMsgType: func(node *Node, message []any, params ...any) {
		log.Println("REQ TYPE!!!!!!!!") // TODO
	},
}
