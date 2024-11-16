package messaging

import (
	"immccc/vdem/peer"
	"iter"
)

const (
	OtherPeersOnNetworkNotificationKind uint16 = 1001
)

func BuildOtherPeersOnNetworkNotificationEvent(prs *iter.Seq[peer.Peer]) Event {

	tags := make([]([]string), 0)

	for pr := range *prs {
		tag := []string{TagTypeUser, pr.PubKey, pr.ToURL()}
		tags = append(tags, tag)
	}

	event := Event{
		Kind:    OtherPeersOnNetworkNotificationKind,
		Tags:    tags,
		Content: "Network notification with the rest of the crew!",
	}

	return event
}
