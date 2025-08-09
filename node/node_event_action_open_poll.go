package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/vote"
	"strings"
)

func ProcessEventOpenPoll(node *Node, event *messaging.Event) error {
	decoder := json.NewDecoder(strings.NewReader(event.Content))

	poll := vote.Poll{}
	err := decoder.Decode(&poll)
	if err != nil {
		return fmt.Errorf("failed to decode poll creation event: %v", err)
	}

	if isNodeNotAllowedInPoll(node, &poll) {
		return errors.New("node is not allowed in the poll")
	}

	if areIntruderPeersInPoll(node, &poll) {
		return fmt.Errorf("intruder peers in poll: %v", poll.AllowedPeers)
	}

	node.AddActivePoll(&poll)
	return nil
}

func isNodeNotAllowedInPoll(node *Node, poll *vote.Poll) bool {
	for peerPubKeyInPoll := range poll.AllowedPeers {
		if peerPubKeyInPoll == node.Config.PubKey {
			return false
		}
	}
	return true
}

func areIntruderPeersInPoll(node *Node, poll *vote.Poll) bool {
	for peerPubKeyInPoll := range poll.AllowedPeers {
		if _, exists := node.PeersByPubKey[peerPubKeyInPoll]; !exists {
			return true
		}
	}
	return false
}