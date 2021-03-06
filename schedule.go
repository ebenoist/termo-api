package main

import (
	"github.com/guregu/null"
	"log"
	"time"
)

type Schedule struct {
	Id         null.Int `db:"id" json:"id"`
	Hour       int      `db:"hour" json:"hour" binding:"required"`
	TargetTemp int      `db:"target_temp" json:"targetTemp" binding:"required"`
	Days       string   `db:"days" json:"days" binding:"required"`
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

func DestroySchedule(d *Database, id int) error {
	connection := d.Connection()
	_, err := connection.Exec(`DELETE FROM schedule WHERE id=?`, id)

	if err != nil {
		log.Printf("db: error %s", err)
	}

	return err
}

func FindAllSchedules(d *Database) ([]Schedule, error) {
	connection := d.Connection()

	var schedules []Schedule
	err := connection.Select(&schedules, `SELECT * FROM schedule`)

	if err != nil {
		log.Printf("db: error %s", err)
	}

	return schedules, err
}

func (s *Schedule) Save(d *Database) error {
	connection := d.Connection()
	result, err := connection.NamedExec(`
		INSERT OR REPLACE INTO
		schedule (id, hour, target_temp, days)
		values (:id, :hour, :target_temp, :days)`, s)

	if err != nil {
		log.Printf("db: error %s", err)
	}

	id, _ := result.LastInsertId()
	s.Id.Scan(id)

	return err
}
