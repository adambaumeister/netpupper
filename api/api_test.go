package api

import (
	"testing"
)

func TestStartServerApi(t *testing.T) {
	go StartServerApi()
	a := APIClient{}
	a.SendRegister()
}
