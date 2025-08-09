package messaging

import (
	"encoding/json"
	"immccc/vdem/vote"
)

const (
	OpenPollKind uint16 = 10003
)

func OpenPollEvent(poll *vote.Poll) Event {
	
	poll_as_json, _ := json.MarshalIndent(poll, " ", " ")

	event := Event{  
		Kind:    OpenPollKind,
		Tags:    [][]string{{}},
		Content: string(poll_as_json),
	}

	return event
}
