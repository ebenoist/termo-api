package main

import (
	"github.com/mrmorphic/hwio"
	"log"
)

type FakeHeater struct {
	on  bool
	pin hwio.Pin
}

func (h *FakeHeater) On() bool {
	return h.on
}

func (h *FakeHeater) TurnOn() {
	log.Printf("TurnOn was called.")
	h.on = true
}

func (h *FakeHeater) TurnOff() {
	log.Printf("TurnOff was called.")
	h.on = false
}

func (h *FakeHeater) ShutDown() {
	h.TurnOff()
	log.Printf("ShutDown was called.")
}

type FakeThermometer struct{}

func (t *FakeThermometer) ReadTemp() int {
	return 17
}
