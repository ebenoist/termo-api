package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, API-KEY")
		if c.Request.Method == "OPTIONS" {
			c.Abort(200)
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("API-KEY") != os.Getenv("API_KEY") {
			c.Abort(401)
			return
		}

		c.Next()
	}
}

func presentThermostat(thermostat *Thermostat) gin.H {
	return gin.H{
		"currentTemp": thermostat.CurrentTemp(),
		"targetTemp":  thermostat.TargetTemp,
		"heaterOn":    thermostat.HeaterOn(),
	}
}

type ThermostatRequest struct {
	TargetTemp int `json:"targetTemp" binding:"required"`
}

type ScheduleResponse struct {
	Id         int    `json:"id"`
	Hour       int    `json:"hour" binding:"required"`
	TargetTemp int    `json:"target_temp" binding:"required"`
	Days       string `json:"days" binding:"required"`
}

func Api(thermostat *Thermostat) *gin.Engine {
	database := &Database{}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(CORSMiddleware())

	if os.Getenv("TERMO_MOCK") != "true" {
		r.Use(AuthMiddleware())
	}

	v1 := r.Group("/v1")
	v1.GET("/thermostat", func(c *gin.Context) {
		c.JSON(200, presentThermostat(thermostat))
	})

	v1.POST("/thermostat", func(c *gin.Context) {
		var json ThermostatRequest
		c.Bind(&json)

		var targetTemp int
		targetTemp = json.TargetTemp
		log.Printf("Setting target temp to: %d", targetTemp)

		thermostat.TargetTemp = json.TargetTemp

		c.JSON(200, json)
	})

	v1.GET("/schedules", func(c *gin.Context) {
		schedules, _ := FindAllSchedules(database)

		c.JSON(200, schedules)
	})

	v1.POST("/schedules", func(c *gin.Context) {
		var schedule Schedule
		bound := c.Bind(&schedule)

		if !bound {
			c.Fail(400, errors.New("Invalid schedule"))
		} else {
			schedule.Save(database)
			c.JSON(200, schedule)
		}
	})

	v1.DELETE("/schedules/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))

		if err != nil {
			c.Fail(400, errors.New("Malformed id"))
		} else {
			DestroySchedule(database, id)
			c.String(200, "")
		}
	})

	return r
}
