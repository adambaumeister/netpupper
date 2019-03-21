package tcpbw

import (
	"testing"
)

func TestTcpbw(t *testing.T) {
	// Start the server
	s := Server{}
	s.Configure()
	s.Config.Address = "127.0.0.1:8080"
	go s.Run()

	// Start the client
	c := Client{}
	c.Configure()
	c.Config.Server = "127.0.0.1:8080"

	c.Run()

}
