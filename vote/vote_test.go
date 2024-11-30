package vote

import (
	"immccc/vdem/messaging"
	"immccc/vdem/node"
	"log"
	"testing"
	"time"
)

func TestVote(t *testing.T) {

	assertConnections := func(t *testing.T, nodes []*node.Node) {
		log.Println("Giving time nodes to be connected...")
		time.Sleep(time.Second * 10)

		node0 := nodes[0]
		node1 := nodes[1]
		node2 := nodes[2]

		//Node 0 creates a voting
		options := []string{"No", "Yes"}
		openEvent := messaging.OpenPollEvent("Testing votes is correct?", options)
		openEvent.Sign(node0.Config.PrivateKey)
		node0.Send(&openEvent)

		// Wait until voting has been registered on 3 nodes
		log.Println("Giving time nodes to receive voting request...")
		time.Sleep(time.Second * 5)

		//Nodes emits an update of it
		voteId := "to be replaced by real vote id"
		updateEvent := messaging.UpdatePollEvent(voteId, 0)
		updateEvent.Sign(node0.Config.PrivateKey)
		node0.Send(&updateEvent)

		updateEvent = messaging.UpdatePollEvent(voteId, 1)
		updateEvent.Sign(node1.Config.PrivateKey)
		node1.Send(&updateEvent)

		updateEvent = messaging.UpdatePollEvent(voteId, 0)
		updateEvent.Sign(node2.Config.PrivateKey)
		node2.Send(&updateEvent)

		//Assert all nodes have latest update
		for _, node := range nodes {

			votesPerOption := node.GetVotes(voteId)
			if votesPerOption["Yes"] != 2 {
				t.Error("Option Yes must be 2, but got", votesPerOption["Yes"])
			}
			if votesPerOption["No"] != 1 {
				t.Error("Option Yes must be 1, but got", votesPerOption["No"])
			}

		}

	}

	node.RunTestOnMultipleNodesSetup(t, assertConnections, &node.Node1Config, &node.Node2Config, &node.Node3Config)

}

func TestVoteFailsFromNewPeer(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestVoteWhenParticipantLeavesAndRejoins(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestVoteFailsWhenClose(t *testing.T) {
	t.Error("Not yet implemented!")
}

func TestVoteWhenPeerIsUnsync(t *testing.T) {

}
