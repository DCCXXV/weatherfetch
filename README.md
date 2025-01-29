## Screenshot

![screenshot](screenshot/img.png)

## Usage
```bash
weatherfetch

      .--.
   .-(    ).
  (___.__)__)

Cloudy
Temperature: 5.10 ºC
Humidity: 97.00 %
Wind speed: 3.70 m/s
```

## Dependencies
* A working Go environment
* A Tomorrow.io API key (see Setup below)

## Setup
Get a Tomorrow.io API key[here](https://www.tomorrow.io/weather-api/)
Put the api along with your coordinates in `config.go`:

```go
var Api = "TOMORROWIO_API_KEY"

var Latitude = "0"
var Longitude = "0"
```

## Installation

Clone the repository:
```bash
git clone https://github.com/DCCXXV/weatherfetch
```

Then build it with:
```bash
cd weatherfetch
go build -o weatherfetch config.go weatherfetch.go
```
And for system-wide access:
```bash
sudo mv weatherfetch /usr/local/bin/
```

## Thanks
* [wego](https://github.com/schachmat/wego) for providing all the ASCII icons
