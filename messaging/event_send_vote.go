package messaging

import "fmt"

const (
	SendVoteKind uint16 = 1001
)

func SendConnectionAttemptEvent(pubKey string, host string, port int) Event {

	event := Event{
		Kind:    ConnectionAttemptKind,
		Tags:    [][]string{{TagTypeUser, pubKey}},
		Content: fmt.Sprintf("ws://%v:%v", host, port),
	}

	return event
}
