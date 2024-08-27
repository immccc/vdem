package event

import "fmt"

const (
	ConnectionAttemptKind uint16 = 20001
)

func BuildConnectionAttemptEvent(pubKey string, host string, port int) Event {

	event := Event{
		Kind:    ConnectionAttemptKind,
		Tags:    [][]string{{TagTypeUser, pubKey}},
		Content: fmt.Sprintf("ws://%v:%v", host, port),
	}

	return event
}
