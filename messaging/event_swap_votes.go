package messaging

const (
	SwapVoteKeysKind uint16 = 20010
)

func SwapVoteKeysEvent(poll string, fromPeerKey string, key string) Event {

	event := Event{
		Kind:    SwapVoteKeysKind,
		Tags:    [][]string{{"p", fromPeerKey}, {"e", poll}},
		Content: key,
	}

	return event
}
