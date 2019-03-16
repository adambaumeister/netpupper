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

func TestStringToByte(t *testing.T) {
	s := "20G"
	r := StringToByte(s)
	fmt.Printf("StringToByte: %v:%v\n\n", s, r)
	back := ByteToString(r)
	fmt.Printf("ByteToString: %v:%v\n\n", r, back)
}
