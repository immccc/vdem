package node

import (
	"immccc/vdem/vote"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVote(t *testing.T) {

	assertConnections := func(t *testing.T, nodes []*Node) {
		t.Log("Giving time nodes to be connected...")
		time.Sleep(time.Second * 10)

		node0 := nodes[0]
		// node1 := nodes[1]
		// node2 := nodes[2]

		//Node 0 creates a voting
		poll := &vote.Poll{
					Id:           uuid.New(),
					Description:  "Testing votes is correct?",
					Options:      []string{"No", "Yes"},

		}
		node0.OpenPoll(poll)

		t.Log("Giving time nodes to receive voting request...")
		time.Sleep(time.Second * 10)

		assertPollIsOpenInNodes(t, nodes)

		//Nodes emits an update of it
		// voteId := "to be replaced by real vote id"
		// updateEvent := messaging.UpdatePollEvent(voteId, 0)
		// updateEvent.Sign(node0.Config.PrivateKey)
		// node0.Send(&updateEvent)

		// updateEvent = messaging.UpdatePollEvent(voteId, 1)
		// updateEvent.Sign(node1.Config.PrivateKey)
		// node1.Send(&updateEvent)

		// updateEvent = messaging.UpdatePollEvent(voteId, 0)
		// updateEvent.Sign(node2.Config.PrivateKey)
		// node2.Send(&updateEvent)

		// //Assert all nodes have latest update
		// for _, node := range nodes {

		// 	votesPerOption := node.GetVotes(voteId)
		// 	if votesPerOption["Yes"] != 2 {
		// 		t.Error("Option Yes must be 2, but got", votesPerOption["Yes"])
		// 	}
		// 	if votesPerOption["No"] != 1 {
		// 		t.Error("Option Yes must be 1, but got", votesPerOption["No"])
		// 	}

		// }

	}

	RunTestOnMultipleNodesSetup(t, assertConnections, &Node1Config, &Node2Config, &Node3Config)

}

func TestVoteFailsFromNewPeer(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestVoteWhenParticipantLeavesAndRejoins(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestPollIsRejectedWhenPeersHaveIntruders(t * testing.T) {
	t.Error("Not yet implemented!")
}

func TestPollIsRejectedWhenCurrentNodeNotPresentInPeers(t* testing.T) {
	t.Error("Not yet implemented!")
}
func TestVoteFailsWhenClose(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestVoteWhenPeerIsUnsync(t *testing.T) {

}


// Assertions

func assertPollIsOpenInNodes(t *testing.T, nodes []*Node) {
	t.Log("Asserting all nodes have the voting registered...")
	for _, node := range nodes {
		if len(node.pollsById) != 1 {
			t.Error("Node ", node.Config.PubKey, "has not registered the voting")
		}

		for pollId, poll := range node.pollsById {
			if pollId != poll.Id.String() {
				t.Error("Node ", node.Config.PubKey, " has not registered the voting with id", pollId, "but got", poll.Id)
			}

			if poll.Description != "Testing votes is correct?" {
				t.Error("Node ", node.Config.PubKey, " has not registered the voting with description 'Testing votes is correct?', but got", poll.Description)
			}

			if poll.Options[0] != "No" {
				t.Error("Node ", node.Config.PubKey, " has not registered the voting with option 'No', but got", poll.Options[0])
			}

			if poll.Options[1] != "Yes" {
				t.Error("Node ", node.Config.PubKey, " has not registered the voting with option 'Yes', but got", poll.Options[1])
			}

			if len(poll.Options) != 2 {
				t.Error("Node ", node.Config.PubKey, " has not registered the voting with 2 options, but got", len(poll.Options))
			}

		}
	}
}
