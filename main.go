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
	TEMP_BUFFER      = 1
)

type Context struct {
	Thermostat  *Thermostat
	Thermometer *Thermometer
}

func buildContext() *Context {
	return &Context{
		Thermostat: &Thermostat{
			On:                false,
			TargetTemperature: 20,
			pin:               openPin(),
		},
		Thermometer: &Thermometer{},
	}
}

func openPin() hwio.Pin {
	hwio.SetDriver(new(hwio.RaspberryPiDTDriver))
	time.Sleep(2 * time.Second)
	pin, err := hwio.GetPinWithMode("gpio17", hwio.OUTPUT)

	if err != nil {
		log.Printf("Error opening pin! %s\n", err)
	}

	time.Sleep(2 * time.Second)
	// PIN OFF ON BOOT
	hwio.DigitalWrite(pin, hwio.HIGH)

	return pin
}

func main() {
	log.Println("Booting termo!")

	if os.Getenv("API_KEY") != "" {
		log.Printf("Using key %s", os.Getenv("API_KEY"))

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
}

func monitorRun(context *Context) {
	for {
		context.Thermometer.setCurrentTemp()
		context.Thermostat.adjust(context.Thermometer.CurrentTemp)
		time.Sleep(MONITOR_INTERVAL)
	}
}
