package collection

import "fmt"

type StationStateCollection struct {
	StationStates []stationState `xml:"station"`
	Updatetime    int            `xml:"updatetime"`
}

func (s StationStateCollection) Print() {
	for i := 0; i < len(s.StationStates); i++ {
		s.StationStates[i].Print()
	}
}

type stationState struct {
	// TODO review which of these fields need to be parsed and which not (we could potentially have different queries for the station state and the station data, as the second will change less frequently or may even not change at all)
	ID                int     `xml:"id"`
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

func (s stationState) Print() {
	fmt.Printf("Id : %v\n", s.ID)
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
