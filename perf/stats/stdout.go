package stats

import (
	"fmt"
	"github.com/adamb/netpupper/perf"
)

type StdoutCollector struct {
}

func (w *StdoutCollector) WriteBwTest(r BpsResult) {
	fmt.Printf("CURRENT : %v\n", perf.ByteToString(uint64(r.Get())))
}
