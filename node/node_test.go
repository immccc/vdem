package node

import (
	"testing"
	"time"
)

func TestNodeConnects(t *testing.T) {
	RunTestOnMultipleNodesSetup(t, assertConnections, &Node1Config, &Node2Config)
	// TODO Test what happens when a node connecting to other is rejected
}

func TestNodesConnectsOnNetworkOfMany(t *testing.T) {
	RunTestOnMultipleNodesSetup(t, assertConnections, &Node1Config, &Node2Config, &Node3Config)
}

func assertConnections(t *testing.T, nodes []*Node) {
	start_secs := time.Now().Unix()
	after_secs := time.Now().Unix()

	for after_secs-start_secs < 20 {
		after_secs = time.Now().Unix()
		for _, node := range nodes {
			if len(node.PeersByPubKey) != len(nodes)-1 {
				time.Sleep(time.Second)
				break
			}
		}
	}

	for idx, node := range nodes {
		if len(node.PeersByPubKey) != len(nodes)-1 {
			t.Error("Node", idx, "has not been updated with others peers.")
		}
	}
}
