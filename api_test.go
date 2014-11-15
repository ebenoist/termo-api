package main_test

import (
	"bytes"
	"encoding/json"
	. "github.com/ebenoist/termo-api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("API", func() {
	database := &Database{}

	BeforeEach(func() {
		database.Create()
	})

	AfterEach(func() {
		database.Destroy()
	})

	Describe("/schedule", func() {
		It("POSTS a new schedule", func() {
			schedule := &Schedule{
				TargetTemp: 10,
				Hour:       8,
				Days:       WEEKENDS,
			}

			json, _ := json.Marshal(schedule)
			reqBody := bytes.NewBuffer(json)

			req, _ := http.NewRequest("POST", "/v1/schedules", reqBody)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			api := Api(&Thermostat{})
			api.ServeHTTP(w, req)
			w.Flush()

			Expect(w.Code).To(Equal(200))

			schedules, _ := FindAllSchedules(database)
			Expect(schedules).To(HaveLen(1))
			Expect(w.Body.String()).To(MatchJSON(`{
				"id": 1,
				"targetTemp": 10,
				"days": "weekends",
				"hour": 8
			}`))
		})

		It("GETS all the schedules", func() {
			schedule := &Schedule{
				TargetTemp: 10,
				Hour:       8,
				Days:       WEEKENDS,
			}

			scheduleTwo := &Schedule{
				TargetTemp: 11,
				Hour:       11,
				Days:       WEEKENDS,
			}

			schedule.Save(database)
			scheduleTwo.Save(database)

			req, _ := http.NewRequest("GET", "/v1/schedules", nil)
			w := httptest.NewRecorder()

			thermo := &Thermostat{
				TargetTemp:  DEFAULT_TARGET_TEMP,
				Heater:      &FakeHeater{},
				Thermometer: &FakeThermometer{},
			}

			api := Api(thermo)

			api.ServeHTTP(w, req)
			w.Flush()

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).To(MatchJSON(`[
				{"id":1,"hour":8,"targetTemp":10,"days":"weekends"},
				{"id":2,"hour":11,"targetTemp":11,"days":"weekends"}
			]`))
		})

		It("DELETES a schedule", func() {
			schedule := &Schedule{
				Hour:       4,
				TargetTemp: 20,
				Days:       WORKDAYS,
			}

			schedule.Save(database)

			schedules, _ := FindAllSchedules(database)
			Expect(schedules).To(HaveLen(1))

			req, _ := http.NewRequest("DELETE", "/v1/schedules/1", nil)
			w := httptest.NewRecorder()

			api := Api(&Thermostat{})

			api.ServeHTTP(w, req)
			w.Flush()

			Expect(w.Code).To(Equal(200))

			schedules, _ = FindAllSchedules(database)
			Expect(schedules).To(HaveLen(0))
		})
	})

	Describe("/thermostat", func() {
		It("GETS the current thermostat", func() {
			req, _ := http.NewRequest("GET", "/v1/thermostat", nil)
			w := httptest.NewRecorder()

			thermo := &Thermostat{
				TargetTemp:  DEFAULT_TARGET_TEMP,
				Heater:      &FakeHeater{},
				Thermometer: &FakeThermometer{},
			}

			api := Api(thermo)

			api.ServeHTTP(w, req)
			w.Flush()

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).To(MatchJSON(`{"currentTemp":17,"heaterOn":false,"targetTemp":20}`))
		})

		It("POSTS the current thermostat", func() {
			reqBody := bytes.NewBufferString(`{"targetTemp":25}`)

			req, _ := http.NewRequest("POST", "/v1/thermostat", reqBody)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			thermo := &Thermostat{
				TargetTemp:  DEFAULT_TARGET_TEMP,
				Heater:      &FakeHeater{},
				Thermometer: &FakeThermometer{},
			}

			api := Api(thermo)

			api.ServeHTTP(w, req)
			w.Flush()

			Expect(w.Code).To(Equal(200))
			Expect(thermo.TargetTemp).To(Equal(25))
		})
	})
})
