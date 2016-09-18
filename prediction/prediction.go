package prediction

import (
	"log"

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

	currentBikes, err := getCurrentBikes(station.ID, stationProvider)
	if err != nil {
		return Prediction{}, err
	}

	probability, err := getProbability(station.ID, time, weather, temperature, currentBikes)

	return Prediction{
		station.ID,
		station.Street + ", " + station.StreetNumber,
		probability,
		station.Lon,
		station.Lat,
		err != nil,
	}, err
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
		currentBikes, err := getCurrentBikes(station.ID, stationProvider)
		if err != nil {
			log.Printf("Error obtaining bikes for station %d\n", station.ID)
			continue
		}

		probability, err := getProbability(station.ID, time, weather, temperature, currentBikes)
		if err != nil { // TODO consider adding a metric
			fitCalculator.ScheduleCalculate(station.ID)
			log.Println("Error getting probability for station", station.ID, err.Error())
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

func getCurrentBikes(stationID uint, stationProvider repository.StationProvider) (int, error) {
	stationState, err := stationProvider.GetCurrentStationStateStationByID(stationID)

	return stationState.Bikes, err
}
