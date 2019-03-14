package main

import (
	"github.com/adamb/netpupper/perf/server"
)

func main() {
	go server.Server()
	server.Client()

}
