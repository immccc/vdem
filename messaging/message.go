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
)

type MessageType string

func BuildMessageWithEvent(event *Event) string {
	eventAsJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("Unable to build Message, can't convert to JSON: %v !", err)
	}

	return fmt.Sprintf("[%v, %v]", EventMsgType, eventAsJson)
}
