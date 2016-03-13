package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexcarol/bicing-oracle/station-state/collection"
)

type stationStateStorage struct{}

// StationStatePersister takes the data from a collection.StationStateCollection and saves it to a persistent storage
type StationStatePersister interface {
	PersistCollection(collection.StationStateCollection) error
}

// NewSQLStorage returns a StationStaterPersister that will persist data in the database/sql passed to it
func NewSQLStorage(db *sql.DB) StationStatePersister {
	db.Exec("CREATE TABLE IF NOT EXISTS `station` ( `id` int(11) NOT NULL, `latitude` float DEFAULT NULL, `longitude` float DEFAULT NULL, `street` varchar(255) DEFAULT NULL, `height` int(11) DEFAULT NULL, `street_number` varchar(255) DEFAULT NULL, `nearby_station_list` varchar(255) DEFAULT NULL, `last_updatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	db.Exec("CREATE TABLE IF NOT EXISTS `station_state` ( `id` int(11) NOT NULL DEFAULT '0', `updatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `slots` int(11) DEFAULT NULL, `bikes` int(11) DEFAULT NULL, PRIMARY KEY (`id`,`updatetime`), CONSTRAINT `station_state_ibfk_1` FOREIGN KEY (`id`) REFERENCES `station` (`id`) ON DELETE CASCADE ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	return sqlStorage{db}
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
		_, err := transaction.Exec("insert into station_state values (?, FROM_UNIXTIME(?), ?, ?, ?);", stationState.ID, collection.Updatetime, stationState.FreeSlots, stationState.Bikes, stationState.Status)
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
