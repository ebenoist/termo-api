package main

import (
	"github.com/mrmorphic/hwio"
	"log"
)

// Heater is a light wrapper around the hwio package
type Heater struct {
	On  bool
	pin hwio.Pin
}

func (h *Heater) TurnOn() {
	h.setPin(true)
	h.On = true
}

func (h *Heater) TurnOff() {
	h.setPin(false)
	h.On = false
}

func (h *Heater) ShutDown() {
	h.TurnOff()
	hwio.CloseAll()
}

func (h *Heater) setPin(on bool) {
	if h.pin == 0 { // pin will be 0 if not set
		h.openPin()
	}

	if on {
		hwio.DigitalWrite(h.pin, hwio.LOW)
	} else {
		hwio.DigitalWrite(h.pin, hwio.HIGH)
	}
}

func (h *Heater) openPin() {
	hwio.SetDriver(new(hwio.RaspberryPiDTDriver))
	pin, err := hwio.GetPinWithMode("gpio17", hwio.OUTPUT)

	if err != nil {
		log.Printf("Error opening pin! %s\n", err)
	}

	h.pin = pin
}
