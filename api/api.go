package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"encoding/csv"
	"encoding/json"
	"net/url"
	"strconv"

	_ "github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/alexcarol/bicing-oracle/db"
	"github.com/alexcarol/bicing-oracle/fitCalculator"
	"github.com/alexcarol/bicing-oracle/prediction"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
)

func main() {
	db, err := db.GetRawDataDBFromEnv()
	if err != nil {
		panic(err)
	}

	var stationProvider = repository.NewSQLStationProvider(db)

	scheduleAllFits(stationProvider)

	http.HandleFunc("/admin/calculateFit", getCalculateFitHandler())
	http.HandleFunc("/checkup", getHealthCheckHandler(db))

	http.HandleFunc("/dumpdata", getDumpDataHandler(stationProvider))

	http.HandleFunc("/prediction", getPredictionHandler(stationProvider))
	http.HandleFunc("/prediction/single", getSingleStationPredictionHandler(stationProvider))

	log.Fatal(http.ListenAndServe(":80", nil))
}

func getCalculateFitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query = r.URL.Query()

		stationID, err := parseRequestInt(query, "stationID")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		from, err := parseRequestInt(query, "from")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		to, err := parseRequestInt(query, "to")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		fitCalculator.CalculateFit(uint(stationID), int64(from), int64(to))
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
	}
}

func getHealthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output string
		var problem = false

		var updatetime int64
		err := db.QueryRow("select UNIX_TIMESTAMP(last_updatetime) from station order by last_updatetime desc limit 1").Scan(&updatetime)
		if err != nil {
			output += err.Error() + "\n"
			problem = true
		} else {
			currentTime := time.Now().Unix()
			twoMinutesInThePast := currentTime - 120
			if twoMinutesInThePast > updatetime {
				problem = true
				output += fmt.Sprintf("Time difference too big, current time : %d, update time: %d, difference: %d\n", currentTime, updatetime, currentTime-updatetime)
			} else {
				output += fmt.Sprintf("Time difference is reasonable, current time : %d, update time: %d, difference: %d", currentTime, updatetime, currentTime-updatetime)
			}
		}

		if problem {
			log.Println(output)
			http.Error(w, output, 500)
		} else {
			fmt.Fprint(w, output)
		}
	}
}

func getDumpDataHandler(stationProvider repository.StationProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query = r.URL.Query()

		timestamp, err := parseRequestInt(query, "time")
		if err != nil {
			output := "Error parsing time"
			http.Error(w, output, 400)
			return
		}

		stationID, err := parseRequestInt(query, "station_id")
		if err != nil {
			output := "Error parsing stationId"
			http.Error(w, output, 400)
			return
		}

		startTime := time.Unix(int64(timestamp), 0)

		durationInSeconds, err := parseRequestInt(query, "duration")
		if err != nil {
			output := "Error parsing stationId"
			http.Error(w, output, 400)
			return
		}
		var duration = time.Duration(durationInSeconds) * time.Second

		var states []repository.StationState
		states, err = stationProvider.GetStationStateByInterval(stationID, startTime, duration)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writer := csv.NewWriter(w)

		writer.Write([]string{"ID", "bikes", "slots", "time"})
		for _, state := range states {
			writer.Write([]string{strconv.Itoa(state.ID), strconv.Itoa(state.Bikes), strconv.Itoa(state.Slots), strconv.Itoa(int(state.Time))})
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func getPredictionHandler(stationProvider repository.StationProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query = r.URL.Query()
		timestamp, err := parseRequestInt(query, "time")
		if err != nil {
			output := "Error parsing time"
			http.Error(w, output, 400)
			return
		}

		lat, err := parseRequestFloat64(query, "lat")
		if err != nil {
			http.Error(w, "Error parsing lat", 400)
			return
		}

		lon, err := parseRequestFloat64(query, "lon")
		if err != nil {
			http.Error(w, "Error parsing lon", 400)
			return
		}

		predictions, err := prediction.GetPredictions(timestamp, lat, lon, stationProvider)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		predictionMap := map[string]([]prediction.Prediction){
			"stations": predictions,
		}
		output, err := json.Marshal(predictionMap)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprint(w, string(output))
	}
}

func getSingleStationPredictionHandler(stationProvider repository.StationProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query = r.URL.Query()
		timestamp, err := parseRequestInt(query, "time")
		if err != nil {
			output := "Error parsing time"
			http.Error(w, output, 400)
			return
		}

		stationID, err := parseRequestInt(query, "stationID")
		if err != nil {
			http.Error(w, "Error parsing stationID", 400)
			return
		}

		singlePrediction, err := prediction.GetStationPrediction(timestamp, uint(stationID), stationProvider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting prediction: %s", err.Error()), 500)
			return
		}

		output, err := json.Marshal(singlePrediction)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprint(w, string(output))
	}
}

func scheduleAllFits(stationProvider repository.StationProvider) {
	stations, err := stationProvider.GetAllStations()
	if err != nil {
		panic(err)
	}

	for _, station := range stations {
		fitCalculator.ScheduleCalculate(station.ID)
	}
}

func parseRequestInt(query url.Values, name string) (int, error) {
	result, err := strconv.Atoi(query.Get(name))
	if err != nil {
		log.Println(err)
	}

	return result, err
}

func parseRequestFloat64(query url.Values, name string) (float64, error) {
	result, err := strconv.ParseFloat(query.Get(name), 64)
	if err != nil {
		log.Println(err)
	}

	return result, err
}
