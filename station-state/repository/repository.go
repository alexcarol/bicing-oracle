package repository

import (
	"database/sql"
	"fmt"
	"time"

	"math"
	"sort"

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
				fmt.Println("Error doing rollback" + rollErr.Error())

				return rollErr
			}

			return err
		}
		_, err := transaction.Exec("insert into station_state values (?, FROM_UNIXTIME(?), ?, ?, ?);", stationState.ID, collection.Updatetime, stationState.FreeSlots, stationState.Bikes, stationState.Status)
		if err != nil {
			fmt.Println("Error executing statement " + err.Error())
			rollErr := transaction.Rollback()
			if rollErr != nil {
				fmt.Println("Error doing rollback" + rollErr.Error())

				return rollErr
			}

			return err
		}
	}

	return transaction.Commit()
}

// StationProvider gives you a list of the nearby stations
type StationProvider interface {
	GetNearbyStations(lat, lon float64, minStations int) ([]Station, error)
	GetStationStateByInterval(stationID int, start time.Time, duration time.Duration) ([]StationState, error)
}

// Station contains info about a station
type Station struct {
	ID           int
	Type         string
	Street       string
	StreetNumber string
	Height       int
	Lon          float64
	Lat          float64
	Distance     float64
}

// StationState represents the state of a current station in a point in time
type StationState struct {
	ID, Bikes, Slots int
	Time             int64
}

// NewSQLStationProvider returns a StationStateProvider that uses mysql to retrieve the information
func NewSQLStationProvider(db *sql.DB) StationProvider {
	return sqlStorage{db}
}

func (storage sqlStorage) GetNearbyStations(lat float64, lon float64, minStations int) ([]Station, error) {
	rows, err := storage.database.Query("SELECT id, latitude, longitude, street, street_number, height  FROM station")
	if err != nil {
		return nil, err
	}

	var stationList = make([]Station, 0, 500) // TODO check if this can be adjusted

	defer rows.Close()
	rows.Columns()

	for rows.Next() {
		var currentStation Station

		err = rows.Scan(&currentStation.ID, &currentStation.Lat, &currentStation.Lon, &currentStation.Street, &currentStation.StreetNumber, &currentStation.Height)
		if err != nil {
			return nil, err
		}

		currentStation.Distance = distance(currentStation, lat, lon)
		stationList = append(stationList, currentStation)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	sort.Stable(byDistance(stationList))

	return stationList[:minStations], nil
}

func (storage sqlStorage) GetStationStateByInterval(stationID int, start time.Time, duration time.Duration) ([]StationState, error) {
	rows, err := storage.database.Query("SELECT id, bikes, slots, UNIX_TIMESTAMP(updatetime) FROM station_state WHERE id=? AND updatetime > ? AND updatetime < ?", stationID, start, start.Add(duration))
	if err != nil {
		return nil, err
	}

	var stationList = make([]StationState, 0, 10000)

	defer rows.Close()
	rows.Columns()

	for rows.Next() {
		var currentStation StationState

		err = rows.Scan(&currentStation.ID, &currentStation.Bikes, &currentStation.Slots, &currentStation.Time)
		if err != nil {
			return nil, err
		}

		stationList = append(stationList, currentStation)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return stationList, nil
}

func distance(s Station, lat, lon float64) float64 {
	latDistance := math.Abs(s.Lat - lat)
	lonDistance := math.Abs(s.Lon - lon)

	return math.Sqrt(latDistance*latDistance + lonDistance*lonDistance)
}

type byDistance []Station

func (s byDistance) Len() int {
	return len(s)
}

func (s byDistance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDistance) Less(i, j int) bool {
	return s[i].Distance < s[j].Distance
}
