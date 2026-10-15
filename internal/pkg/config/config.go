package config

import (
	"encoding/json"
	"os"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	MongoURI          string
	DBName            string
	RedisUri          string
	JwtKey            string
	GoogleOauthConfig *oauth2.Config
	OriginUrl         string
}

type GoogleOAuthCredentials struct {
	Web struct {
		JavaScriptOrigins []string `json:"javascript_origins"`
	} `json:"web"`
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
		googleOauthConfigJson := os.Getenv("GOOGLE_OAUTH_SECRET")
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

		var googleOauthConfig *oauth2.Config
		oauthConfig, err := google.ConfigFromJSON(
			[]byte(googleOauthConfigJson),
			"email",
		)
		if err != nil {
			googleOauthConfig = nil
		} else {
			googleOauthConfig = oauthConfig
		}

		var originUrl string

		var credentials GoogleOAuthCredentials
		err = json.Unmarshal([]byte(googleOauthConfigJson), &credentials)
		if err == nil {
			originUrl = credentials.Web.JavaScriptOrigins[0]
		} else {
			originUrl = "http://localhost:3000/"
		}

		config = &Config{
			MongoURI:          mongoURI,
			DBName:            dbName,
			RedisUri:          redisUri,
			JwtKey:            jwtKey,
			GoogleOauthConfig: googleOauthConfig,
			OriginUrl:         originUrl,
		}
	})

	return config
}
