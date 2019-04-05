package udpr

import (
	"encoding/binary"
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
	Config     *ServerConfig `yaml:"tcpbw"`
}

type ServerConfig struct {
	Address string
}

/*
Configure the UDP Reliability test
Returns TRUE if this method matches the requested config
*/
func (s *Server) Configure(cf string) {
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
		fmt.Printf("yep got here: %v\n", s.Config.Address)
		errors.CheckError(err)
	}
}

func (s *Server) Run() {
	// Define channels
	s.stopChan = make(chan bool)
	s.notifyChan = make(chan bool)

	addr, err := net.ResolveUDPAddr("udp", s.Config.Address)
	conn, err := net.ListenUDP("udp", addr)
	errors.CheckError(err)
	for {
		packet := make([]byte, 1500)
		// Read the max number of bytes in a datagram(1500) into a variable length slice of bytes, 'Buffer'
		// Also set the total number of bytes read so we can check it later
		n, addr, _ := conn.ReadFromUDP(packet)
		fmt.Printf("Got %v bytes from %v.", n, addr.IP)
	}
}

type Client struct {
	notifyChan    chan bool
	stopChan      chan bool
	testCollector stats.Collector

	Config *ClientConfig
}

type ClientConfig struct {
	Server string
}

/*
Configure the UDPR Client
*/
func (c *Client) Configure(cf string) {
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
		data, err := ioutil.ReadFile(f)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, c.Config)
		errors.CheckError(err)
	}
}

func (c *Client) Run() {
	c.stopChan = make(chan bool)
	c.notifyChan = make(chan bool)

	conn, err := net.Dial("udp", c.Config.Server)
	errors.CheckError(err)
	b := []byte{255, 255}
	err = binary.Write(conn, binary.BigEndian, b)
}
