package node

import (
	"immccc/vdem/messaging"
	"log"
	"net/http"
)

type actionFunctionForMessage func(node *Node, message []any, params ...any)

// TODO As node function, like actions_per_event
var ActionsPerMessage = map[string]actionFunctionForMessage{
	messaging.EventMsgType: func(node *Node, message []any, params ...any) {
		r := params[0].(*http.Request)
		err := node.ParseEvent(&message[0], r)
		if err != nil {
			log.Println("Unable to parse event: ", err)
		}
	},
	messaging.OkType: func(node *Node, message []any, params ...any) {
		eventId := message[0].(string)
		accepted := message[1].(bool)
		err := node.ChangeEventAcceptance(eventId, accepted)
		if err != nil {
			log.Printf("Unable to mark event %s as %v, %v", eventId, accepted, err)
		}
	},
	messaging.ReqMsgType: func(node *Node, message []any, params ...any) {
		log.Println("REQ TYPE!!!!!!!!") // TODO
	},
}
