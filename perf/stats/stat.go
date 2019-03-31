package stats

import (
	"fmt"
	"time"
)

type Collector interface {
	WriteBwTest(BpsResult)
	WriteSummary(SummaryResult)
}

type Test struct {
	StartTime int64
	Queue     []BwTestResult
	InMsgs    chan BwTestResult
	InReqs    chan string

	Stop    chan bool
	EndTime int64

	Collector Collector
}

type IntTestResult interface {
	Get() uint64
}

type BwTestResult struct {
	Bytes   int
	Elapsed uint64
}

type BpsResult struct {
	Bps float64
}

type SummaryResult struct {
	Bps float64
}

func (b *BpsResult) Get() uint64 {
	return uint64(b.Bps)
}

func InitTest() *Test {
	t := Test{}
	t.Collector = &StdoutCollector{}

	t.Queue = []BwTestResult{}
	t.InMsgs = make(chan BwTestResult)
	t.InReqs = make(chan string)
	t.Stop = make(chan bool)
	t.StartTime = time.Now().UnixNano()
	go t.Listen()
	go t.IntervalReport(1)
	return &t
}

/*
Parse the Current() results every i time interval (seconds)
Sends a call to the channel that listens for reques
*/
func (t *Test) IntervalReport(i int) {
	// While the test is still running
	for t.EndTime == 0 {
		select {
		// Stop on a Stop signal
		case <-t.Stop:
			fmt.Printf("Periodic reporting stopped\n")
			return
		default:
			time.Sleep(time.Duration(i) * time.Second)
			// If the test still aint over
			if t.EndTime == 0 {
				t.InReqs <- "stats"
			}
		}
	}
}

func (t *Test) Current() {
	// total bytes
	tb := 0
	// Elapsed time
	et := uint64(0)
	for _, tr := range t.Queue {
		tb = tb + tr.Bytes
		et = et + tr.Elapsed
	}

	// divide the total bytes by the total elapsed time
	cbps := float64(tb) / float64(et)
	// Convert from nano to reg seconds
	cbps = cbps * 1000000000
	tr := BpsResult{
		Bps: cbps,
	}
	t.Collector.WriteBwTest(tr)

}

func (t *Test) Summary() {
	tb := 0
	for _, tr := range t.Queue {
		tb = tb + tr.Bytes
	}
	e := uint64(t.EndTime - t.StartTime)
	cbps := float64(tb) / float64(e)
	// Convert from nano to reg seconds
	cbps = cbps * 1000000000
	//fmt.Printf("Test summary (len: %v): %v\n", len(t.Queue), perf.ByteToString(uint64(cbps)))
	sr := SummaryResult{
		Bps: cbps,
	}
	t.Collector.WriteSummary(sr)
}

func (t *Test) End() {
	fmt.Printf("Client end\n")
	t.EndTime = time.Now().UnixNano()
	t.Stop <- true
}

/*
Listen to incoming test results
Also listen for incoming requests from the other programs for stats details
*/
func (t *Test) Listen() {
	for {
		select {
		case <-t.InReqs:
			t.Current()
		case m := <-t.InMsgs:
			t.AddResult(m)
		case <-t.Stop:
			fmt.Printf("Test finished.\n")
			return
		}
	}
}

func (t *Test) AddResult(tr BwTestResult) {
	t.Queue = append(t.Queue, tr)
}
