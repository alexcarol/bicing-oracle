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
			http.Error(w, output, 500)
		} else {
			fmt.Fprint(w, output)
		}
	})

	http.HandleFunc("/prediction", func(w http.ResponseWriter, r *http.Request) {
		//var i = r.URL.Query()

		a, err := prediction.GetPredictions()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		output, err := json.Marshal(a)
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
