package messaging

import "fmt"

const (
	ConnectionAttemptKind uint16 = 20001
)

func BuildConnectionAttemptEvent(pubKey string, host string, port int) Event {

	event := Event{
		Kind:    ConnectionAttemptKind,
		Tags:    [][]string{{TagTypeUser, pubKey, fmt.Sprintf("%v:%v", host, port)}},
		Content: "Connection attempt",
	}

	return event
}
