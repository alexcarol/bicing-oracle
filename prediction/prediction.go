package prediction

import (
	"math/rand"

	"github.com/alexcarol/bicing-oracle/station-state/repository"
)

// Prediction contains a prediction for a station at a certain time
type Prediction struct {
	ID      int     `json:"id"`
	Address string  `json:"address"`
	Slots   int     `json:"slots"`
	Bikes   int     `json:"bikes"`
	Lon     float64 `json:"lon"`
	Lat     float64 `json:"lat"`
}

// GetPredictions Returns an array of Prediction if everything goes alright
func GetPredictions(time int, lat float64, lon float64, stationProvider repository.StationProvider) ([]Prediction, error) {
	stations, err := stationProvider.GetNearbyStations(lat, lon, 3)
	if err != nil {
		return nil, err
	}

	var predictions = make([]Prediction, len(stations))

	for i, station := range stations {
		predictions[i] = Prediction{station.ID, station.Street + ", " + station.StreetNumber, rand.Int() % 20, rand.Int() % 20, station.Lon, station.Lat}
	}

	return predictions, nil
}
