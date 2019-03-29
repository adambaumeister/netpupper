package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamb/netpupper/controller"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/tcpbw"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type API struct {
	Controller *controller.Controller

	Config *APIConfig
}
type APIConfig struct {
	Servers []string
}

/*
Configure the API
Returns TRUE if this method matches the requested config
*/
func (a *API) Configure(cf string) {
	fmt.Printf("Got yaml configuration %v\n", cf)
	var serverFile string
	if len(cf) > 0 {
		serverFile = cf
	} else if len(os.Getenv("NETP_CONFIG")) > 0 {
		serverFile = os.Getenv("NETP_CONFIG")
	}
	// First, try bootstrapping from the YAML server file
	a.Config = &APIConfig{}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); os.IsExist(err) {
		data, err := ioutil.ReadFile(serverFile)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, a.Config)
		errors.CheckError(err)
	}
}
func (a *API) Run() {
	// Register to the provided list of servers, if any
	for _, s := range a.Config.Servers {
		a.SendRegister(s)
	}
	StartServerApi("8999")
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
func (a *API) SendRegister(server string) {
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
	resp, err := http.Post(fmt.Sprintf("http://%v/register", server), "application/json", bytes.NewBuffer(b))
	errors.CheckError(err)

	msg := Message{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &msg)
	fmt.Printf("DEBUG: %v\n", msg.Value)
}
