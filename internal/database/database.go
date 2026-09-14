package database

import (
	"database/sql"
	"os"
	"strconv"
	"sync"

	"github.com/lib/pq"
)

var (
	instance *sql.DB
	once     sync.Once
	err      error
)

func GetInstance() (*sql.DB, error) {
	once.Do(func() {
		port := uint16(5432)
		sslMode := os.Getenv("UCT_DB_SSLMODE")
		if sslMode == "" {
			sslMode = "disable"
		}
		if configuredPort := os.Getenv("UCT_DB_PORT"); configuredPort != "" {
			parsedPort, parseErr := strconv.ParseUint(configuredPort, 10, 16)
			if parseErr != nil {
				err = parseErr
				return
			}
			port = uint16(parsedPort)
		}

		cfg := pq.Config{
			User:     os.Getenv("UCT_DB_USER"),
			Password: os.Getenv("UCT_DB_PASSWORD"),
			Database: os.Getenv("UCT_DB_NAME"),
			Host:     os.Getenv("UCT_DB_HOST"),
			Port:     port,
			SSLMode:  pq.SSLMode(sslMode),
		}

		c, connectorErr := pq.NewConnectorConfig(cfg)
		if connectorErr != nil {
			err = connectorErr
			return
		}
		instance = sql.OpenDB(c)
		instance.SetMaxOpenConns(25)
		instance.SetMaxIdleConns(10)

		err = instance.Ping()
		if err != nil {
			instance.Close()
			instance = nil
		}
	})
	return instance, err
}
