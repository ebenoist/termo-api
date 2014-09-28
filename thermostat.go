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
	hwio.DigitalWrite(t.pin, hwio.LOW)
	t.On = true
}

func (t *Thermostat) turnOff() {
	hwio.DigitalWrite(t.pin, hwio.HIGH)
	t.On = false
}

func (t *Thermostat) adjust(currentTemp int) {
	if t.On {
		if currentTemp > t.TargetTemperature {
			log.Printf("Turning OFF the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemperature)
			t.turnOff()
		}
	} else {
		// fall TEMP_BUFFER bellow the target, turn heat on
		if currentTemp+TEMP_BUFFER < t.TargetTemperature {
			log.Printf("Turning ON the heat, currentTemp: %d, targetTemp: %d", currentTemp, t.TargetTemperature)
			t.turnOn()
		}
	}
}
