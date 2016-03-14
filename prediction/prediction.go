package prediction

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
func GetPredictions(time int, lat float64, lon float64) ([]Prediction, error) {
	stations, err := getNearbyStations(lat, lon)
	if err != nil {
		return nil, err
	}

	var predictions = make([]Prediction, len(stations))

	for i, station := range stations {
		predictions[i] = Prediction{station.ID, station.Street + ", " + station.StreetNumber, station.FreeSlots, station.Bikes, station.Lon, station.Lat}
	}

	return predictions, nil
}

func getNearbyStations(lat float64, lon float64) ([]Station, error) {

	return []Station{
		{1, "Mecanic", "Carrer Joan Blanques", 123, "23", "", 12, 4, lon - 0.1, lat + 0.1},
		{2, "Mecanic", "Carrer Joan Blanques", 123, "23", "", 12, 4, lon - 0.1, lat + 0.1},
	}, nil
}

// Station contains info about a station
type Station struct {
	ID                int
	Type              string
	Street            string
	Height            int
	StreetNumber      string
	NearbyStationList string
	FreeSlots         int
	Bikes             int
	Lon               float64
	Lat               float64
}
