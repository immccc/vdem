package node

import (
	"immccc/vdem/peer"
	"sync"
	"testing"
)

const (
	portNode1   int    = 3394
	portNode2   int    = 3395
	portNode3   int    = 3396
	pubKeyNode1 string = "035653c0e9250de356c13e11dc4887e9d07f571c7fc9d761e8d8fdd6ad231c388c"
	host        string = "localhost"
)

var (
	Node1Config NodeConfig = NodeConfig{
		PrivateKey:              "Kx67AX7YZ6VCvBR7qGz35wxVaRku4Gvg5Pa445TEGonWYCG8AZmL",
		PubKey:                  pubKeyNode1, // TODO PubKey to be generated, and not manually assigned?
		ServerPort:              portNode1,
		ServerPublicHost:        host,
		ForceConnectionRequests: true,
		ForceAcknowledge:        true,
	}
	Node2Config NodeConfig = NodeConfig{
		PrivateKey:              "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je",
		PubKey:                  "0319b82bc08e8dd49f089ff839bada602c2c42e4166f4916924464c63b6d3b4c31",
		ServerPort:              portNode2,
		ServerPublicHost:        host,
		ForceConnectionRequests: true,
		ForceAcknowledge:        true,
	}
	Node3Config NodeConfig = NodeConfig{
		PrivateKey:              "L3FRLDYALav5dKi6MgEKvfRaAP3jgeatRnU44uopzNvFyetW55E4",
		PubKey:                  "029979a9c18820104050fb337e6bb4e229f138f377996c6dbffe1857d6c5d52eb6",
		ServerPort:              portNode3,
		ServerPublicHost:        host,
		ForceConnectionRequests: true,
		ForceAcknowledge:        true,
	}
)

// Performs a test into a node group setup
func RunTestOnMultipleNodesSetup(t *testing.T, testFunc func(_ *testing.T, _ []*Node), nodeConfigs ...*NodeConfig) {
	var wg sync.WaitGroup

	peerNode0 := peer.Peer{
		PubKey: nodeConfigs[0].PubKey,
		Port:   nodeConfigs[0].ServerPort,
		Host:   nodeConfigs[0].ServerPublicHost,
	}

	nodes := make([]*Node, 0)

	for idx, nodeConfig := range nodeConfigs {
		node := CreateAndStartNode(&wg, nodeConfig)
		defer node.Close()

		nodes = append(nodes, node)
		if idx > 0 {
			node.AddPeer(&peerNode0, true)
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunc(t, nodes[:])
	}()

	wg.Wait()

}

func CreateAndStartNode(wg *sync.WaitGroup, config *NodeConfig) *Node {
	node := Node{Config: *config}

	wg.Add(1)
	go node.Start(wg)
	wg.Wait()
	return &node
}
