package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func apiFakeDataProvider() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
}

func doCurl() []byte {
	response, err := http.Get("http://wservice.viabicing.cat/v1/getstations.php?v=1")
	if err != nil {
		fmt.Printf("Error with the request %s", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error with the request %s", err)
		os.Exit(1)
	}

	return contents
}

func main() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				data := obtainApiData()
				persistData(data)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}

func persistData(collection StationStateCollection) {
	// implement this
}

func obtainApiData() StationStateCollection {
	startTime := time.Now()
	apiData := doCurl()
	requestEndTime := time.Now()

	var stationCollection StationStateCollection

	err := xml.Unmarshal(apiData, &stationCollection)
	if err != nil {
		fmt.Printf("Unmarshal error: %v, structure :%v", err, apiData)
		return stationCollection
	}

	fmt.Printf("Data successfully received, request time: %v, unmarshalling time: %v\n", requestEndTime.Sub(startTime), time.Since(requestEndTime))
	return stationCollection
}

type StationStateCollection struct {
	StationStates []StationState `xml:"station"`
}

func (s StationStateCollection) Print() {
	for i := 0; i < len(s.StationStates); i++ {
		s.StationStates[i].Print()
	}
}

type StationState struct {
	// TODO review which of these fields need to be parsed and which not (we could potentially have different queries for the station state and the station data, as the second will change less frequently or may even not change at all)
	Id                int     `xml:"id"`
	Type              string  `xml:"type"`
	Latitude          float64 `xml:"lat"`
	Longitude         float64 `xml:"long"`
	Street            string  `xml:"street"`
	Height            int     `xml:"height"`
	StreetNumber      string  `xml:"streetNumber"` // Temporary, sometimes it is not set
	NearbyStationList string  `xml:"nearbyStationList"`
	Status            string  `xml:"status"`
	FreeSlots         int     `xml:"slots"`
	Bikes             int     `xml:"bikes"`
}

func (s StationState) Print() {
	fmt.Printf("Id : %v\n", s.Id)
	fmt.Printf("Type : %v\n", s.Type)
	fmt.Printf("Latitude : %v\n", s.Latitude)
	fmt.Printf("Longitude : %v\n", s.Longitude)
	fmt.Printf("Street : %v\n", s.Street)
	fmt.Printf("Height : %v\n", s.Height)
	fmt.Printf("StreetNumber : %v\n", s.StreetNumber)
	fmt.Printf("NearbyStationList : %v\n", s.NearbyStationList)
	fmt.Printf("Status : %v\n", s.Status)
	fmt.Printf("FreeSlots : %v\n", s.FreeSlots)
	fmt.Printf("Bikes : %v\n", s.Bikes)
}
