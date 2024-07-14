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
)

func TestNodeConnects(t *testing.T) {

	var wg sync.WaitGroup

	node1 := createAndStartNode(&wg, "Kx67AX7YZ6VCvBR7qGz35wxVaRku4Gvg5Pa445TEGonWYCG8AZmL", "1BE5XHY3AAeZ72pBbKggjCeUxRAQk8XY6x", portNode1, nil)
	node2 := createAndStartNode(
		&wg, "KxLaBhSXFaosxuxXzhmTsGoLd6FEA9g3J9coZaY87smykZ6JC9je", "1FLYde1vFvtNDiasdxkC9jEBH7v69nKa1x", portNode2,
		[]peer.Peer{
			{
				Port: portNode1,
				Host: "localhost",
			},
		},
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

func createAndStartNode(wg *sync.WaitGroup, privKey string, pubKey string, port int, peers []peer.Peer) *Node {
	node := Node{Config: NodeConfig{
		PrivateKey: privKey,
		PubKey:     pubKey,
		ServerPort: port,
		ForceConnectionRequests: true,
	}}

	for _, pr := range peers {
		node.AddPeer(pr)
	}

	wg.Add(1)
	go node.Start(wg)
	wg.Wait()
	return &node
}
