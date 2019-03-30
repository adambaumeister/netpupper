package api

import (
	"fmt"
	"testing"
)

func _TestStartServerApi(t *testing.T) {
	go StartServerApi("8999")
	a := API{}
	a.SendRegister("127.0.0.1:8999")
	a.StartbwTest(
		"127.0.0.1:8999",
		"127.0.0.1:8080",
		"500m",
	)
}

func TestConfigure(t *testing.T) {
	a := API{}
	a.Configure("..\\daemon.yml")
	for _, s := range a.Config.Servers {
		fmt.Printf("Server;%v\n", s)
	}
}
