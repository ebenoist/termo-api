package main

import (
	"database/sql"
	"log"
	"time"
)

type Schedule struct {
	Id         sql.NullInt64 `db:"id"`
	Hour       int           `db:"hour"`
	TargetTemp int           `db:"target_temp"`
	Days       string        `db:"days"`
}

const (
	WEEKENDS = "weekends"
	WORKDAYS = "workdays"
)

func FindScheduleByTime(d *Database, t time.Time) *Schedule {
	connection := d.Connection()

	hour := t.Hour()
	var days string
	if t.Weekday() == 0 || t.Weekday() == 6 {
		days = WEEKENDS
	} else {
		days = WORKDAYS
	}

	var result Schedule
	connection.Get(&result, `SELECT * FROM schedule
		WHERE days=$1 AND hour<=$2
		ORDER BY hour DESC LIMIT 1`, days, hour)
	return &result
}

func (s *Schedule) Save(d *Database) error {
	connection := d.Connection()
	result, err := connection.NamedExec(`
		INSERT OR REPLACE INTO
		schedule (hour, target_temp, days)
		values (:hour, :target_temp, :days)`, s)

	id, _ := result.LastInsertId()
	s.Id.Scan(id)

	if err != nil {
		log.Printf("db: error %s", err)
	}

	return err
}
