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

func BuildEventMessage(event *Event) []byte {
	eventAsJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("Unable to build Message, can't convert to JSON: %v !", err)
	}

	return []byte(fmt.Sprintf(`["%s", %s]`, EventMsgType, eventAsJson))
}
