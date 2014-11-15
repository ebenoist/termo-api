package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
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

func Api(thermostat *Thermostat) *gin.Engine {
	r := gin.Default()
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

	return r
}
