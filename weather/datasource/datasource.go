package datasource

import (
	"errors"
	"log"

	"github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/briandowns/openweathermap"
)

// Weather represents the current weather data
type Weather struct {
	Temperature     float64
	Type            int
	WindSpeed       float64
	WindDegree      float64
	CloudPercentage int
	Time            int
}

const temperatureUnit = "C"
const language = "ES"
const barcelonaWeatherID = 3128760

// GetForecast returns the weather forecast for a time
func GetForecast(time int) (int, float64, error) {
	// this uses only the daily prediction, better use 3-hour precision prediction
	forecastPredictor, err := openweathermap.NewForecast(temperatureUnit, language)
	if err != nil {
		return 0, 0, err
	}

	// TODO cache this to prevent too many api calls
	err = forecastPredictor.DailyByID(barcelonaWeatherID, 5)
	if err != nil {
		return 0, 0, err
	}

	var minI int
	var minDiff = -1

	for i, forecast := range forecastPredictor.List {
		if minDiff == -1 {
			minI = i
			minDiff = abs(time - forecast.Dt)
		} else if abs(time-forecast.Dt) < minDiff {
			minI = i
			minDiff = abs(time - forecast.Dt)
		}
	}

	forecast := forecastPredictor.List[minI]
	if len(forecast.Weather) < 0 {
		return 0, 0, errors.New("Weather length should not be 0")
	}
	temperature := forecast.Temp.Day

	return convertWeatherID(forecast.Weather[0].ID), temperature, nil
}

func abs(i int) int {
	if i > 0 {
		return int(i)
	}

	return int(-i)
}

// GetWeatherData returns the weather data
func GetWeatherData() (Weather, error) {
	var currentWeather Weather

	w, err := openweathermap.NewCurrent(temperatureUnit, language)
	if err != nil {
		return currentWeather, err
	}

	err = w.CurrentByID(barcelonaWeatherID) // Barcelona ID
	if err != nil {
		return currentWeather, err
	}
	currentWeather.Temperature = w.Main.Temp
	currentWeather.WindDegree = w.Wind.Deg
	currentWeather.WindSpeed = w.Wind.Speed
	currentWeather.CloudPercentage = w.Clouds.All
	currentWeather.Time = w.Dt

	if len(w.Weather) >= 1 {
		currentWeather.Type = convertWeatherID(w.Weather[0].ID)
	} else {
		log.Println("Weather only ")
		currentWeather.Type = -1
	}

	return currentWeather, nil
}

func convertWeatherID(id int) int {
	/**
	weather codes:
	2xx: Thunderstorm
	3xx: Drizzle
	5xx: Rain
	6xx: Snow
	7xx: Atmosphere
	800: Clear
	801-804: Clouds -> few, scattered, broken, overcast (negres)
	90x: Extreme
	910-999: Additional (calm, light breeze, gentle breeze, ..., storm, violent storm, hurricane
	*/
	switch id {
	case 781:
		return 13 // Tornado
	case 500:
	case 501:
		return 14 // slight rain
	case 502:
	case 503:
	case 504:
		return 15 // heavy rain
	case 511:
		return 16 // freezing rain
	case 520:
	case 521:
		return 17 // light to normal showers
	case 522:
		return 18 // heavy intensity shower
	case 531:
		return 19 // ragged shower rain

	case 800:
		return 4 // clear
	case 801:
		return 5 // few clouds
	case 802:
		return 6 // scattered clouds
	case 803:
		return 7 // broken clouds
	case 804:
		return 8 // overcast clouds

	}

	switch {

	case id >= 200 && id < 300:
		return 0 // Thunderstorm
	case id >= 300 && id < 400:
		return 1 // Drizzle
	case id >= 600 && id < 700:
		return 2 // Snow
	case id >= 700 && id <= 780:
		return 3 // Atmosphere (mist, fog, sand, ...)
	case id >= 900 && id < 910:
		// extreme weather
		return 9
	case id >= 951 && id <= 955:
		// gentle weather
		return 10
	case id >= 956 && id <= 959:
		// windy weather
		return 11
	case id >= 960 && id <= 962:
		// storm to hurricane
		return 12
	}

	return id
}
