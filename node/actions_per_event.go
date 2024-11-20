package node

import (
	"errors"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"maps"
)

type actionFunction func(*Node, *messaging.Event) error

func (node *Node) getActionsPerEvent(event *messaging.Event) (actionFunction, error) {

	switch event.Kind {
	case messaging.ConnectionAttemptKind:
		return func(node *Node, event *messaging.Event) error {
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
		}, nil
	case messaging.OtherPeersOnNetworkNotificationKind:
		return func(node *Node, event *messaging.Event) error {
			for _, tag := range event.Tags {
				tagType, pubKey, url := tag[0], tag[1], tag[2]

				if tagType != messaging.TagTypeUser {
					continue
				}

				if node.Config.PubKey == pubKey {
					continue
				}

				pr := peer.New(pubKey, url, nil)
				node.AddPeer(&pr, false)
			}

			if node.eventsWaitingForRegisterSenderById[event.PubKey] == nil {
				return nil
			}

			for _, pendingEvent := range node.eventsWaitingForRegisterSenderById[event.PubKey] {
				node.ConfirmEvent(&pendingEvent)
			}

			clear(node.eventsWaitingForRegisterSenderById[event.PubKey])
			return nil
		}, nil

	}

	return nil, errors.New("unregistered event kind")
}
