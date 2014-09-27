package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ThermometerResponse struct {
	Thermometer *Thermometer `json:"thermometer"`
}

type ThermostatResponse struct {
	Thermostat *Thermostat `json:"thermostat"`
}

func apiRun(context *Context) {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/thermometer", func(c *gin.Context) {
		t := new(ThermometerResponse)
		t.Thermometer = context.Thermometer
		c.JSON(200, t)
	})

	v1.GET("/thermostat", func(c *gin.Context) {
		t := new(ThermostatResponse)
		t.Thermostat = context.Thermostat
		c.JSON(200, t)
	})

	v1.POST("/thermostat", func(c *gin.Context) {
		var json ThermostatResponse
		var targetTemp int

		c.Bind(&json)
		targetTemp = json.Thermostat.TargetTemperature
		log.Printf("Setting target temp to: %d", targetTemp)
		context.Thermostat.setTargetTemp(json.Thermostat.TargetTemperature)
		c.JSON(200, json)
	})

	r.Run(":8080")
}
