package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"database/sql"
	"time"
	_ "github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
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
		var updatetime int64
		err := db.QueryRow("select UNIX_TIMESTAMP(last_updatetime) from station order by last_updatetime desc limit 1").Scan(&updatetime);
		if err != nil {
			http.Error(w, err.Error(), 500)

			return;
		}

		currentTime := time.Now().Unix()
		twoMinutesInThePast := currentTime - 120
		if twoMinutesInThePast > updatetime {
			http.Error(w, fmt.Sprintf("Time difference too big, current time : %d, update time: %d, difference: %d", currentTime, updatetime, currentTime - updatetime), 500)

			return;
		} else {
			fmt.Fprintf(w, "Time difference is reasonable, current time : %d, update time: %d, difference: %d", currentTime, updatetime, currentTime - updatetime)

			return;
		}
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
