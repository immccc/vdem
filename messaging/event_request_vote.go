package messaging

const (
	RequestVoteKind uint16 = 1002
)

func BuildRequestVoteEvent() Event {

	event := Event{
		Kind:    RequestVoteKind,
		Tags:    [][]string{},
		Content: "Vote request! THIS CONTENT SHOULD BE REPLACED WITH SOME REAL INFORMATION!",
	}

	return event
}
