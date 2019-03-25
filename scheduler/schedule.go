package scheduler

import "time"

/*
TestSchedule builds a queue of tests
Tests have a scheduled interval and tests must leave
*/
type TestSchedule struct {
	Tests []*ScheduleItem

	interval int
	buffer   int
}

// Init the schedule and build the slot layout
func InitSchedule(cf string) TestSchedule {
	ts := TestSchedule{}
	return ts
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
	currentSec := time.Now().Second() + (time.Now().Minute() * 60) + (time.Now().Hour() * 3600)
	for _, si := range t.Tests {
		st := si.Time
		if currentSec >= st+t.buffer {

		}
	}
}

/*
ScheduleItem represents a single test run and the time at which it's scheduled
*/
type ScheduleItem struct {
	Name string
	Time int

	lastRun int64
	test    TestDefinition
}

/*
Test definition is the actual test to run.
*/
type TestDefinition interface {
	Run()
}
