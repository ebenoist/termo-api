package main

import (
	"github.com/mrmorphic/hwio"
	"log"
)

type Heater interface {
	TurnOn()
	TurnOff()
	ShutDown()
	On() bool
}

// Heater is a light wrapper around the hwio package
type RealHeater struct {
	on  bool
	pin hwio.Pin
}

func (h *RealHeater) On() bool {
	return h.on
}

func (h *RealHeater) TurnOn() {
	h.setPin(true)
	h.on = true
}

func (h *RealHeater) TurnOff() {
	h.setPin(false)
	h.on = false
}

func (h *RealHeater) ShutDown() {
	h.TurnOff()
	hwio.CloseAll()
}

func (h *RealHeater) setPin(on bool) {
	if h.pin == 0 { // pin will be 0 if not set
		h.openPin()
	}

	if on {
		hwio.DigitalWrite(h.pin, hwio.LOW)
	} else {
		hwio.DigitalWrite(h.pin, hwio.HIGH)
	}
}

func (h *RealHeater) openPin() {
	hwio.SetDriver(new(hwio.RaspberryPiDTDriver))
	pin, err := hwio.GetPinWithMode("gpio17", hwio.OUTPUT)

	if err != nil {
		log.Printf("Error opening pin! %s\n", err)
	}

	h.pin = pin
}
