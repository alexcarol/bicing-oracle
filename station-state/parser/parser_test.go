package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXMLParsesReturnsErrWhenNotAssignedParsableXML(t *testing.T) {
	_, err := ParseXML([]byte("<asdf>"))

	require.NotNil(t, err)
}

func TestParseXMLReturnsAProperCollection(t *testing.T) {
	var correctData = []byte(`<?xml version="1.0" encoding="UTF-8"?>
    <bicing_stations>
     <updatetime><![CDATA[1415996588]]></updatetime>
      <station>
        <id>1</id>
          <type>BIKE</type>
            <lat>41.397952</lat>
              <long>2.180042</long>
                <street><![CDATA[Gran Via Corts Catalanes]]></street>
                  <height>21</height>
                    <streetNumber>760</streetNumber>
                      <nearbyStationList>24, 369, 387, 426</nearbyStationList>
                        <status>OPN</status>
                          <slots>0</slots>
                            <bikes>24</bikes>
                             </station>
                             </bicing_stations>`)

	collection, err := ParseXML(correctData)

	require.Nil(t, err)
	require.Equal(t, 1415996588, collection.Updatetime)
	require.Len(t, collection.StationStates, 1)

	stationState := collection.StationStates[0]

	require.Equal(t, 1, stationState.ID)
	require.Equal(t, 41.397952, stationState.Latitude)
	require.Equal(t, 2.180042, stationState.Longitude)
	require.Equal(t, "Gran Via Corts Catalanes", stationState.Street)
	require.Equal(t, 21, stationState.Height)
	require.Equal(t, "760", stationState.StreetNumber)
	require.Equal(t, "24, 369, 387, 426", stationState.NearbyStationList)
	require.Equal(t, "OPN", stationState.Status)
	require.Equal(t, 0, stationState.FreeSlots)
	require.Equal(t, 24, stationState.Bikes)
}
