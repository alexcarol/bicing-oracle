package repository

import (
	"database/sql"
	"github.com/alexcarol/bicing-oracle/weather/api"
)

type weatherPersister interface {
	PersistWeather(api.Weather) error
}

type sqlStorage struct {
	database *sql.DB
}

func (storage sqlStorage) PersistWeather(w api.Weather) error {

	_, err := storage.database.Exec("insert into weather values (?, ?, ?, ?, ?, FROM_UNIXTIME(?));", w.Type, w.Temperature, w.CloudPercentage, w.WindDegree, w.WindSpeed, w.Time)

	return err
}

func NewSQLStorage(db *sql.DB) weatherPersister {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `weather` ( `type` int(11) NOT NULL, `temperature` float DEFAULT NULL, `cloud_percentage` int(11) NOT NULL, `wind_degree` float DEFAULT NULL, `wind_speed` float DEFAULT NULL, `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`type`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		panic(err);
	}

	return sqlStorage{db}
}