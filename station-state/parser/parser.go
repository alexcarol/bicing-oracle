package parser

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/alexcarol/bicing-oracle/station-state/collection"
)

// ParseXML parses the xml bicing api data and returns it as an StationStateCollection
func ParseXML(apiData []byte) (collection.StationStateCollection, error) {
	startTime := time.Now()
	requestEndTime := time.Now()

	var stationCollection collection.StationStateCollection

	err := xml.Unmarshal(apiData, &stationCollection)
	if err != nil {
		return stationCollection, fmt.Errorf("Unmarshal error: %v\n, structure :%s", err, apiData)
	}

	fmt.Printf("Data successfully received, request time: %v, unmarshalling time: %v\n", requestEndTime.Sub(startTime), time.Since(requestEndTime))
	return stationCollection, nil
}
