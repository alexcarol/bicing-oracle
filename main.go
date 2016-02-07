package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/go-sql-driver/mysql"

	"os"

	"strconv"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
	"github.com/alexcarol/bicing-oracle/station-state/parser"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
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

	storage := repository.NewSQLStorage(db)

	pollingTime, err := strconv.Atoi(getEnv("BICING_API_POLLING_TIME", "5"))
	if err != nil {
		panic("Error converting ascii to integer " + err.Error())
	}

	ticker := time.NewTicker(time.Duration(pollingTime) * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				apiData, err := datasource.APIData()
				if err != nil {
					fmt.Println(err)
					break
				}

				data, err := parser.ParseXML(apiData)
				if err != nil {
					fmt.Println("Error parsing xml")
					break
				}

				storage.PersistCollection(data)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
