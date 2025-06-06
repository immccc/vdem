package messaging

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	EventMsgType = "EVENT"
	ReqMsgType   = "REQ"
	CloseType    = "CLOSE"
	OkType       = "OK"
)

func BuildEventMessage(event *Event) []byte {
	eventAsJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("Unable to build Message, can't convert to JSON: %v !", err)
	}

	return fmt.Appendf(nil, `["%s", %s]`, EventMsgType, eventAsJson)
}

func BuildOkMessage(eventId string, accepted bool) []byte {
	accepted_str := "false"
	if accepted {
		accepted_str = "true"
	}
	return fmt.Appendf(nil, `["%s", "%s", %s]`, OkType, eventId, accepted_str)
}

func BuildReqMessage(subscriptionId string) []byte {
	return fmt.Appendf(nil, `["%s", "%s"]`, ReqMsgType, subscriptionId)
}
