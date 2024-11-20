package node

type NodeConfig struct {
	PrivateKey              string
	PubKey                  string
	ServerPublicHost		string
	ServerPort              int
	ForceConnectionRequests bool
	ForceAcknowledge		bool
}
