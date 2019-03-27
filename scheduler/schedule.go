package scheduler

import (
	"fmt"
	"time"
)

/*
TickerControl
Used to control the ticker via the Channel.
*/
const STOP = 0

type TickerControl struct {
	Type int
}

/*
Ticker
Runs the scheduler
*/
func (t *TestSchedule) Ticker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case tc := <-t.TickControl:
			if tc.Type == STOP {
				ticker.Stop()
				return
			}
		case <-ticker.C:
			t.RunTests(time.Now().Unix())
		}
	}
}

/*
TestSchedule builds a queue of tests
Tests have a scheduled interval and tests must leave
*/
type TestSchedule struct {
	Tests       []*ScheduleItem
	TickControl chan TickerControl

	Interval int
	Buffer   int
}

// Init the schedule and build the slot layout
func InitSchedule(cf string) TestSchedule {
	ts := TestSchedule{}
	return ts
}

// Run all tests at this current tick
// 		ti: Tick interval
func (t *TestSchedule) RunTests(ti int64) {
	for _, si := range t.Tests {
		if si.Time <= int(ti) {
			si.test.Run()
			// Update the schedule and lastrun
			si.Time = si.Time + t.Interval
			si.lastRun = si.Time
		}
	}
}

// Add a test to the schedule
/*
Ex. buffer is 60 seconds, interval is 2 hours
Test1 is scheduled at 9:00am
 first: 32400 seconds
 second: 39600
 etc.

Test2 is scheduled at 10am
 first: 36000

check if 36000 is greather than 32460 (test1 sched + buffer)
if true, schedule for that time.

*/

func (t *TestSchedule) ScheduleTest(td TestDefinition) {
	// Get the current second of the day
	currentSec := int(time.Now().Unix())
	// If there are no other tests scheduled for just now,
	if t.TestConflicts(currentSec) {
		si := &ScheduleItem{
			Time: currentSec,
			test: td,
		}
		t.Tests = append(t.Tests, si)
		currentSec = currentSec + t.Interval
		// if there are conflicts
	} else {
		for t.TestConflicts(currentSec) == false {
			currentSec = currentSec + t.Buffer
		}
		si := &ScheduleItem{
			Time: currentSec,
			test: td,
		}
		t.Tests = append(t.Tests, si)
		currentSec = currentSec + t.Interval
	}
}

/*
Check if a proposed time conflicts with any existing scheduled items
All checks are run against the offset
In a way this means tests
*/
func (t *TestSchedule) TestConflicts(secs int) bool {
	for _, si := range t.Tests {
		if secs > si.Time && secs <= si.Time+t.Buffer {
			return false
		}
	}
	return true
}

func (t *TestSchedule) PrintSchedule() {
	for _, si := range t.Tests {
		fmt.Printf("%v\n", si.Time)
	}
}

/*
ScheduleItem represents a single test run and the time at which it's scheduled
*/
type ScheduleItem struct {
	Name string
	Time int

	lastRun int
	test    TestDefinition
}

/*
Test definition is the actual test to run.
*/
type TestDefinition interface {
	Run()
}
