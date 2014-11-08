package main

import (
	"log"
	"time"
)

type Thermostat struct {
	TargetTemp  int
	Heater      Heater
	Thermometer Thermometer
}

func (t *Thermostat) Run() {
	for {
		t.Adjust()
		time.Sleep(MONITOR_INTERVAL)
	}
}

func (t *Thermostat) Stop() {
	t.Heater.ShutDown()
}

func (t *Thermostat) CurrentTemp() int {
	return t.Thermometer.ReadTemp()
}

func (t *Thermostat) HeaterOn() bool {
	return t.Heater.On()
}

func (t *Thermostat) Adjust() {
	var currentTemp = t.CurrentTemp()
	if t.HeaterOn() {
		if currentTemp > t.TargetTemp {
			log.Printf("Turning OFF the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemp)
			t.Heater.TurnOff()
		}
	} else {
		// fall TEMP_BUFFER bellow the target, turn heat on
		if t.CurrentTemp()+TEMP_BUFFER < t.TargetTemp {
			log.Printf("Turning ON the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemp)
			t.Heater.TurnOn()
		}
	}
}
