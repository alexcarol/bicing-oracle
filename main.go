package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
	"github.com/alexcarol/bicing-oracle/station-state/parser"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
)

func main() {

	storage, err := repository.NewSQLStorage("mysql:3306")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(5 * time.Second)

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
