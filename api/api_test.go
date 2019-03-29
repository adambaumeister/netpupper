package api

import (
	"testing"
)

func TestStartServerApi(t *testing.T) {
	go StartServerApi("8999")
	a := API{}
	a.SendRegister()
	a.StartbwTest(
		"127.0.0.1:8999",
		"127.0.0.1:8080",
		"500m",
	)
}
