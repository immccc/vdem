// Actions on node regarding voting

package node

import (
	"crypto/rand"
	"errors"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
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

func (node *Node) RequestShuffleVoteWithRandomPeer(vote *vote.Vote, peers[] *peer.Peer) error {
	// Shuffle the vote with a random peer of the network. This step is necessary to ensure that the vote is not directly linked to the voter.

	idxPeer, err := rand.Int(rand.Reader, big.NewInt(int64(len(peers))))
    if err != nil {
		return fmt.Errorf("failed to get arandom peer: %v", err)
    }

	peer := peers[idxPeer.Int64()]

	event := messaging.SwapVotesEvent(peer.PubKey, vote.EncodedSelection, vote.EncodedSelection)
	peer.SendMessage(messaging.BuildEventMessage(&event))

	return nil
	
}