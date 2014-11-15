package main

import (
	"log"
	"time"
)

// Scheduler will on a 2 minute delay search for a current
// schedule. If one is found for the current time frame, a
// it is compared against the last found schedule and if it
// differs, it should set the current thermostat's TargetTemp
func SchedulerRun(t *Thermostat) {
	database := &Database{}
	currentSchedule := &Schedule{}

	for {
		found := FindScheduleByTime(database, time.Now())
		if found.TargetTemp != 0 && found.TargetTemp != currentSchedule.TargetTemp {
			log.Printf("Found new schedule %q", found)
			currentSchedule = found
			t.TargetTemp = found.TargetTemp
		}

		time.Sleep(2 * time.Minute)
	}
}
