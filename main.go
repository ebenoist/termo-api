package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	MONITOR_INTERVAL    = 30 * time.Second
	TEMP_BUFFER         = 2
	DEFAULT_TARGET_TEMP = 20
	DB_NAME             = "db/termo-db"
)

func buildThermostat() *Thermostat {
	if ENV == "DEVELOPMENT" {
		log.Println("Booting in mock mode")

		return &Thermostat{
			TargetTemp:  DEFAULT_TARGET_TEMP,
			Heater:      &FakeHeater{},
			Thermometer: &FakeThermometer{},
		}
	}

	return &Thermostat{
		TargetTemp:  DEFAULT_TARGET_TEMP,
		Heater:      &RealHeater{},
		Thermometer: &RealThermometer{},
	}
}

var (
	ENV     = "DEVELOPMENT"
	API_KEY = ""
)

func setupEnv() {
	if os.Getenv("ENV") != "" {
		ENV = os.Getenv("ENV")
	}

	if os.Getenv("API_KEY") == "" && ENV == "production" {
		panic("Could not find API_KEY")
	}

	log.Printf("Using key %s", os.Getenv("API_KEY"))
	log.Printf("Booting termo in %s", ENV)
}

func cleanUp(t *Thermostat) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Caught %v", sig)
			log.Println("Exiting Termo!!")
			t.Stop()
			os.Exit(1)
		}
	}()
}

func main() {
	setupEnv()
	thermostat := buildThermostat()
	cleanUp(thermostat)

	go thermostat.Run()
	go SchedulerRun(thermostat)
	api := Api(thermostat)
	api.Run(":8080")
}
