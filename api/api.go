package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"encoding/json"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
	"github.com/alexcarol/bicing-oracle/station-state/parser"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/stations", getAllStations)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllStations(w http.ResponseWriter, r *http.Request) {
	apiData := datasource.FixtureData()

	stationStateCollection, error := parser.ParseXML(apiData)
	if error != nil {
		fmt.Fprintf(w, "Error executing %q", html.EscapeString(r.URL.Path))
		fmt.Print(error)
	}

	json.NewEncoder(w).Encode(stationStateCollection)
}
