package messaging

const (
	UpdatePollKind uint16 = 30001
)

func UpdatePollEvent(votedId string, selection uint8) Event {

	event := Event{
		Kind:    OpenPollKind,
		Tags:    [][]string{},
		Content: "Vote update! THIS CONTENT SHOULD BE REPLACED WITH SOME REAL INFORMATION!",
	}

	return event
}
