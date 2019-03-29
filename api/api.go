package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamb/netpupper/controller"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/tcpbw"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type API struct {
	Controller *controller.Controller
}

func StartServerApi(p string) API {
	con := controller.Controller{}
	a := API{
		Controller: &con,
	}

	http.HandleFunc("/register", a.register)
	http.HandleFunc("/tcpbw", a.bwtest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), nil))

	return a
}

/*
Receive a register from a client
*/
func (a *API) register(w http.ResponseWriter, r *http.Request) {
	reg := controller.Client{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &reg)
	a.Controller.AddClient(reg)

	reg.Addr = strings.Split(r.RemoteAddr, ":")[0]

	m := Message{
		Value: fmt.Sprintf("%v registered to controller.", reg.Addr),
	}
	b, err := json.Marshal(m)
	errors.CheckError(err)
	w.Write(b)

}

func (a *API) bwtest(w http.ResponseWriter, r *http.Request) {
	cc := tcpbw.ClientConfig{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &cc)
	fmt.Printf("DEBUG: Got a test request destination: %v\n", cc.Server)
}

/*
Start a bandwidth test from a client
*/
func (a *API) StartbwTest(addr string, server string, byteCount string) {
	cc := tcpbw.ClientConfig{
		Server: server,
		Bytes:  byteCount,
	}
	b, err := json.Marshal(cc)
	errors.CheckError(err)
	resp, err := http.Post(fmt.Sprintf("http://%v/tcpbw", addr), "application/json", bytes.NewBuffer(b))
	errors.CheckError(err)

	msg := Message{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &msg)
	fmt.Printf("DEBUG: %v\n", msg.Value)
}

/*
Send a register
*/
func (a *API) SendRegister() {
	host, _ := os.Hostname()
	gt := controller.Tag{
		Name:  "name",
		Value: host,
	}
	tags := []controller.Tag{gt}
	reg := controller.Client{
		Tags: tags,
	}
	b, err := json.Marshal(reg)
	errors.CheckError(err)
	resp, err := http.Post("http://127.0.0.1:8999/register", "application/json", bytes.NewBuffer(b))
	errors.CheckError(err)

	msg := Message{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &msg)
	fmt.Printf("DEBUG: %v\n", msg.Value)
}
