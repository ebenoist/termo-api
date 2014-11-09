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
)

func buildThermostat() *Thermostat {
	if os.Getenv("TERMO_MOCK") != "" {
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

func main() {
	log.Println("Booting termo!")
	if os.Getenv("API_KEY") == "" {
		panic("Could not find API_KEY")
	}

	log.Printf("Using key %s", os.Getenv("API_KEY"))
	thermostat := buildThermostat()

	// CLEAN UP
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Caught %v", sig)
			log.Println("Exiting Termo!!")
			thermostat.Stop()
			os.Exit(1)
		}
	}()

	go thermostat.Run()
	apiRun(thermostat)
}
