package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamb/netpupper/controller"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/tcpbw"
	"github.com/adamb/netpupper/perf/udpr"
	"github.com/adamb/netpupper/scheduler"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type API struct {
	Controller *controller.Controller

	Config    *APIConfig
	TcpBwAddr string
	UdpAddr   string

	configFile string

	rw http.ResponseWriter
}
type APIConfig struct {
	Servers  []string
	ApiAddr  string
	Schedule *scheduler.TestSchedule
	Tags     map[string]string
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
func (a *API) Configure(cf string) bool {

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
		ApiAddr: ":8999",
	}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); err == nil {
		data, err := ioutil.ReadFile(serverFile)
		errors.CheckError(err)
		err = yaml.Unmarshal(data, a.Config)
		errors.CheckError(err)
		// Initialize the scheduler
		//a.Schedule = &sch

		return true

	}
	return false
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
	var result bool
	s := tcpbw.Server{}
	result = s.Configure(a.configFile)
	if result {
		fmt.Printf("Start TCPBW port: %v\n", s.Config.Address)
		go s.Run()
	}

	a.TcpBwAddr = s.Config.Address
	// UDP SERVER
	us := udpr.Server{}
	result = us.Configure(a.configFile)
	if result {
		fmt.Printf("Start UDPR port: %v\n", us.Config.Address)
		go us.Run()
	}
	a.UdpAddr = us.Config.Address

	fmt.Printf("Started API server on %v\n", a.Config.ApiAddr)
	a.StartServerApi(a.Config.ApiAddr)
}

func (a *API) StartServerApi(p string) {
	con := controller.Controller{}
	a.Controller = &con

	http.HandleFunc("/register", a.register)
	http.HandleFunc("/tcpbw", a.bwtest)
	http.HandleFunc("/udpr", a.udptest)
	http.HandleFunc("/clients", a.getclients)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v", p), nil))

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

	// ApiAddr will usually be a port number spec only, so this is buggy
	reg.ApiAddr = fmt.Sprintf("%v%v", reg.Addr, reg.ApiAddr)

	// If the tcpbwaddr string is an fqdn or IP:port combo use it directly
	// Otherwise get the hostname for this server and set it as that
	if len(strings.Split(a.TcpBwAddr, ":")) < 2 {
		host, _ := os.Hostname()
		reg.Server = fmt.Sprintf("%v%v", host, a.TcpBwAddr)
	} else {
		reg.Server = a.TcpBwAddr
	}

	// Same but for UDP tests
	if len(strings.Split(a.UdpAddr, ":")) < 2 {
		host, _ := os.Hostname()
		reg.UdpServer = fmt.Sprintf("%v%v", host, a.UdpAddr)
	} else {
		reg.UdpServer = a.UdpAddr
	}

	// Basic
	reg.ByteCount = "150M"
	reg.Rate = 1000
	reg.PacketCount = uint64(10000)

	sch := a.Config.Schedule
	sch.ScheduleTest(reg.StartbwTest)
	sch.ScheduleTest(reg.StartUdpTest)
	//a.Schedule.ScheduleTest(func() {})
	sch.PrintSchedule()

	go sch.Ticker()

	m := Message{
		Value: fmt.Sprintf("%v registered to controller.", reg.Addr),
	}
	b, err := json.Marshal(m)
	errors.CheckError(err)
	w.Write(b)

}

/*
Start a UDP reliability test from the given client
*/
func (a *API) udptest(w http.ResponseWriter, r *http.Request) {
	var cc udpr.ClientConfig
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
	fmt.Printf("DEBUG: Got a UDP test request destination: %v\n", cc.Server)
	ac.SetResponse(w)

	c := udpr.Client{}
	c.Configure(a.configFile)
	c.Config = &cc
	//c.SetTestCollector(&ac)
	c.Run()
}

func (a *API) bwtest(w http.ResponseWriter, r *http.Request) {
	var cc tcpbw.ClientConfig
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &cc)

	// Setup an API test collector. Stats will be sent back to the client.
	//ac := ApiCollector{}
	//enableCors(&w)

	// Required for CORS to work.
	if r.Method == "OPTIONS" {
		return
	}
	fmt.Printf("DEBUG: Got a test request destination: %v\n", cc.Server)

	c := tcpbw.Client{}
	c.Configure(a.configFile)
	c.Config = &cc
	//c.SetTestCollector(&ac)

	c.Run()
}

func (a *API) getclients(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	b, err := json.Marshal(a.Controller.Clients)
	errors.CheckError(err)
	w.Write(b)
}

/*
Send a register
*/
func (a *API) SendRegister(server string) {
	host, _ := os.Hostname()
	tags := a.Config.Tags
	tags["name"] = host
	reg := controller.Client{
		ApiAddr: a.Config.ApiAddr,
		Tags:    tags,
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
