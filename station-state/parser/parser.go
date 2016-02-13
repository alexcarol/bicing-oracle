package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/alexcarol/bicing-oracle/station-state/collection"
)

// ParseXML parses the xml bicing api data and returns it as an StationStateCollection
func ParseXML(apiData []byte) (collection.StationStateCollection, error) {
	var stationCollection collection.StationStateCollection

	err := xml.Unmarshal(apiData, &stationCollection)
	if err != nil {
		return stationCollection, fmt.Errorf("Unmarshal error: %v\n, structure :%s", err, apiData)
	}

	return stationCollection, nil
}
