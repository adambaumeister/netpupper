package tcpbw

import (
	"fmt"
	"testing"
)

func _TestTcpbw(t *testing.T) {
	// Start the server
	s := Server{}
	s.Configure("..\\daemon_server.yml")
	s.Config.Address = "127.0.0.1:8080"
	go s.Run()

	// Start the client
	c := Client{}
	c.Configure("..\\daemon.yml")
	c.Config.Server = "127.0.0.1:8080"

	c.Run()

}

func TestClient_Configure(t *testing.T) {
	s := Server{}
	s.Configure("C:/users/adam/go/src/github.com/adamb/netpupper/daemon_server.yml")
	fmt.Printf("%v\n", s.Config.Address)
}
