package main

import (
	"time"

	"github.com/alexcarol/bicing-api/station-state/datasource"
	"github.com/alexcarol/bicing-api/station-state/parser"
	"github.com/alexcarol/bicing-api/station-state/repository"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	storage := repository.NewStorage()

	go func() {
		for {
			select {
			case <-ticker.C:
				data := parser.ParseXML(datasource.FixtureData())
				storage.PersistCollection(data)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}
