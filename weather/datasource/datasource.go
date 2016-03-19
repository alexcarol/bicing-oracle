package datasource

import (
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

// GetWeatherData returns the weather data
func GetWeatherData() (Weather, error) {
	w, err := openweathermap.NewCurrent("C", "ES")

	var currentWeather Weather
	if err != nil {
		return currentWeather, err
	}

	err = w.CurrentByID(3128760) // Barcelona ID
	if err != nil {
		return currentWeather, err
	}
	temp := w.Main.Temp

	weatherID := convertWeatherID(w.Weather[0].ID)
	return Weather{temp, weatherID, w.Wind.Speed, w.Wind.Deg, w.Clouds.All, w.Dt}, nil
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
