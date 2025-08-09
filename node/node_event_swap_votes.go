package node

import (
	"fmt"
	"immccc/vdem/messaging"
	"log"
	"maps"
)

func ProcessSwapVotesEvent(node *Node, event *messaging.Event) error {

	pollId := event.Tags[1][1]

	poll, found := node.pollsById[pollId]
	if !found {
		return fmt.Errorf("received vote swap for unknown poll %s", pollId)
	}

	ongoingShufflingPeer := poll.GetOnGoingShufflingPeer()
	if ongoingShufflingPeer != nil && *ongoingShufflingPeer != event.PubKey {
		return fmt.Errorf("awaiting shuffle response from %s, but got response from %s", *ongoingShufflingPeer, event.PubKey)
	}

	if len(poll.GetPeersWithUnshuffledVote()) == 0 {
		return fmt.Errorf("attempting to swap votes after all votes have been shuffled")
	}

	fromPeerKey := event.Tags[0][1]

	vote := poll.GetLocalVote()
	if vote == nil {
		return fmt.Errorf("attempting to swap before a vote has been created")
	}

	confirmSwapEvent := messaging.SwapVoteKeysEvent(pollId, node.Config.PubKey, poll.HoldKey)
	peerDest, found := node.PeersByPubKey[fromPeerKey]

	if !found {
		return fmt.Errorf("attempting to swap votes with unknown peer %s", fromPeerKey)
	}
	node.Send(&confirmSwapEvent, &peerDest)
	

	poll.HoldKey = event.Content

	maps.DeleteFunc(
		poll.GetPeersWithUnshuffledVote(),
		func(peerId string, _ int) bool {
			return peerId == fromPeerKey
		},
	)

	if ongoingShufflingPeer != nil {
		poll.ResetOngoingshufflingPeer()
	} else {
		poll.SetOngoingShufflingPeer(fromPeerKey)
	}

	log.Printf("node: %s, peer: %s, poll: %s, status: VOTES_SHUFFLED", node.Config.PubKey, node.Config.PubKey, pollId)

	return nil
}
