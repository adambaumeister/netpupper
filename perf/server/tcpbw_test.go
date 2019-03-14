package server

import (
	"testing"
)

func TestTcpbw(t *testing.T) {
	go Server()
	Client()

}
