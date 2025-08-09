package messaging

const (
	SwapVotesKind uint16 = 20010
)

func SwapVotesEvent(poll string, fromPeerKey string, encodedSelection string) Event {

	event := Event{
		Kind:    SwapVotesKind,
		Tags:    [][]string{{"p", fromPeerKey}, {"e", poll}},
		Content: encodedSelection,
	}

	return event
}
