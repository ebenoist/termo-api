package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Thermometer struct {
	CurrentTemp int    `json:"currentTemp"` // cached value
	ReadTime    string `json:"readTime"`
}

func (t *Thermometer) setCurrentTemp() {
	t.ReadTime = time.Now().UTC().Format(time.RFC3339)

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

	log.Printf("Setting current temp to: %d", temp)
	t.CurrentTemp = temp
}
