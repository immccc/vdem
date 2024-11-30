package node

import (
	"log"
)

type VoteResult map[string]uint64

func (node *Node) GetVotes(voteId string) VoteResult {
	log.Println("Not yet implemented")
	empty := make(VoteResult)
	return empty
}
