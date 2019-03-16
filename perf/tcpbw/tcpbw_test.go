package tcpbw

import (
	"fmt"
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

func TestConvertByteDec(t *testing.T) {
	s := "20M"
	r := ConvertByteDec(s)
	fmt.Printf("ConvertByteDec: %v:%v\n\n", s, r)
}
