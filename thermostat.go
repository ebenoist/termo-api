package main

import (
	"github.com/mrmorphic/hwio"
	"log"
)

type Thermostat struct {
	TargetTemperature int  `json:"targetTemperature"`
	On                bool `json:on"`
	pin               hwio.Pin
}

func (t *Thermostat) setTargetTemp(targetTemp int) {
	t.TargetTemperature = targetTemp
}

func (t *Thermostat) turnOn() {
	t.On = true
}

func (t *Thermostat) turnOff() {
	t.On = false
}

func (t *Thermostat) adjust(currentTemp int) {
	if currentTemp+TEMP_BUFFER > t.TargetTemperature {
		if t.On {
			hwio.DigitalWrite(t.pin, hwio.HIGH)
			log.Printf("Turning OFF the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemperature)
			t.turnOff()
		}
	}

	if currentTemp-TEMP_BUFFER < t.TargetTemperature {
		if !t.On {
			hwio.DigitalWrite(t.pin, hwio.LOW)
			log.Printf("Turning ON the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemperature)
			t.turnOn()
		}
	}

	log.Printf("Leaving the heat %t, currentTemp: %d, targetTemp: %d", t.On, currentTemp, t.TargetTemperature)
}
