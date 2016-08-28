package db

import (
	"database/sql"
	"fmt"
	"os"

	// loads the mysql driver
	_ "github.com/alexcarol/bicing-oracle/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
)

// GetRawDataDBFromEnv returns a sql driver obtained using environment variables
func GetRawDataDBFromEnv() (*sql.DB, error) {
	dbName := getEnv("MYSQL_RAW_DATA_NAME", "bicing_raw")

	username := getEnv("MYSQL_RAW_DATA_USER", "root")
	password := getEnv("MYSQL_RAW_DATA_PASSWORD", "")

	port := getEnv("MYSQL_RAW_DATA_ADDRESS", "localhost:3306")

	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, port, dbName))
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
