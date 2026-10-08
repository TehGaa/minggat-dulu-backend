package repository

import (
	"context"
	"encoding/json"
	"log"
	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/user/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	StoreToken(ctx context.Context, token model.Token) bool
	GetToken(ctx context.Context, key string) string
	DeleteToken(ctx context.Context, key string) bool
	CreateToken(ctx context.Context, userId, email string) model.Token
}

type AuthRepositoryImpl struct {
	client *mongo.Client
	conf   *config.Config
	r      *redis.Client
}

func NewAuthRepository(client *mongo.Client, conf *config.Config, r *redis.Client) AuthRepository {
	return &AuthRepositoryImpl{
		client: client,
		conf:   conf,
		r:      r,
	}
}

func (a *AuthRepositoryImpl) StoreToken(ctx context.Context, token model.Token) bool {

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return false
	}

	err = a.r.Set(ctx, "session:"+token.RefreshToken, tokenBytes, 24*time.Hour).Err()
	if err != nil {
		return false
	}

	return true
}

func (a *AuthRepositoryImpl) GetToken(ctx context.Context, key string) string {
	token, err := a.r.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return token
}

func (a *AuthRepositoryImpl) DeleteToken(ctx context.Context, key string) bool {
	err := a.r.Del(ctx, key).Err()
	if err != nil {
		return false
	}

	return true
}

func (a *AuthRepositoryImpl) CreateToken(ctx context.Context, userId, email string) model.Token {

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":   userId,
		"email": email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}).SignedString([]byte(a.conf.JwtKey))
	if err != nil {
		log.Fatalf("Error generating access token: %v", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":   userId,
		"email": email,
		"exp":   time.Now().Add(15 * time.Hour).Unix(),
	}).SignedString([]byte(a.conf.JwtKey))
	if err != nil {
		log.Fatalf("Error generating refresh token: %v", err)
	}

	token := model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AtExpires:    time.Now().Add(15 * time.Minute),
		RtExpires:    time.Now().Add(15 * time.Hour),
	}

	return token
}
