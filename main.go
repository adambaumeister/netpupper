package main

import (
	"flag"
	"fmt"
	"github.com/adamb/netpupper/api"
	"github.com/adamb/netpupper/perf/tcpbw"
	"github.com/adamb/netpupper/perf/udpr"
)

type Runner interface {
	Configure(string) bool
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
	testType := flag.String("test", "TCP", "Run either a TCP Bandwidth test or a UDP reliability test (default: TCP)")

	// CMDLINE Args
	addr := flag.String("address", ":8080", "Address to bind server daemon OR server address to test against.")
	bytes := flag.String("bytes", "1G", "Total bytes to send.")

	flag.Parse()
	if *clientMode {
		fmt.Printf("Started in CLIENT mode\n")
		if *testType == "UDP" {
			c := &udpr.Client{}
			c.Configure(*cfgFile)
			c.Config.Server = *addr
			c.Config.Rate = 5000
			c.Config.PacketCount = uint64(20000)
			return c
		} else {
			c := &tcpbw.Client{}
			c.Configure(*cfgFile)
			c.Config.Server = *addr
			c.Config.Bytes = *bytes
			c.Config.Reverse = *reverse
			return c
		}
	} else if *serverMode {
		fmt.Printf("Started in SERVER mode\n")
		if *testType == "UDP" {
			s := &udpr.Server{}
			s.Configure(*cfgFile)
			s.Config.Address = *addr
			return s
		} else {
			s := &tcpbw.Server{}
			s.Configure(*cfgFile)
			s.Config.Address = *addr
			return s
		}
	} else if *daemon {
		a := &api.API{}
		a.Configure(*cfgFile)
		fmt.Printf("Started in Daemon mode.\n")
		return a
	} else {
		return nil
	}
}
