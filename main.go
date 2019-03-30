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

	if *clientMode {
		fmt.Printf("Started in CLIENT mode\n")
		c := &tcpbw.Client{}
		c.Configure(*cfgFile)
		c.Config.Server = *addr
		c.Config.Bytes = *bytes
		c.Config.Reverse = *reverse
		return c
	} else if *serverMode {
		s := &tcpbw.Server{}
		s.Configure(*cfgFile)
		s.Config.Address = *addr
		fmt.Printf("Started in SERVER mode\n")
		return s
	} else if *daemon {
		a := &api.API{}
		a.Configure(*cfgFile)
		fmt.Printf("Started in Daemon mode.\n")
		return a
	} else {
		return nil
	}
}
