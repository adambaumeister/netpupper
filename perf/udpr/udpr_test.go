package udpr

import (
	"testing"
	"time"
)

func TestClient_Run(t *testing.T) {
	s := Server{}
	s.Configure("")
	s.Config.Address = "127.0.0.1:9500"
	go s.Run()

	c := Client{}
	c.Configure("")
	c.Config.Server = "127.0.0.1:9500"
	c.Config.PacketCount = 50000
	c.Run()
	time.Sleep(1 * time.Second)
}

func TestCalcLimiter(t *testing.T) {
	CalcLimiter(512000)
}
