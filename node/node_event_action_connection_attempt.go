package node

import (
	"errors"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"maps"
)

func ProcessEventConnectionAttempt(node *Node, event *messaging.Event) error {
	if !node.Config.ForceConnectionRequests {
		// TODO Implement connection to networks by trusted certificates
		return errors.New("trusted connections not implemented yet")
	}

	pr := peer.New(event.Tags[0][1], event.Tags[0][2], nil)
	node.AddPeer(&pr, false)

	//TODO Dedicated private function to map for the build other peers event
	peers := make([]peer.Peer, 0)
	nodeAsPeer := peer.New(node.Config.PubKey, node.Config.ServerPublicHost, &node.Config.ServerPort)
	peers = append(peers, nodeAsPeer)
	for peer := range maps.Values(node.PeersByPubKey) {
		peers = append(peers, peer)
	}

	networkEvent := messaging.BuildOtherPeersOnNetworkNotificationEvent(peers[:])
	err := node.Send(&networkEvent)

	if err != nil {
		return err
	}

	return nil
}
