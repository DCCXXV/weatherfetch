package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type WeatherResponse struct {
	Timelines struct {
		Hourly []struct {
			Time   string `json:"time"`
			Values struct {
				Temperature float64 `json:"temperature"`
				Humidity    float64 `json:"humidity"`
				WindSpeed   float64 `json:"windSpeed"`
				WeatherCode int     `json:"weatherCode"`
			} `json:"values"`
		} `json:"hourly"`
	} `json:"timelines"`
}

type CurrentWeather struct {
	Time   string
	Values struct {
		Temperature float64
		Humidity    float64
		WindSpeed   float64
		WeatherCode int
		Description string
	}
}

func fetchData() (CurrentWeather, error) {
	url := fmt.Sprintf("https://api.tomorrow.io/v4/weather/forecast?location=%s,%s&fields=temperature,humidity&timesteps=1h&units=metric&apikey=%s",
		Latitude, Longitude, Api)

	now := time.Now().UTC()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Error creating the request: %v", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error doing the request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading the response: %v", err)
	}

	// parsing to json
	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)

	if err != nil {
		log.Fatal("Error parsing JSON: %n", err)
	}

	minDiff := time.Duration(math.MaxInt64)
	var currentWeather CurrentWeather

	for _, hourly := range weatherData.Timelines.Hourly {
		entryTime, err := time.Parse(time.RFC3339, hourly.Time)

		if err != nil {
			log.Fatal("Error parsing time")
		}

		diff := entryTime.Sub(now).Abs()
		if diff < minDiff {
			minDiff = diff
			currentWeather.Time = hourly.Time
			currentWeather.Values.Humidity = hourly.Values.Humidity
			currentWeather.Values.Temperature = hourly.Values.Temperature
			currentWeather.Values.WindSpeed = hourly.Values.WindSpeed
			currentWeather.Values.WeatherCode = hourly.Values.WeatherCode
		}
	}

	currentWeather.Values.Description = getDescription(strconv.Itoa(currentWeather.Values.WeatherCode))

	return currentWeather, err
}

func getDescription(code string) string {
	jsonData := `{
	"0": "Unknown",
	"1000": "Clear, Sunny",
	"1100": "Mostly Clear",
	"1101": "Partly Cloudy",
	"1102": "Mostly Cloudy",
	"1001": "Cloudy",
	"2000": "Fog",
	"2100": "Light Fog",
	"4000": "Drizzle",
	"4001": "Rain",
	"4200": "Light Rain",
	"4201": "Heavy Rain",
	"5000": "Snow",
	"5001": "Flurries",
	"5100": "Light Snow",
	"5101": "Heavy Snow",
	"6000": "Freezing Drizzle",
	"6001": "Freezing Rain",
	"6200": "Light Freezing Rain",
	"6201": "Heavy Freezing Rain",
	"7000": "Ice Pellets",
	"7101": "Heavy Ice Pellets",
	"7102": "Light Ice Pellets",
	"8000": "Thunderstorm"
	}`

	var weatherCodes map[string]string
	err := json.Unmarshal([]byte(jsonData), &weatherCodes)
	if err != nil {
		log.Fatal("Error parsing JSON ", err)
	}

	description, _ := weatherCodes[code]

	return description
}

func printAscii(weatherCode string) {
	switch weatherCode {
	case "0":
		unk := `
    .-.
     __)
    (
     '-'
      •
      `
		fmt.Print(unk)
	case "1000", "1100":
		clear := `
    \   /
     .-.
  ― (   ) ―
     '-’
    /   \
    `
		fmt.Println(clear)
	case "1101", "1102", "1001":
		cloudy := `
      .--.
   .-(    ).
  (___.__)__)
    `
		fmt.Println(cloudy)
	case "2000", "2100":
		fog := `
	 _ - _ - _ -
	  _ - _ - _
	 _ - _ - _ -
	`
		fmt.Print(fog)
	case "4000", "4001", "4200", "4201", "6000", "6001", "6200", "6201":
		rain := `
     .-.
    (   ).
   (___(__)
    ' ' ' '
   ' ' ' '
   `
		fmt.Print(rain)
	case "5000", "5001", "5100", "5101":
		snow := `
      .-.
     (   ).
    (___(__)
     *  *  *
    *  *  *
    `
		fmt.Print(snow)
	case "7000", "7101", "7102":
		icePellet := `
      .-.
     (   ).
    (___(__)
     |  |  |
    |  |  |
    `
		fmt.Print(icePellet)
	}
}

func printInfo(weather CurrentWeather) {
	printAscii(strconv.Itoa(weather.Values.WeatherCode))
	fmt.Println(weather.Values.Description)
	fmt.Printf("Temperature: %.2f ºC\n", weather.Values.Temperature)
	fmt.Printf("Humidity: %.2f %%\n", weather.Values.Humidity)
	fmt.Printf("Wind speed: %.2f m/s\n", weather.Values.WindSpeed)
}

func main() {
	weather, _ := fetchData()
	printInfo(weather)
}
