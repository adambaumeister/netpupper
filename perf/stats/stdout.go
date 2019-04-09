package stats

import (
	"fmt"
	"github.com/adamb/netpupper/perf"
)

type StdoutCollector struct {
}

func (w *StdoutCollector) WriteBwTest(r BpsResult) {
	fmt.Printf("%v\n", perf.ByteToString(uint64(r.Get())))
}

func (w *StdoutCollector) WriteBwSummary(r BpsSummaryResult) {
	fmt.Printf("SUMMARY: %v\n", perf.ByteToString(uint64(r.Bps)))
}

func (w *StdoutCollector) WriteReliabilityTest(r ReliabilityResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
func (w *StdoutCollector) WriteReliabilitySummary(r ReliabilitySummaryResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
