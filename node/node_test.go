package node

import (
	"fmt"
	http "net/http"
	"testing"
)

func TestNodeListens(t *testing.T) {
	config1 := NodeConfig{
		PrivateKey: "abcd",
		PubKey:     "dbca",
		ServerPort: 3394,
	}

	config2 := NodeConfig{
		PrivateKey: "abcd",
		PubKey:     "dbca",
		ServerPort: 3395,
	}

	quit_node_1 := make(chan bool)
	quit_node_2 := make(chan bool)

	go func() {
		go Start(&config1)

		for {
			select {
			case <-quit_node_1:
				return
			}
		}
	}()

	_, err := http.Get(fmt.Sprintf("http://localhost:%v", config.ServerPort))
	quit_node_1 <- true

	if err != nil {
		t.Error("Request to server could not be performed", err)
	}

	if NodeConfiguration != config {
		t.Error("Server configuration has not been registered")
	}

}
