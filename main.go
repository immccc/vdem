package main

import (
	"flag"
	"immccc/vdem/node"
	"immccc/vdem/peer"
	"log"
)

func main() {
	log.Println("Welcome!!!")

	portPtr := flag.Int("port", 3333, "The port the server will be running on")
	publicKeyPtr := flag.String("publicKey", "", "Public key for signing within the network")
	privateKeyPtr := flag.String("privateKey", "", "Private key for signing within the network")

	peerPortPtr := flag.Int("peerPort", 0, "The port to connect to a peer")
	peerHostPtr := flag.String("peerHost", "", "The host the peer will be connected to")

	flag.Parse()

	config := node.NodeConfig{
		ServerPort: *portPtr,
		PubKey:     *publicKeyPtr,
		PrivateKey: *privateKeyPtr,
	}

	peer := peer.Peer{
		Host: *peerHostPtr,
		Port: *peerPortPtr,
	}

	node := node.Node{Config: config}
	node.AddPeer(peer)
	node.Start(nil)

}
