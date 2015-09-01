package parser

import (
    "time"
    "fmt"
    "encoding/xml"
    "github.com/alexcarol/bicing-api/station-state/collection"
)

func Parse(apiData []byte) collection.StationStateCollection {
	startTime := time.Now()
	requestEndTime := time.Now()

	var stationCollection collection.StationStateCollection

	err := xml.Unmarshal(apiData, &stationCollection)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n, structure :%s", err, apiData)
		return stationCollection
	}

	fmt.Printf("Data successfully received, request time: %v, unmarshalling time: %v\n", requestEndTime.Sub(startTime), time.Since(requestEndTime))
	return stationCollection
}
