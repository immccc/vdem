package messaging

const (
	connectionAttemptKind uint16 = 20001
)

func BuildConnectionAttemptMsg(pubKey string) Message {

	msg := Message{
		Kind:    connectionAttemptKind,
		Tags:    [][]string{[]string{TagTypeUser, pubKey}},
		Content: "Please let me connect!",
	}

	return msg
}
