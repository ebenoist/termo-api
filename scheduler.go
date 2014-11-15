package main

import (
	"log"
	"time"
)

type Scheduler struct {
	database   *Database
	thermostat *Thermostat
	schedule   *Schedule
}

// Scheduler will on a 2 minute delay search for a current
// schedule. If one is found for the current time frame, a
// it is compared against the last found schedule and if it
// differs, it should set the current thermostat's TargetTemp
func (s *Scheduler) Run(t *Thermostat) {
	database := &Database{}

	for {
		found := FindScheduleByTime(database, time.Now())
		if found.TargetTemp != 0 && found != s.schedule {
			log.Printf("Found new schedule %q", found)
			s.schedule = found
			s.thermostat.TargetTemp = found.TargetTemp
		}

		time.Sleep(2 * time.Minute)
	}
}
