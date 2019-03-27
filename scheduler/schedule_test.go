package scheduler

import (
	"fmt"
	"testing"
	"time"
)

type dummyTD struct {
}

func (d dummyTD) Run() {
	fmt.Printf("Ran at %v\n", time.Now().Unix())
	return
}
func TestTestSchedule_ScheduleTest(t *testing.T) {
	td := dummyTD{}
	ts := InitSchedule("none")
	ts.Interval = 30
	ts.Buffer = 5
	ts.ScheduleTest(td)
	time.Sleep(1 * time.Second)
	ts.ScheduleTest(td)

	ts.PrintSchedule()
	go ts.Ticker()
	time.Sleep(10 * time.Second)

}
