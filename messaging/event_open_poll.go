package messaging


const (
	OpenPollKind uint16 = 1003
)

func OpenPollEvent(description string, options []string) Event {

	event := Event{
		Kind:    OpenPollKind,
		Tags:    [][]string{{description}},
		Content: "Vote request! THIS CONTENT SHOULD BE REPLACED WITH SOME REAL INFORMATION!",
	}

	return event
}
