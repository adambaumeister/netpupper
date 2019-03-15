package tcpbw

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type Runner interface {
	Configure()
	Run()
}

/*
Main struct for tcp bandwidth server
*/
type Server struct {
	Config *ServerConfig
}

type ServerConfig struct {
	Address string
}

/*
Configure the TCPBW Server
Returns TRUE if this method matches the requested config
*/
func (s *Server) Configure() {
	serverFile := "./server.yml"
	// First, try bootstrapping from the YAML server file
	s.Config = &ServerConfig{}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); os.IsExist(err) {
		data, err := ioutil.ReadFile("./server.yml")
		errors.CheckError(err)

		err = yaml.Unmarshal(data, s.Config)
		errors.CheckError(err)
	}
}

func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.Config.Address)
	if err != nil {
		errors.RaiseError("Failed to open socket.")
	}
	for {
		conn, err := ln.Accept()
		errors.CheckError(err)

		h := ReadHeader(conn)
		switch {
		case h.PacketType.Value == OPEN_TYPE:
			var o = Open{}
			o = ReadOpen(conn)
			fmt.Printf("Got a connection from: %v, Packet Type: %v Data to follow: %v bytes\n", conn.RemoteAddr(),
				h.PacketType.Value, o.DataLength)
			SendConfirm(conn)

			timedRead(conn, o.DataLength)
		}
	}
}

/*
TCPBW Client struct
*/
type Client struct {
	Config *clientConfig
}

type clientConfig struct {
	Server string
	Bytes  uint64
}

/*
Configure the TCPBW Client
Returns TRUE if a client mode is requested
*/
func (c *Client) Configure() {
	f := "./client.yml"
	// First, try bootstrapping from the YAML server file
	c.Config = &clientConfig{}
	// If the yaml file exists
	if _, err := os.Stat(f); os.IsExist(err) {
		data, err := ioutil.ReadFile("./server.yml")
		errors.CheckError(err)

		err = yaml.Unmarshal(data, c.Config)
		errors.CheckError(err)
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp", c.Config.Server)
	if err != nil {
		errors.RaiseError("Failed to open connection!")
	}
	fmt.Printf("Succesfully connected to: %v\n", conn.RemoteAddr())

	// Send the open message, request to start
	SendOpen(conn)
	// Wait for a confirmation
	h := ReadHeader(conn)
	switch {
	case h.PacketType.Value == CONFIRM_TYPE:
		fmt.Printf("OPEN Request confirmed. Sending data...\n")
		// Test by splitting up the data
		conn.Write([]byte{255, 255})
		time.Sleep(2 * time.Second)
		conn.Write([]byte{255, 255})
		time.Sleep(1 * time.Second)
	}
}

func timedRead(conn net.Conn, rl uint64) {
	start := time.Now().UnixNano()

	// Chunk size is how much we read at each time interval
	chunk := rl / 4
	data := make([]byte, rl)

	currentChunk := 1
	// Read each chunk until we've read the entire thing
	for currentChunk <= 4 {
		chunkData := make([]byte, chunk)
		conn.Read(chunkData)
		// Append each read chunk to the full data array
		data = append(data, chunkData...)
		currentChunk = currentChunk + 1
	}
	t := time.Now().UnixNano()
	elapsed := t - start
	fmt.Printf("Read took %v ns\n", elapsed)

}
