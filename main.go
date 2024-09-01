package main

import (
	"flag"
	"immccc/vdem/node"
	"immccc/vdem/peer"
	"log"
)

func main() {
	//TODO WIP. Testing of nodes and peers are performed on tests.
	log.Println("Welcome!!!")

	portPtr := flag.Int("port", 3333, "The port the server will be running on")
	publicKeyPtr := flag.String("publicKey", "", "Public key for signing within the network")
	privateKeyPtr := flag.String("privateKey", "", "Private key for signing within the network")

	peerPubKey := flag.String("pubKey", "", "The pubkey of the peer")
	peerPortPtr := flag.Int("peerPort", 0, "The port to connect to a peer")
	peerHostPtr := flag.String("peerHost", "", "The host the peer will be connected to")

	flag.Parse()

	config := node.NodeConfig{
		ServerPort: *portPtr,
		PubKey:     *publicKeyPtr,
		PrivateKey: *privateKeyPtr,
	}

	peer := peer.Peer{
		PubKey: *peerPubKey,
		Host:   *peerHostPtr,
		Port:   *peerPortPtr,
	}

	node := node.Node{Config: config}
	node.AddPeer(&peer, false)
	node.Start(nil)

}
