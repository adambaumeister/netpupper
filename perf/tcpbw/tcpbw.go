package tcpbw

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf"
	"github.com/adamb/netpupper/perf/stats"
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
	notifyChan: bool channel, notifies the server to do "something", depending on the context
	stopChan: Bool channel, notifies the receiving function to stop.
*/
type Server struct {
	notifyChan chan bool
	stopChan   chan bool
	Config     *ServerConfig
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
	// Define channels
	s.stopChan = make(chan bool)
	s.notifyChan = make(chan bool)

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
			if o.Reverse == 0 {
				s.ClientToServer(conn, &o)
			} else {
				s.ServerToClient(conn, &o)
			}
		}
	}
}
func (s *Server) ClientToServer(conn net.Conn, o *Open) {

	fmt.Printf("Got a connection from: %v, Data to follow: %v bytes\n", conn.RemoteAddr(), o.DataLength)
	SendConfirm(conn)

	// Initilize a test for storing the results
	test := stats.InitTest()
	go timedRead(conn, o.DataLength, test.InMsgs, s.stopChan)
	// Schedule the test interval function
	s.GetUserInput(test.InReqs)
	// End the test and print the summary
	test.End()
	conn.Close()
	test.Summary()
}
func (s *Server) ServerToClient(conn net.Conn, o *Open) {
	fmt.Printf("Got a connection from: %v, Data to follow: %v bytes REVERSE\n", conn.RemoteAddr(), o.DataLength)
	SendConfirm(conn)
	// Initilize a test for storing the results
	test := stats.InitTest()
	go timedSend(conn, o.DataLength, test.InMsgs, s.stopChan)
	// Schedule the test interval function
	s.GetUserInput(test.InReqs)
	// End the test and print the summary
	test.End()
	conn.Close()
	test.Summary()
}

/*
GetUserInput: Retrieves user input from stdin and sends to the provided channel based on said input.
*/
func (s *Server) GetUserInput(c chan string) {
	// Temp channel
	tc := make(chan string)
	reader := bufio.NewScanner(os.Stdin)
	for {
		// Non-blocking call to wait for user input
		go func() {
			reader.Scan()
			s := reader.Text()

			//s = strings.TrimSuffix(s, "\n")
			//s = strings.TrimSuffix(s, "\r")
			tc <- s
			return
		}()

		// This part here WILL block.
		var str string
		select {
		// If we get user input
		case str = <-tc:
			switch str {
			case "stats":
				fmt.Printf("Stats requested.\n")
				c <- str
			}
		// If we get signalled to stop
		case <-s.stopChan:
			return
		}
	}
}

/*
TCPBW Client struct
*/
type Client struct {
	notifyChan chan bool
	stopChan   chan bool
	Config     *clientConfig
}

type clientConfig struct {
	Server  string
	Bytes   string
	Reverse bool
}

/*
Configure the TCPBW Client
Configuration can set external to this function, however, should be done after this call.
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
	c.stopChan = make(chan bool)
	c.notifyChan = make(chan bool)

	dl := uint64(perf.StringToByte(c.Config.Bytes))
	if dl < perf.MEGABYTE {
		errors.RaiseError("Error: passed byte count too small!")
	}
	conn, err := net.Dial("tcp", c.Config.Server)
	if err != nil {
		errors.RaiseError("Failed to open connection!")
	}
	fmt.Printf("Succesfully connected to: %v\n", conn.RemoteAddr())

	// Send the open message, request to start
	if c.Config.Reverse {
		SendOpen(conn, dl, 1)
	} else {
		SendOpen(conn, dl, 0)
	}

	// Wait for a confirmation
	h := ReadHeader(conn)
	switch {
	case h.PacketType.Value == CONFIRM_TYPE:
		if c.Config.Reverse {
			fmt.Printf("OPEN Request for reverse mode confirmed. Receiving data %v...\n", dl)
			// Initilize a test for storing the results
			test := stats.InitTest()
			timedRead(conn, dl, test.InMsgs, c.stopChan)

			test.End()
			test.Summary()
			return
		} else {
			fmt.Printf("OPEN Request confirmed. Sending data...\n")
			// Initilize a test for storing the results
			test := stats.InitTest()
			timedSend(conn, dl, test.InMsgs, c.stopChan)
			test.End()
			test.Summary()
			return
		}

	}
}

/*
TimedRead implements a timed read of rl bytes from conn net.Conn
It uses a chan (NC: NotifyChannel) to send stats.
It also uses a chan to signal the running function that it's done
To allow this to be a a part of the read flow, this method splits the receipt of data into discrete chunks.
*/
func timedRead(conn net.Conn, rl uint64, nc chan stats.BwTestResult, sc chan bool) {
	start := time.Now().UnixNano()
	fmt.Printf("Read start: %v\n", start)
	lt := start
	// Chunk size is how much we read at each time interval
	chunk := perf.MEGABYTE

	currentChunk := uint64(0)
	// Read each chunk until we've read the entire thing
	for currentChunk < rl {
		chunkData := make([]byte, chunk)

		rc := make(chan bool)
		go basicRead(conn, chunkData, rc)

		var e uint64
		select {
		case _ = <-rc:
			// Read the current time
			t := time.Now().UnixNano()
			// Get the elapsed from the last chunk time
			e = uint64(t - lt)

			if e > 0 {
				// Bytes transferred per nanosecond
				tr := stats.BwTestResult{
					Bytes:   chunk,
					Elapsed: e,
				}
				// Send the result to the given notify channel as type stats.TestResuly
				nc <- tr
			}
			//fmt.Printf("Read chunk at %v BPS\n", ByteToString(uint64(cbps)))
			// Set the last time to the time of this chunk's finished read
			lt = t
		}

		// Append each read chunk to the full data array
		// Don't do this, obviously it fills ya memory up fam
		//data = append(data, chunkData...)
		currentChunk = currentChunk + uint64(chunk)

	}
	// THIS IS BROKEN IN REVERSE FOR SOME REASON
	//sc <- true
	fmt.Printf("read end\n\n")
	return
}

/*
TimedSend implements a timed send method from Conn of sl bytes.
It uses a chan (NC: NotifyChannel) to listen for requests for stats.
It also uses a chan to signal the running function that it's done
To allow this to be a a part of the read flow, this method splits the receipt of data into discrete chunks.
*/
func timedSend(conn net.Conn, sl uint64, nc chan stats.BwTestResult, sc chan bool) {
	start := time.Now().UnixNano()
	lt := start
	// Chunk size is how much we read at each time interval
	chunk := perf.MEGABYTE

	currentChunk := uint64(0)
	// Read each chunk until we've read the entire thing
	var e uint64
	var cbps float64

	for currentChunk < sl {
		chunkData := make([]byte, chunk)
		rand.Read(chunkData)
		rc := make(chan bool)
		go basicWrite(conn, chunkData, rc)
		select {
		case _ = <-rc:
			// Read the current time
			t := time.Now().UnixNano()
			// Get the elapsed from the last chunk time
			e = uint64(t - lt)
			if e > 0 {
				// Bytes transferred per nanosecond
				cbps = float64(chunk) / float64(e)
				// Convert from nano to reg seconds
				cbps = cbps * 1000000000
				//fmt.Printf("time: %v, BPS: %v, chunk: %v\n", e, int(cbps), int(chunk))
				// Bytes transferred per nanosecond
				tr := stats.BwTestResult{
					Bytes:   chunk,
					Elapsed: e,
				}
				// Send the result to the given notify channel as type stats.TestResuly
				nc <- tr
			}
			//fmt.Printf("Read chunk at %v BPS\n", ByteToString(uint64(cbps)))
			// Set the last time to the time of this chunk's finished read
			lt = t
		}

		// Append each read chunk to the full data array
		// Don't do this, obviously it fills ya data up fam
		//data = append(data, chunkData...)
		currentChunk = currentChunk + uint64(chunk)
	}
	t := time.Now().UnixNano()
	e = uint64(t - start)
	cbps = float64(sl) / float64(e)
	// Convert from nano to reg seconds
	cbps = cbps * 1000000000
	// TX and RX will be slightly different as the timing here does not include the time the last chunk is read.

	sc <- true
	fmt.Printf("TX: %v Bps\n", perf.ByteToString(uint64(cbps)))
	return
}

// Both methods below use a channel to indicate their status
func basicRead(conn net.Conn, data []byte, c chan bool) {
	conn.Read(data)
	c <- true
	return
}
func basicWrite(conn net.Conn, data []byte, c chan bool) {
	conn.Write(data)
	c <- true
	return
}
