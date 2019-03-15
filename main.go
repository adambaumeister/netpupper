package main

import (
	"flag"
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
	serverMode := flag.Bool("server", false, "a bool")
	if *serverMode {
		s := &tcpbw.Server{}
		return s
	} else {
		c := &tcpbw.Client{}
		return c
	}
}
