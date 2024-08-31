package node

import (
	"immccc/vdem/peer"
	"sync"
	"testing"
	"time"
)

const (
	portNode1 int = 3394
	portNode2 int = 3395
	portNode3 int = 3396
)

func TestNodeConnects(t *testing.T) {

	var wg sync.WaitGroup

	node1 := createAndStartNode(
		&wg, &NodeConfig{
			PrivateKey:              "Kx67AX7YZ6VCvBR7qGz35wxVaRku4Gvg5Pa445TEGonWYCG8AZmL",
			PubKey:                  "1BE5XHY3AAeZ72pBbKggjCeUxRAQk8XY6x",
			ServerPort:              portNode1,
			ForceConnectionRequests: true,
		},
	)

	node2 := createAndStartNode(
		&wg, &NodeConfig{
			PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:                  "1FLYde1vFvtNDiasdxkC9jEBH7v69nKa1x",
			ServerPort:              portNode2,
			ForceConnectionRequests: true,
		},
	)

	defer node1.Close()
	defer node2.Close()

	node2.AddPeer(
		&peer.Peer{
			Port: portNode1,
			Host: "localhost",
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

			if len(node1.Peers) != 1 {
				continue
			}

			if len(node2.Peers) != 1 {
				continue
			}
		}

		if len(node1.Peers) != 1 {
			t.Error("Node1 is not updated with peer from Node2.")
		}

		if len(node2.Peers) != 1 {
			t.Error("Node2 is not updated with peer from Node1.")
		}
	}()

	wg.Wait()

}

func TestNodesPerformVoting(t *testing.T) {
	peerNode1 := peer.Peer{
		Port: portNode1,
		Host: "localhost",
	}

	var wg sync.WaitGroup
	node1 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:                  "1BE5XHY3AAeZ72pBbKggjCeUxRAQk8XY6x",
			ServerPort:              portNode1,
			ForceConnectionRequests: true,
		},
	)
	node2 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey: "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
			PubKey:     "1FLYde1vFvtNDiasdxkC9jEBH7v69nKa1x",
			ServerPort: portNode2,
		},
	)

	node3 := createAndStartNode(
		&wg,
		&NodeConfig{
			PrivateKey: "L3FRLDYALav5dKi6MgEKvfRaAP3jgeatRnU44uopzNvFyetW55E4",
			PubKey:     "123CjemJyn9ZPdQXRGDbF3kSdhQKCELSFT",
			ServerPort: portNode3,
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

			if len(node1.Peers) != expected_peers {
				continue
			}

			if len(node2.Peers) != expected_peers {
				continue
			}

			if len(node3.Peers) != expected_peers {
				continue
			}
		}

		if len(node1.Peers) != expected_peers {
			t.Errorf("Node1 has only %d peers, expected %d.", len(node1.Peers), expected_peers)
		}

		if len(node2.Peers) != expected_peers {
			t.Errorf("Node2 has only %d peers, expected %d.", len(node2.Peers), expected_peers)
		}

		if len(node3.Peers) != expected_peers {
			t.Errorf("Node3 has only %d peers, expected %d.", len(node3.Peers), expected_peers)
		}

	}()

	wg.Wait()

	// TODO Remove!
	t.Error("Test to be written yet!")

	// TODO what to test
	// node 2 requests voting
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
