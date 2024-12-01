package node

import (
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
)

func ProcessEventOtherPeersOnNetwork(node *Node, event *messaging.Event) error {
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
}
