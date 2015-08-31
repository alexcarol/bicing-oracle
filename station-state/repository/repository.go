package repository

import (
	"github.com/alexcarol/bicing-api/station-state/collection"
)


func NewStorage() *StationStateStorage {
	return &StationStateStorage{}
}

type StationStateStorage struct{}

func (storage *StationStateStorage) PersistCollection(collection collection.StationStateCollection) {

}
