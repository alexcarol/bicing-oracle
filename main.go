package main

import (
	"fmt"
	"time"

	"os"

	"strconv"

	"log"

	"github.com/alexcarol/bicing-oracle/db"
	"github.com/alexcarol/bicing-oracle/station-state/datasource"
	"github.com/alexcarol/bicing-oracle/station-state/parser"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
	weatherDatasource "github.com/alexcarol/bicing-oracle/weather/datasource"
	weatherRepository "github.com/alexcarol/bicing-oracle/weather/repository"
)

func main() {
	db, err := db.GetRawDataDBFromEnv()
	if err != nil {
		panic(err)
	}

	storage := repository.NewSQLStorage(db)

	weatherStorage := weatherRepository.NewSQLStorage(db)

	pollingTime, err := strconv.Atoi(getEnv("BICING_API_POLLING_TIME", "45"))
	if err != nil {
		panic("Error converting ascii to integer " + err.Error())
	}

	ticker := time.NewTicker(time.Duration(pollingTime) * time.Second)

	var dataProvider datasource.BicingDataProvider
	if getEnv("BICING_API_FETCH_REAL_DATA", "1") == "1" {
		dataProvider = datasource.ProvideAPIData
	} else {
		dataProvider = datasource.ProvideFakeData
	}

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				weather, err := weatherDatasource.GetWeatherData()
				if err != nil {
					log.Print("Error getting weather data")
					log.Println(err)
				}

				err = weatherStorage.PersistWeather(weather)
				if err != nil {
					log.Println("Error persisting weather", err)
				}

				apiData, err := dataProvider()
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
