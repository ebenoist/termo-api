package main

import (
	"github.com/mrmorphic/hwio"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	TEMP_SENSOR      = "/sys/bus/w1/devices/28-0000061504ee/w1_slave"
	START_TEMP       = 69
	END_TEMP         = 71
	MONITOR_INTERVAL = 30 * time.Second
	TEMP_BUFFER      = 2 // Buffer temp by 2 degrees celsius
)

type Context struct {
	Thermostat  *Thermostat
	Thermometer *Thermometer
}

func buildContext() *Context {
	context := new(Context)

	context.Thermostat = new(Thermostat)
	context.Thermostat.On = false
	context.Thermostat.TargetTemperature = 20
	context.Thermostat.pin = openPin()

	context.Thermometer = new(Thermometer)

	return context
}

func openPin() hwio.Pin {
	hwio.SetDriver(new(hwio.RaspberryPiDTDriver))
	pin, err := hwio.GetPinWithMode("gpio17", hwio.OUTPUT)

	if err != nil {
		log.Printf("Error opening pin! %s\n", err)
	}

	// PIN OFF ON BOOT
	hwio.DigitalWrite(pin, hwio.HIGH)

	return pin
}

func main() {
	log.Println("Booting termo!")

	context := buildContext()

	// CLEAN UP
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Caught %v", sig)
			log.Println("Exiting Termo!!")

			context.Thermostat.turnOff()

			hwio.CloseAll()
			os.Exit(1)
		}
	}()

	go monitorRun(context)
	apiRun(context)
}

func monitorRun(context *Context) {
	for {
		context.Thermometer.setCurrentTemp()
		context.Thermostat.adjust(context.Thermometer.CurrentTemp)
		time.Sleep(MONITOR_INTERVAL)
	}
}
