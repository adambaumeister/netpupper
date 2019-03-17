package stats

import (
	"fmt"
	"github.com/adamb/netpupper/perf"
	"time"
)

type TestResults struct {
	StartTime int64
	Queue     []TestResult
	InMsgs    chan TestResult
	EndTime   int64
}

type TestResult struct {
	Bytes   int
	Elapsed uint64
}

func InitTest() *TestResults {
	t := TestResults{}
	t.Queue = []TestResult{}
	t.InMsgs = make(chan TestResult)

	t.StartTime = time.Now().UnixNano()
	go t.Listen()
	return &t
}

func (t *TestResults) Summary() {
	tb := 0
	for _, tr := range t.Queue {
		tb = tb + tr.Bytes
	}
	e := uint64(t.EndTime - t.StartTime)
	cbps := float64(tb) / float64(e)
	// Convert from nano to reg seconds
	cbps = cbps * 1000000000
	fmt.Printf("Test summary: %v\n", perf.ByteToString(uint64(cbps)))
}

func (t *TestResults) End() {
	t.EndTime = time.Now().UnixNano()
}

func (t *TestResults) Listen() {
	for {
		t.AddResult(<-t.InMsgs)
	}
}

func (t *TestResults) AddResult(tr TestResult) {
	t.Queue = append(t.Queue, tr)
}
