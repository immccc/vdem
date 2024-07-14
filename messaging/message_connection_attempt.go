package messaging

import "fmt"

const (
	ConnectionAttemptKind uint16 = 20001
)

func BuildConnectionAttemptMsg(pubKey string, host string, port int) Message {

	msg := Message{
		Kind:    ConnectionAttemptKind,
		Tags:    [][]string{{TagTypeUser, pubKey}},
		Content: fmt.Sprintf("ws://%v:%v", host, port),
	}

	return msg
}
