// Actions on node regarding voting

package node

import (
	"crypto/rand"
	"errors"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/vote"
	"math/big"
)

func (node *Node) AddActivePoll(poll *vote.Poll) {
	if node.pollsById == nil {
		node.pollsById = make(map[string]vote.Poll)
	}

	node.pollsById[poll.Id.String()] = *poll

}

func (node *Node) OpenPoll(poll *vote.Poll) error {

	if poll.AllowedPeers != nil {
		return errors.New("poll must not have allowed peers set before opening")
	}

	allowedPeers := node.getPeersAsArrayOfPointers()
	poll.SetPeers(allowedPeers...)

	openEvent := messaging.OpenPollEvent(poll)
	openEvent.Sign(node.Config.PrivateKey)
	node.Send(&openEvent, allowedPeers...)

	// TODO What if poll is rejected by a peer? It should be removed
	node.AddActivePoll(poll)

	return nil
}

func (node *Node) RequestShuffleVoteWithRandomPeer(pollId string) error {
	// Shuffle the vote with a random peer of the network. This step is necessary to ensure that the vote is not directly linked to the voter.
	poll, ok := node.pollsById[pollId]
	if !ok {
		return fmt.Errorf("poll with id %s not found", pollId)
	}

	if poll.GetOnGoingShufflingPeer() != nil {
		return fmt.Errorf("shuffling ongoing on poll with id %s", pollId)
	}

	peerKeys := make([]string, 0, len(poll.GetPeersWithUnshuffledVote()))
	for peerKey := range poll.GetPeersWithUnshuffledVote() {
		peerKeys = append(peerKeys, peerKey)
	}

	randomKeyIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(peerKeys))))
	if err != nil {
		return fmt.Errorf("failed to generate random index: %w", err)
	}

	peerKey := peerKeys[randomKeyIdx.Int64()]
	peer := node.PeersByPubKey[peerKey]

	vote := poll.GetLocalVote()
	if vote == nil {
		return fmt.Errorf("no local vote found for poll with id %s", pollId)
	}

	event := messaging.SwapVoteKeysEvent(peer.PubKey, peerKey, poll.HoldKey)
	peer.SendMessage(messaging.BuildEventMessage(&event))

	return nil

}
