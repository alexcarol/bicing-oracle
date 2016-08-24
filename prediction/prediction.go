package prediction

import (
	"github.com/alexcarol/bicing-oracle/station-state/repository"
	"github.com/alexcarol/bicing-oracle/weather/datasource"
)

// Prediction contains a prediction for a station at a certain time
type Prediction struct {
	ID              uint    `json:"id"`
	Address         string  `json:"address"`
	BikeProbability float64 `json:"probability"`
	Lon             float64 `json:"lon"`
	Lat             float64 `json:"lat"`
}

// GetPredictions Returns an array of Prediction if everything goes alright
func GetPredictions(time int, lat float64, lon float64, stationProvider repository.StationProvider) ([]Prediction, error) {
	stations, err := stationProvider.GetNearbyStations(lat, lon, 3)
	if err != nil {
		return nil, err
	}

	var predictions = make([]Prediction, len(stations))

	weather, err := datasource.GetWeatherData()
	if err != nil {
		return nil, err
	}

	for i, station := range stations {
		probability, err := getBikeProbability(station.ID, time, weather.Type)
		if err != nil { // TODO consider ignoring failed cases but adding a metric
			return nil, err
		}

		predictions[i] = Prediction{
			station.ID,
			station.Street + ", " + station.StreetNumber,
			probability,
			station.Lon,
			station.Lat,
		}
	}

	return predictions, nil
}
