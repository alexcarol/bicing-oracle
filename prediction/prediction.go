package prediction

import (
	"github.com/alexcarol/bicing-oracle/fitCalculator"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
	"github.com/alexcarol/bicing-oracle/weather/datasource"
)

// Prediction contains a prediction for a station at a certain time
type Prediction struct {
	ID              uint    `json:"id"`
	Address         string  `json:"address"`
	BikeProbability float64 `json:"bike-probability"`
	Lon             float64 `json:"lon"`
	Lat             float64 `json:"lat"`
	Failure         bool    `json:"failure"`
}

// GetStationPrediction returns the prediction for a station
func GetStationPrediction(time int, stationID uint, stationProvider repository.StationProvider) (Prediction, error) {
	station, err := stationProvider.GetStationByID(stationID)
	if err != nil {
		return Prediction{}, err
	}

	weather, temperature, err := datasource.GetForecast(time)
	if err != nil {
		return Prediction{}, err
	}

	bikes, err := getProbability(station.ID, time, weather, temperature)

	return Prediction{
		station.ID,
		station.Street + ", " + station.StreetNumber,
		bikes,
		station.Lon,
		station.Lat,
		err != nil,
	}, nil
}

// GetPredictions Returns an array of Prediction if everything goes alright
func GetPredictions(time int, lat float64, lon float64, stationProvider repository.StationProvider) ([]Prediction, error) {
	stations, err := stationProvider.GetNearbyStations(lat, lon, 3)
	if err != nil {
		return nil, err
	}

	var predictions = make([]Prediction, len(stations))

	weather, temperature, err := datasource.GetForecast(time)
	if err != nil {
		return nil, err
	}

	for i, station := range stations {
		probability, err := getProbability(station.ID, time, weather, temperature)
		if err != nil { // TODO consider adding a metric
			fitCalculator.ScheduleCalculate(station.ID)
		}

		predictions[i] = Prediction{
			station.ID,
			station.Street + ", " + station.StreetNumber,
			probability,
			station.Lon,
			station.Lat,
			err != nil,
		}
	}

	return predictions, nil
}
