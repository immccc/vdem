package node

import (
	"errors"
	"immccc/vdem/messaging"
)

type actionFunction func(*Node, *messaging.Event) error

func (node *Node) getActionsPerEvent(event *messaging.Event) (actionFunction, error) {

	switch event.Kind {
	case messaging.ConnectionAttemptKind:
		return ProcessEventConnectionAttempt, nil
	case messaging.OtherPeersOnNetworkNotificationKind:
		return ProcessEventOtherPeersOnNetwork, nil
	case messaging.OpenPollKind:
		return ProcessEventOpenPoll, nil
	case messaging.SwapVoteKeysKind:
		return ProcessSwapVotesEvent, nil
	default:
		return nil, errors.New("unregistered event kind")
	}
}
