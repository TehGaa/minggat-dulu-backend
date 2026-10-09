package config

import (
	"os"
	"sync"
)

type Config struct {
	MongoURI string
	DBName   string
	RedisUri string
	JwtKey   string
}

var (
	config   *Config
	syncOnce sync.Once
)

func GetConfig() *Config {
	syncOnce.Do(func() {
		mongoURI := os.Getenv("MONGODB_URI")
		dbName := os.Getenv("MONGODB_DBNAME")
		redisUri := os.Getenv("REDIS_URI")
		rabbitMQUri := os.Getenv("RABBITMQ_URI")
		jwtKey := os.Getenv("JWT_KEY")
		if mongoURI == "" {
			mongoURI = "mongodb://root:password@localhost:27017"
		}
		if dbName == "" {
			dbName = "minggat_dulu"
		}
		if redisUri == "" {
			redisUri = "redis://localhost:6379/0"
		}
		if rabbitMQUri == "" {
			rabbitMQUri = "amqp://guest:guest@localhost:5672/"
		}

		config = &Config{
			MongoURI: mongoURI,
			DBName:   dbName,
			RedisUri: redisUri,
			JwtKey:   jwtKey,
		}
	})

	return config
}
