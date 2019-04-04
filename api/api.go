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

	Config     *APIConfig
	configFile string

	rw http.ResponseWriter
}
type APIConfig struct {
	Servers []string
	ApiPort string
	Tags    []controller.Tag
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
	(*w).Header().Set("Content-type", "application/json")
}

/*
Configure the API
Returns TRUE if this method matches the requested config
*/
func (a *API) Configure(cf string) {

	a.configFile = cf
	var serverFile string
	if len(cf) > 0 {
		serverFile = cf
	} else if len(os.Getenv("NETP_CONFIG")) > 0 {
		serverFile = os.Getenv("NETP_CONFIG")
	}
	fmt.Printf("Got yaml configuration %v\n", cf)
	// First, try bootstrapping from the YAML server file
	// Defaults below
	a.Config = &APIConfig{
		ApiPort: "8999",
	}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); err == nil {
		data, err := ioutil.ReadFile(serverFile)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, a.Config)
		errors.CheckError(err)
	}
}

/*
Start the API listener.
*/
func (a *API) Run() {
	// Register to the provided list of servers, if any
	for _, s := range a.Config.Servers {
		a.SendRegister(s)
	}
	// Start the test listeners
	// TCP BANDWIDTH SERVER
	s := tcpbw.Server{}
	s.Configure(a.configFile)
	fmt.Printf("Start TCPBW port: %v\n", s.Config.Address)
	go s.Run()
	fmt.Printf("Started API server on %v\n", a.Config.ApiPort)
	StartServerApi(a.Config.ApiPort)
}

func StartServerApi(p string) API {
	con := controller.Controller{}
	a := API{
		Controller: &con,
	}

	http.HandleFunc("/register", a.register)
	http.HandleFunc("/tcpbw", a.bwtest)
	http.HandleFunc("/clients", a.getclients)
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
	reg.Addr = strings.Split(r.RemoteAddr, ":")[0]
	a.Controller.AddClient(reg)

	m := Message{
		Value: fmt.Sprintf("%v registered to controller.", reg.Addr),
	}
	b, err := json.Marshal(m)
	errors.CheckError(err)
	w.Write(b)

}

func (a *API) bwtest(w http.ResponseWriter, r *http.Request) {
	var cc tcpbw.ClientConfig
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &cc)

	// Setup an API test collector. Stats will be sent back to the client.
	ac := ApiCollector{}
	enableCors(&w)

	// Required for CORS to work.
	if r.Method == "OPTIONS" {
		return
	}
	fmt.Printf("DEBUG: Got a test request destination: %v\n", cc.Server)
	ac.SetResponse(w)

	c := tcpbw.Client{
		Config: &cc,
	}
	c.SetTestCollector(&ac)
	c.Run()
}

func (a *API) getclients(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	b, err := json.Marshal(a.Controller.Clients)
	errors.CheckError(err)
	w.Write(b)
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
	tags := append(a.Config.Tags, gt)
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
