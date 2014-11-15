package main_test

import (
	"bytes"
	. "github.com/ebenoist/termo-api"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("API", func() {
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

			gin.SetMode(gin.ReleaseMode)
			api := Api(thermo)

			api.ServeHTTP(w, req)
			w.Flush()
			Expect(w.Code).To(Equal(200))

			Expect(thermo.TargetTemp).To(Equal(25))
		})
	})
})
