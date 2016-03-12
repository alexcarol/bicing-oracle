package prediction

// Prediction contains a prediction for a station at a certain time
type Prediction struct {
	Address string  `json:"address"`
	Slots   int     `json:"slots"`
	Bikes   int     `json:"bikes"`
	Lon     float32 `json:"lon"`
	Lat     float32 `json:"lat"`
}

// GetPredictions Returns an array of Prediction if everything goes alright
func GetPredictions() ([]Prediction, error) {
	return []Prediction{
		{"Gran via 123", 4, 3, 1.4, 1.3},
		{"Gran via 145", 3, 2, 1.1, 1.7},
	}, nil
}
