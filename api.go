package main

import "fmt"
import "os"
import "time"
import "github.com/gin-gonic/gin"
import "strconv"

const (
	TEMP_SENSOR = "/sys/bus/w1/devices/28-0000061504ee/w1_slave"
	START_TEMP  = 69
	END_TEMP    = 71
)

type Temperature struct {
	ReadTime   string `json:"readTime"`
	Fahrenheit int    `json:"fahrenheit"`
	Celsius    int    `json:"celsius"`
}

func main() {
	fmt.Println("Booting termo!")
	r := gin.Default()

	r.GET("/v1/temperature", func(c *gin.Context) {
		t := new(Temperature)
		t.ReadTime = time.Now().UTC().Format(time.RFC3339)
		t.Celsius = readCurrentTemp()
		t.Fahrenheit = celsiusToF(t.Celsius)

		c.JSON(200, gin.H{"temperature": t})
	})

	r.Run(":8080")
}

func readCurrentTemp() int {
	file, err := os.Open(TEMP_SENSOR)
	if err != nil {
		fmt.Println("Error! %s", err)
	}

	data := make([]byte, 75)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Error! %s", err)
	}

	return parseTemp(data)
}

func parseTemp(data []byte) int {
	raw_temp := string(data[START_TEMP:END_TEMP])
	temp, _ := strconv.Atoi(raw_temp)
	return temp
}

func celsiusToF(celsius int) int {
	return (celsius*9)/5 + 32
}
