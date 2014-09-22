package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	TEMP_SENSOR      = "/sys/bus/w1/devices/28-0000061504ee/w1_slave"
	START_TEMP       = 69
	END_TEMP         = 71
	MONITOR_INTERVAL = 30 * time.Second
	TEMP_KEY         = "termo-api:temperature"
	TEMP_TIME_KEY    = "termo-api:read-time"
	TARGET_TEMP_KEY  = "termo-api:target-temperature"
	THERMO_ON_KEY    = "termo-api:thermo-on"
	TEMP_BUFFER      = 2 // Buffer temp by 2 degrees celsius
)

type Context struct {
	Redis redis.Conn
}

func redisConnection() redis.Conn {
	redisConnection, err := redis.Dial("tcp", ":6379") // TODO: error?
	if err != nil {
		log.Printf("Error connecting to redis: %q")
	}

	// teardown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			redisConnection.Close()
			log.Println("Exiting... %q\n", sig)
			os.Exit(1)
		}
	}()

	return redisConnection
}

func main() {
	log.Println("Booting termo!")
	context := buildContext()
	go monitorRun(context)
	apiRun(context)
}

func buildContext() *Context {
	context := new(Context)
	context.Redis = redisConnection()
	return context
}

// -------------------------------------------------
// API

type TemperatureResponse struct {
	Temperature Temperature `json:"temperature"`
}

type Temperature struct {
	ReadTime   string `json:"readTime"`
	Fahrenheit int    `json:"fahrenheit"`
	Celsius    int    `json:"celsius"`
}

type ThermostatResponse struct {
	Thermostat Thermostat `json:"thermostat"`
}
type Thermostat struct {
	TargetTemperature int  `json:"targetTemperature"`
	On                bool `json:on"`
}

func apiRun(context *Context) {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/temperature", func(c *gin.Context) {
		t := new(TemperatureResponse)
		t.Temperature.ReadTime, _ = redis.String(context.Redis.Do("GET", TEMP_TIME_KEY))
		t.Temperature.Celsius, _ = redis.Int(context.Redis.Do("GET", TEMP_KEY))
		t.Temperature.Fahrenheit = celsiusToF(t.Temperature.Celsius)

		c.JSON(200, t)
	})

	v1.GET("/thermostat", func(c *gin.Context) {
		t := new(ThermostatResponse)
		t.Thermostat.TargetTemperature, _ = redis.Int(context.Redis.Do("GET", TARGET_TEMP_KEY))
		t.Thermostat.On, _ = redis.Bool(context.Redis.Do("GET", THERMO_ON_KEY))

		c.JSON(200, t)
	})

	v1.POST("/thermostat", func(c *gin.Context) {
		var json ThermostatResponse
		c.Bind(&json)

		context.Redis.Do("SET", TARGET_TEMP_KEY, json.Thermostat.TargetTemperature)

		c.JSON(200, json)
	})

	r.Run(":8080")
}

// -------------------------------------------------
// MONITOR
func monitorRun(context *Context) {
	var pin gpio.Pin

	pin = openPin()

	for {
		var currentTemp int
		var targetTemp int
		var thermoOn bool

		currentTemp = readCurrentTemp()

		context.Redis.Do("SET", TEMP_KEY, currentTemp)
		context.Redis.Do("SET", TEMP_TIME_KEY, time.Now().UTC().Format(time.RFC3339))

		targetTemp, _ = redis.Int(context.Redis.Do("GET", TARGET_TEMP_KEY))
		thermoOn, _ = redis.Bool(context.Redis.Do("GET", THERMO_ON_KEY))

		var turnedOn bool
		turnedOn = adjustThermostat(currentTemp, targetTemp, thermoOn, pin)
		context.Redis.Do("SET", THERMO_ON_KEY, turnedOn)

		time.Sleep(MONITOR_INTERVAL)
	}

	pin.Close()
}

func openPin() gpio.Pin {
	pin, err := gpio.OpenPin(rpi.GPIO17, gpio.ModeOutput)
	if err != nil {
		log.Printf("Error opening pin! %s\n", err)
		return nil
	}

	// Turn off HEAT at exit
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func() {
	// for _ = range c {
	// log.Printf("\nClearing and unexporting the pin.\n")
	// pin.Set()
	// pin.Close()
	// }
	// }()

	return pin
}

func adjustThermostat(currentTemp int, targetTemp int, thermoOn bool, pin gpio.Pin) bool {
	if currentTemp+TEMP_BUFFER > targetTemp {
		if thermoOn {
			pin.Set()
			log.Printf("Turning OFF the heat, currentTemp: %d, targetTemp: %d", currentTemp, targetTemp)
			return false
		}
	}

	if currentTemp-TEMP_BUFFER < targetTemp {
		if !thermoOn {
			pin.Clear()
			log.Printf("Turning ON the heat, currentTemp: %d, targetTemp: %d", currentTemp, targetTemp)
			return true
		}
	}

	return thermoOn // No change
}

func readCurrentTemp() int {
	file, err := os.Open(TEMP_SENSOR)
	if err != nil {
		log.Println("Error! %s", err)
	}

	data := make([]byte, 75)
	_, err = file.Read(data)
	if err != nil {
		log.Println("Error! %s", err)
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
