Termo-API
---

A [Go](http://golang.org/) based API for a DIY internet controlled thermostat.

### Hardware
<img src="http://erikbenoist.com/thermo-final.jpg" width=320 alt="diy-thermostat">
- Raspberry Pi (Model B)
- Sainsmart 4 Relay Module
- DS18B20 Temperature Sensor

### JSON API
#### GET /v1/thermometer
```JSON
{
  "thermometer": {
    "readTime": "2014-09-20T20:45:40Z",
    "fahrenheit": 72,
    "celsius": 22
  }
}
```

#### GET /v1/thermostat
```JSON
{
  "thermostat": {
    "targetTemperature": 68,
    "on" : true
  }
}
```

### POST /v1/thermostat
> set the target temperature

```JSON
{
  "thermostat": {
    "targetTemperature": 30
  }
}
```
