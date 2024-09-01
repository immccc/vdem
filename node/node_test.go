package node

import (
	"immccc/vdem/peer"
	"sync"
	"testing"
	"time"
)

const (
	portNode1   int    = 3394
	portNode2   int    = 3395
	portNode3   int    = 3396
	pubKeyNode1 string = "035653c0e9250de356c13e11dc4887e9d07f571c7fc9d761e8d8fdd6ad231c388c"
)

func TestNodeConnects(t *testing.T) {

	var wg sync.WaitGroup

	node1 := createAndStartNode(
		&wg, &NodeConfig{
			PrivateKey:              "Kx67AX7YZ6VCvBR7qGz35wxVaRku4Gvg5Pa445TEGonWYCG8AZmL",
			PubKey:                  pubKeyNode1, // TODO PubKey to be generated, and not manually assigned?
			ServerPort:              portNode1,
			ForceConnectionRequests: true,
			ForceAcknowledge:        true,
		},
	)

	node2 := createAndStartNode(
		&wg, &NodeConfig{
			PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:                  "0319b82bc08e8dd49f089ff839bada602c2c42e4166f4916924464c63b6d3b4c31",
			ServerPort:              portNode2,
			ForceConnectionRequests: true,
			ForceAcknowledge:        true,
		},
	)

	defer node1.Close()
	defer node2.Close()

	node2.AddPeer(
		&peer.Peer{
			PubKey: pubKeyNode1,
			Port:   portNode1,
			Host:   "localhost",
		},
		true,
	)

	wg.Add(1)
	go func() {

		defer wg.Done()

		start_secs := time.Now().Unix()
		after_secs := time.Now().Unix()

		for after_secs-start_secs < 20 {
			after_secs = time.Now().Unix()
			time.Sleep(time.Second)

			if len(node1.PeersByPubKey) != 1 {
				continue
			}

			if len(node2.PeersByPubKey) != 1 {
				continue
			}
		}

		if len(node1.PeersByPubKey) != 1 {
			t.Error("Node1 is not updated with peer from Node2.")
		}

		if len(node2.PeersByPubKey) != 1 {
			t.Error("Node2 is not updated with peer from Node1.")
		}
	}()

	wg.Wait()

	// Test what happens when a node connecting to other is rejected

}

func TestNodesConnectsOnNetworkOfMany(t *testing.T) {
	peerNode1 := peer.Peer{
		PubKey: pubKeyNode1,
		Port:   portNode1,
		Host:   "localhost",
	}

	var wg sync.WaitGroup
	node1 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:                  pubKeyNode1,
			ServerPort:              portNode1,
			ForceConnectionRequests: true,
			ForceAcknowledge:        true,
		},
	)
	node2 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:                  "0319b82bc08e8dd49f089ff839bada602c2c42e4166f4916924464c63b6d3b4c31",
			ServerPort:              portNode2,
			ForceConnectionRequests: true,
			ForceAcknowledge:        true,
		},
	)

	node3 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey:              "L3FRLDYALav5dKi6MgEKvfRaAP3jgeatRnU44uopzNvFyetW55E4",
			PubKey:                  "29979a9c18820104050fb337e6bb4e229f138f377996c6dbffe1857d6c5d52eb6",
			ServerPort:              portNode3,
			ForceConnectionRequests: true,
			ForceAcknowledge:        true,
		},
	)

	defer node1.Close()
	defer node2.Close()
	defer node3.Close()

	node2.AddPeer(&peerNode1, true)
	node3.AddPeer(&peerNode1, true)

	wg.Add(1)
	go func() {

		defer wg.Done()

		const expected_peers = 2

		start_secs := time.Now().Unix()
		after_secs := time.Now().Unix()

		for after_secs-start_secs < 20 {
			after_secs = time.Now().Unix()
			time.Sleep(time.Second)

			if len(node1.PeersByPubKey) != expected_peers {
				continue
			}

			if len(node2.PeersByPubKey) != expected_peers {
				continue
			}

			if len(node3.PeersByPubKey) != expected_peers {
				continue
			}
		}

		if len(node1.PeersByPubKey) != expected_peers {
			t.Errorf("Node1 has only %d peers, expected %d.", len(node1.PeersByPubKey), expected_peers)
		}

		if len(node2.PeersByPubKey) != expected_peers {
			t.Errorf("Node2 has only %d peers, expected %d.", len(node2.PeersByPubKey), expected_peers)
		}

		if len(node3.PeersByPubKey) != expected_peers {
			t.Errorf("Node3 has only %d peers, expected %d.", len(node3.PeersByPubKey), expected_peers)
		}

	}()

	wg.Wait()
}

func testPerformVoting(t *testing.T) {
	// TODO Remove!
	t.Error("Test to be written yet!")

	// TODO what to test
	// node 2 requests voting
	//event := messaging.BuildRequestVoteEvent()
	//event.Sign(node2.Config.PrivateKey)
	//node2.Send(&event)

	// nodes 1 and 2 accepts voting
	// node 2 emits event with voting options to 1 and 2
	// all nodes casts their vote
	// One node does not vote has it hits expiration time. It's counted as "-"
	// Rest of nodes emit their options before timeout
	// In voting backlog, results are stored for all nodes accepted in the network,
	// and those that didn't accept to vote are tracked as "NO SHOW".
	// Backlog is sent to the whole network.
}

func createAndStartNode(wg *sync.WaitGroup, config *NodeConfig) *Node {
	node := Node{Config: *config}

	wg.Add(1)
	go node.Start(wg)
	wg.Wait()
	return &node
}
