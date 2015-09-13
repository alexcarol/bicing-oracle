package main

import (
	"time"

	"github.com/alexcarol/bicing-oracle/station-state/datasource"
	"github.com/alexcarol/bicing-oracle/station-state/parser"
	"github.com/alexcarol/bicing-oracle/station-state/repository"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	storage := repository.NewStorage()

	go func() {
		for {
			select {
			case <-ticker.C:
				data, err := parser.ParseXML(datasource.FixtureData())
				if err != nil {
					panic("Error parsing xml")
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
