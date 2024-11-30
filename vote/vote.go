package vote

import uuid "github.com/google/uuid"

type PollStatus int

const (
	Open PollStatus = 0
	Closed
)

type Poll struct {
	Id           uuid.UUID
	Description  string
	Options      []string
	Status       PollStatus
	AllowedPeers []string
}

type Vote struct {
	Height         uint64
	Comment        string
	Selection      uint
	UpdatedOptions []string
}
