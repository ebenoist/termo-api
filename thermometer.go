package main

import (
	"log"
	"os"
	"strconv"
)

const (
	TEMP_SENSOR = "/sys/bus/w1/devices/28-0000061504ee/w1_slave"
	START_TEMP  = 69 // numer of characters to read before temp data
	END_TEMP    = 71 // read until the 71st character
)

type Thermometer struct {
	lastReadTemp int
}

// ReadTemp will return the last read temperature
// and asynchronously schedule a new read
// in the case of no known temperature, ReadTemp will
// synchronously read from the sensor which can cost ~ 500ms
func (t *Thermometer) ReadTemp() int {
	if t.lastReadTemp != 0 {
		go t.readFromSensor()
		return t.lastReadTemp
	} else {
		t.readFromSensor()
		return t.lastReadTemp
	}
}

func (t *Thermometer) readFromSensor() {
	file, err := os.Open(TEMP_SENSOR)
	if err != nil {
		log.Println("Error! %s", err)
	}

	data := make([]byte, 75)
	_, err = file.Read(data)
	if err != nil {
		log.Println("Error! %s", err)
	}

	raw_temp := string(data[START_TEMP:END_TEMP])
	temp, _ := strconv.Atoi(raw_temp)

	t.lastReadTemp = temp
}
