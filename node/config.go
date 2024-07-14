package node

type NodeConfig struct {
	PrivateKey              string
	PubKey                  string
	ServerPort              int
	ForceConnectionRequests bool
}
