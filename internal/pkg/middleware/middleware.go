package middleware

import (
	"fmt"
	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/pkg/database"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type UserSession struct {
	UserID    string `json:"_id"`
	UserEmail string `json:"email"`
}

var (
	redisClient *redis.Client
	conf        *config.Config
	syncOnce    sync.Once
)

func InitMiddleware() {
	syncOnce.Do(func() {
		conf = config.GetConfig()
		redisClient = database.GetRedisClient(conf.RedisAddr, conf.RedisPassword, conf.RedisDB)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header not found", http.StatusUnauthorized)
			return
		}

		accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(conf.JwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		_id, ok := claims["_id"].(string)
		if !ok || _id == "" {
			http.Error(w, "Token payload missing user id", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Token payload missing email", http.StatusUnauthorized)
			return
		}

		var accessUUID string
		uuidClaim, ok := claims["access_uuid"].(string)
		if ok {
			accessUUID = uuidClaim
		} else {
			http.Error(w, "Token payload missing access_uuid", http.StatusUnauthorized)
			return
		}

		blacklisted, err := redisClient.Get(r.Context(), "auth:blacklist:"+accessUUID).Result()
		if err == nil && blacklisted == "1" {
			http.Error(w, "Token has been compromised", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
