package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"encoding/json"
	_ "github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/alexcarol/bicing-oracle/prediction"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
	"net/url"
	"strconv"
)

func main() {
	dbName := getEnv("MYSQL_RAW_DATA_NAME", "bicing_raw")

	username := getEnv("MYSQL_RAW_DATA_USER", "root")
	password := getEnv("MYSQL_RAW_DATA_PASSWORD", "")

	port := getEnv("MYSQL_RAW_DATA_ADDRESS", "localhost:3306")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, port, dbName))
	if err != nil {
		panic(err)
	}

	var stationProvider = repository.NewSQLStationProvider(db)

	http.HandleFunc("/checkup", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/prediction", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
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
