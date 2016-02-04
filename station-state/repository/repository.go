package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexcarol/bicing-oracle/station-state/collection"
)

func NewStorage() StationStatePersister {
	return stationStateStorage{}
}

type stationStateStorage struct{}

type StationStatePersister interface {
	PersistCollection(collection.StationStateCollection)
}

func (storage stationStateStorage) PersistCollection(collection collection.StationStateCollection) {

}

type sqlStorage struct {
	database *sql.DB
}

func (storage sqlStorage) PersistCollection(collection collection.StationStateCollection) {
	collection.Print()
}

func NewSQLStorage(location string) (StationStatePersister, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s)/test", location))
	if err != nil {
		return nil, err
	}

	return sqlStorage{db}, nil
}
