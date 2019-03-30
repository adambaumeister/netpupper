package main

import (
	"flag"
	"fmt"
	"github.com/adamb/netpupper/api"
	"github.com/adamb/netpupper/perf/tcpbw"
)

type Runner interface {
	Configure(string)
	Run()
}

func main() {
	r := ParseArgs()
	r.Run()
}

/*
Parse the cmdline args to determine what mode to run in
*/
func ParseArgs() Runner {
	serverMode := flag.Bool("server", false, "Run as netpupper server.")
	clientMode := flag.Bool("client", false, "Run as netpupper client.")
	reverse := flag.Bool("reverse", false, "Reverse data direction.")
	cfgFile := flag.String("config", "./daemon_server.yml", "YAML Configuration file")
	daemon := flag.Bool("daemon", false, "Run in DAEMON mode and start the API.")

	// CMDLINE Args
	addr := flag.String("address", ":8080", "Address to bind server daemon OR server address to test against.")
	bytes := flag.String("bytes", "1G", "Total bytes to send.")

	flag.Parse()
	s := &tcpbw.Server{}
	s.Configure(*cfgFile)
	c := &tcpbw.Client{}
	c.Configure(*cfgFile)

	a := &api.API{}
	a.Configure(*cfgFile)
	s.Config.Address = *addr
	c.Config.Server = *addr
	c.Config.Bytes = *bytes
	c.Config.Reverse = *reverse

	if *clientMode {
		fmt.Printf("Started in CLIENT mode\n")
		return c
	} else if *serverMode {
		fmt.Printf("Started in SERVER mode\n")
		return s
	} else if *daemon {
		fmt.Printf("Started in Daemon mode.\n")
		return a
	} else {
		return s
	}
}
