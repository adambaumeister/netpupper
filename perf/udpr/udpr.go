package udpr

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/stats"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
)

type Server struct {
	notifyChan chan bool
	stopChan   chan bool
	Config     *ServerConfig `yaml:"udpr"`
	Influx     *stats.Influx `yaml:"influx"`

	testCollector stats.Collector
}

type ServerConfig struct {
	Address string
}

/*
Configure the UDP Reliability test
Returns TRUE if this method matches the requested config
*/
func (s *Server) Configure(cf string) bool {
	var serverFile string
	if len(cf) > 0 {
		serverFile = cf
	} else if len(os.Getenv("NETP_CONFIG")) > 0 {
		serverFile = os.Getenv("NETP_CONFIG")
	}
	// First, try bootstrapping from the YAML server file
	s.Config = &ServerConfig{}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); err == nil {
		data, err := ioutil.ReadFile(serverFile)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, s)
		errors.CheckError(err)

		if s.Config.Address != "" {
			return true
		}
	}
	return false
}

func (s *Server) Run() {
	// Define channels
	s.stopChan = make(chan bool)
	s.notifyChan = make(chan bool)

	addr, err := net.ResolveUDPAddr("udp", s.Config.Address)
	errors.CheckError(err)
	for {
		conn, err := net.ListenUDP("udp", addr)
		fmt.Printf("Waiting for UDP client.\n")
		packet := make([]byte, 1500)
		_, addr, err := conn.ReadFromUDP(packet)

		errors.CheckError(err)
		h := ReadHeader(packet)
		switch {
		case h.PacketType.Value == OPEN_TYPE:
			var o Open
			o = ReadOpen(packet)
			fmt.Printf("Got UDP Open from %v : %v : %v\n.", addr.IP, o.DataLength, o.AckCount)
			SendConfirm(conn, addr)
			ut := InitUdpSm(conn, addr, o.AckCount, o.DataLength)
			test := stats.InitTest()
			ut.countedRead(conn, test)
			test.End()
			conn.Close()
		}
	}
}

type Client struct {
	notifyChan    chan bool
	stopChan      chan bool
	testCollector stats.Collector

	Config *ClientConfig `yaml:"udpr"`
	Influx *stats.Influx `yaml:"influx"`
	Tags   map[string]string
}

type ClientConfig struct {
	Server      string
	PacketCount uint64
	Rate        int
}

/*
Configure the UDPR Client
*/
func (c *Client) Configure(cf string) bool {
	var f string
	if len(cf) > 0 {
		f = cf
	} else if len(os.Getenv("NETP_CONFIG")) > 0 {
		f = os.Getenv("NETP_CONFIG")
	}
	// First, try bootstrapping from the YAML server file
	c.Config = &ClientConfig{}
	// If the yaml file exists
	if _, err := os.Stat(cf); err == nil {
		fmt.Printf("Yep got here..\n")
		data, err := ioutil.ReadFile(f)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, c)
		errors.CheckError(err)

		if c.Influx != nil {

			host, _ := os.Hostname()
			c.Tags["name"] = host
			c.Influx.Tags = c.Tags
			c.testCollector = c.Influx
		}

		return true
	}
	return false
}

func (c *Client) Run() {
	c.stopChan = make(chan bool)
	c.notifyChan = make(chan bool)

	conn, err := net.Dial("udp", c.Config.Server)
	errors.CheckError(err)
	SendOpen(conn, c.Config.PacketCount, 1000)
	packet := make([]byte, 1500)
	_, err = conn.Read(packet)
	errors.CheckError(err)

	h := ReadHeader(packet)
	switch {
	case h.PacketType.Value == CONFIRM_TYPE:
		fmt.Printf("UDP stream confirmed.\n")
		// Problem is here - test collector is reinitlized
		test := stats.InitTest()
		if c.testCollector != nil {
			fmt.Printf("Using INFLUXDB as stats collector for UDP...%v\n", c.Influx.Database)
			test.Collector = c.testCollector
		}
		addr, _ := net.ResolveUDPAddr("udp", c.Config.Server)

		// start the state machine
		ut := InitUdpSm(conn, addr, 1000, c.Config.PacketCount)
		ut.countedSend(test, c.Config.Rate)
		test.End()
		test.Summary()
	}
}

func (c *Client) SetTestCollector(sc stats.Collector) {
	c.testCollector = sc
}
