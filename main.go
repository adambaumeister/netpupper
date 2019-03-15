package main

import (
	"flag"
	"fmt"
	"github.com/adamb/netpupper/perf/tcpbw"
)

func main() {
	r := ParseArgs()
	r.Run()
}

/*
Parse the cmdline args to determine what mode to run in
*/
func ParseArgs() tcpbw.Runner {
	serverMode := flag.Bool("server", false, "Run as netpupper server.")
	clientMode := flag.Bool("client", false, "Run as netpupper client.")

	s := &tcpbw.Server{}
	s.Configure()
	c := &tcpbw.Client{}
	c.Configure()
	// CMDLINE Args
	addr := flag.String("address", ":8080", "Address to bind server daemon OR server address to test against.")

	flag.Parse()
	s.Config.Address = *addr
	c.Config.Server = *addr

	if *clientMode {
		fmt.Printf("Started in CLIENT mode")
		return c
	} else if *serverMode {
		fmt.Printf("Started in SERVER mode")
		return s
	} else {
		return s
	}
}
