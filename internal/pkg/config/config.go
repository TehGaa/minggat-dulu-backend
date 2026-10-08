package config

import (
	"os"
	"sync"
)

type Config struct {
	MongoURI string
	DBName   string
}

var (
	config   *Config
	syncOnce sync.Once
)

func GetConfig() *Config {
	syncOnce.Do(func() {
		mongoURI := os.Getenv("MONGODB_URI")
		dbName := os.Getenv("MONGODB_DBNAME")
		if mongoURI == "" {
			mongoURI = "mongodb://root:password@localhost:27017"
		}
		if dbName == "" {
			dbName = "minggat_dulu"
		}

		config = &Config{
			MongoURI: mongoURI,
			DBName:   dbName,
		}
	})

	return config
}
