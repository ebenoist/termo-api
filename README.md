Termo-API
---

A [Go](http://golang.org/) based API for a DIY internet controlled thermostat.

### Hardware
<img src="http://erikbenoist.com/thermo-final.jpg" width=320 alt="diy-thermostat">
- Raspberry Pi (Model B)
- Sainsmart 4 Relay Module
- DS18B20 Temperature Sensor

### Dependencies
* [Go](http://golang.org)
* [gin](https://github.com/gin-gonic/gin) a Go web framework.
* [hwio](https://github.com/mrmorphic/hwio) a Go GPIO library.
* [Termo-UI](https://github.com/ebenoist/termo-ui) is an EmberJS powered UI for this API.

### Building (on OS X 10.9.5)
* Download and mount the arm5 [crosscompile toolchain](http://www.jaredwolff.com/toolchains/)
* `$ script/build_for_pi`

*Note:* The toolchain is only needed if there is a cgo dependency.

### Deploying
* Create `./api_key` see `./api_key.example`
* Edit `script/deploy` variables as needed
* `$ script/deploy`

*Note:* Currently termo-api is deployed as root so as to interface with GPIO pins, this is awful, and will be remedied soon.

### Starting
`$ API_KEY=foo ./termo-api`

### Issues
* Permissions issues are preventing running the server as a non superuser
* Hardware build needs to be documented.
* No unit tests

### JSON API

#### GET /v1/thermostat
```JSON
{
  {
    "targetTemp": 20,
    "heaterStatus" : true,
    "currentTemp": 18
  }
}
```

### POST /v1/thermostat
> set the target temperature

```JSON
{
  {
    "targetTemp": 30
  }
}
```
