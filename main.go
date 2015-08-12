package main

import (
	"encoding/xml"
	"fmt"

	curl "github.com/andelf/go-curl"
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

func doCurl() {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, "http://wservice.viabicing.cat/v1/getstations.php?v=1")

	// TODO find out how we can connect this to unmarshalling
	// make a callback function
	fooTest := func(buf []byte, userdata interface{}) bool {
		//println("DEBUG: size=>", len(buf))
		//        print(string(buf))
		return true
	}

	// this is most likely unnecessary, try to remove
	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func main() {
	apiData := apiFakeDataProvider()
	var stationCollection StationStateCollection

	err := xml.Unmarshal(apiData, &stationCollection)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	stationCollection.Print()
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
	StreetNumber      int     `xml:"streetNumber"`
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
