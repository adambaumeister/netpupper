package api

import (
	"fmt"
	"testing"
)

// Test the API server
// Uses the example daemon* configs, found in the root of this repo
func TestStartServerApi(t *testing.T) {
	// Server
	a := API{}
	a.Configure("../daemon_server.yml")
	go a.Run()
	// client
	ca := API{}
	ca.Configure("../daemon.yml")
}

func TestConfigure(t *testing.T) {
	a := API{}
	a.Configure("..\\daemon.yml")
	for _, s := range a.Config.Servers {
		fmt.Printf("Server;%v\n", s)
	}
}
