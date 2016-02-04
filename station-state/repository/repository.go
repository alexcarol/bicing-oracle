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
	PersistCollection(collection.StationStateCollection) error
}

func (storage stationStateStorage) PersistCollection(collection collection.StationStateCollection) error {
	return nil
}

type sqlStorage struct {
	database *sql.DB
}

func (storage sqlStorage) PersistCollection(collection collection.StationStateCollection) error {
	transaction, err := storage.database.Begin()
	if nil != err {
		return err
	}
	for _, stationState := range collection.StationStates {
		_, stationInsertErr := transaction.Exec("insert into station values (?, ?, ?, ?, ?, ?, ?, FROM_UNIXTIME(?)) ON DUPLICATE KEY UPDATE last_updatetime = FROM_UNIXTIME(?);", stationState.ID, stationState.Latitude, stationState.Longitude, stationState.Street, stationState.Height, stationState.StreetNumber, stationState.NearbyStationList, collection.Updatetime, collection.Updatetime)
		if stationInsertErr != nil {
			fmt.Println("Error executing statement " + stationInsertErr.Error())
			rollErr := transaction.Rollback()
			if rollErr != nil {
				fmt.Println("Error doing rollback" + err.Error())
				return rollErr
			}

			return err
		}
		_, err := transaction.Exec("insert into station_state values (?, FROM_UNIXTIME(?), ?, ?);", stationState.ID, collection.Updatetime, stationState.FreeSlots, stationState.Bikes)
		if err != nil {
			fmt.Println("Error executing statement " + err.Error())
			rollErr := transaction.Rollback()
			if rollErr != nil {
				fmt.Println("Error doing rollback" + err.Error())
				return rollErr
			}

			return err
		}
	}

	return transaction.Commit()
}

func NewSQLStorage(db *sql.DB) StationStatePersister {
	return sqlStorage{db}
}
