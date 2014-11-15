package main_test

import (
	"database/sql"
	. "github.com/ebenoist/termo-api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Schedule", func() {
	database := &Database{}

	BeforeEach(func() {
		database.Create()
	})

	AfterEach(func() {
		database.Destroy()
	})

	It("Has an Hour and TargetTemp", func() {
		schedule := &Schedule{TargetTemp: 20, Hour: 2}

		Expect(schedule.TargetTemp).To(Equal(20))
		Expect(schedule.Hour).To(Equal(2))
	})

	It("Has an list of days", func() {
		schedule := &Schedule{Days: WEEKENDS}

		Expect(schedule.Days).To(Equal(WEEKENDS))
	})

	It("Can be saved", func() {
		schedule := &Schedule{
			TargetTemp: 20,
			Hour:       2,
			Days:       WEEKENDS,
		}

		schedule.Save(database)

		var id int64
		id = 1

		Expect(schedule.Id.Value()).To(Equal(id))

		var result Schedule
		database.Connection().Get(&result, "SELECT * FROM schedule WHERE id=$1", id)

		Expect(result.Hour).To(Equal(2))
		Expect(result.Days).To(Equal("weekends"))
		Expect(result.TargetTemp).To(Equal(20))
	})

	It("Can be updated", func() {
		schedule := &Schedule{
			TargetTemp: 20,
			Hour:       2,
			Days:       WEEKENDS,
		}

		schedule.Save(database)
		schedule.TargetTemp = 10
		schedule.Save(database)

		var result Schedule
		database.Connection().Get(&result, "SELECT * FROM schedule WHERE id=1")

		Expect(result.Hour).To(Equal(2))
		Expect(result.Days).To(Equal("weekends"))
		Expect(result.TargetTemp).To(Equal(10))
	})

	It("Can be removed", func() {
		schedule := &Schedule{
			TargetTemp: 20,
			Hour:       2,
			Days:       WEEKENDS,
		}

		schedule.Save(database)
		DestroySchedule(database, 1)

		var result Schedule
		err := database.Connection().Get(&result, "SELECT * FROM schedule WHERE id=1")

		Expect(result.TargetTemp).To(Equal(0))
		Expect(err).To(Equal(sql.ErrNoRows))
	})

	It("Can be found by a given time", func() {
		schedule := &Schedule{
			TargetTemp: 10,
			Hour:       10,
			Days:       WEEKENDS,
		}

		distractingSchedule := &Schedule{
			TargetTemp: 11,
			Hour:       11,
			Days:       WEEKENDS,
		}

		schedule.Save(database)
		distractingSchedule.Save(database)

		searchTime, _ := time.Parse(time.RFC3339, "2014-11-15T10:23:00-06:00")
		foundTime := FindScheduleByTime(database, searchTime)

		Expect(foundTime.TargetTemp).To(Equal(10))
	})

	It("Finds the currently running schedule", func() {
		schedule := &Schedule{
			TargetTemp: 10,
			Hour:       8,
			Days:       WEEKENDS,
		}

		distractingOne := &Schedule{
			TargetTemp: 11,
			Hour:       11,
			Days:       WEEKENDS,
		}

		distractingTwo := &Schedule{
			TargetTemp: 11,
			Hour:       7,
			Days:       WEEKENDS,
		}

		distractingThree := &Schedule{
			TargetTemp: 11,
			Hour:       10,
			Days:       WORKDAYS,
		}

		schedule.Save(database)
		distractingOne.Save(database)
		distractingTwo.Save(database)
		distractingThree.Save(database)

		searchTime, _ := time.Parse(time.RFC3339, "2014-11-15T10:23:00-06:00")
		foundTime := FindScheduleByTime(database, searchTime)

		Expect(foundTime.TargetTemp).To(Equal(10))
	})
})
