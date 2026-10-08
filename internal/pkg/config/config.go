package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	MongoURI      string
	DBName        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	JwtKey        string
}

var (
	config   *Config
	syncOnce sync.Once
)

func GetConfig() *Config {
	syncOnce.Do(func() {
		mongoURI := os.Getenv("MONGODB_URI")
		dbName := os.Getenv("MONGODB_DBNAME")
		redisAddr := os.Getenv("REDIS_ADDR")
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		jwtKey := os.Getenv("JWT_KEY")
		if mongoURI == "" {
			mongoURI = "mongodb://root:password@localhost:27017"
		}
		if dbName == "" {
			dbName = "minggat_dulu"
		}
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}
		if redisPassword == "" {
			redisPassword = ""
		}
		if redisDB == 0 {
			redisDB = 0
		}

		config = &Config{
			MongoURI:      mongoURI,
			DBName:        dbName,
			RedisAddr:     redisAddr,
			RedisPassword: redisPassword,
			RedisDB:       redisDB,
			JwtKey:        jwtKey,
		}
	})

	return config
}
