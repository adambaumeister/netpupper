package stats

import (
	"fmt"
	"testing"
)

func TestInflux_Configure(t *testing.T) {
	ic := Influx{}
	ic.Configure("C:/users/adam/go/src/github.com/adamb/netpupper/daemon_server.yml")
}

func TestInflux(t *testing.T) {
	ic := Influx{}
	ic.Configure("C:/users/adam/go/src/github.com/adamb/netpupper/daemon_server.yml")
	r := BpsResult{
		Bps: 10000000,
	}
	fmt.Printf("%v\n", ic.Config.HTTPConfig.Addr)
	ic.WriteBwTest(r)
	// This method does not require the sum result
	ic.WriteBwSummary(BpsSummaryResult{})
}
