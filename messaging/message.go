package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

	eventAsJsonParsed := []byte(strings.ReplaceAll(string(eventAsJson), `"`, `'`))


	return []byte(fmt.Sprintf(`["%s", "%s"]`, EventMsgType, eventAsJsonParsed))
}
