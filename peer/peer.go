package peer

import (
	"fmt"
	"strconv"
	"strings"
)

type Peer struct {
	Port int
	Host string
}

func (peer *Peer) ToURL() string {
	return fmt.Sprintf("ws://%v:%v", peer.Host, peer.Port)
}

func FromHost(url string) Peer {
	hostAndPort := strings.Split(url, ":")

	host := hostAndPort[0]
	portStr := hostAndPort[1]
	port, _ := strconv.Atoi(portStr)

	return Peer{
		Host: host,
		Port: port,
	}
}
