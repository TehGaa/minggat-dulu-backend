package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/user/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/google/uuid"
)

type AuthRepository interface {
	StoreToken(ctx context.Context, token model.Token) bool
	DeleteToken(ctx context.Context, userId string, oldRefreshUUID string) bool
	CreateToken(ctx context.Context, userId, email string, user bson.M) model.Token
	CheckTokenExist(ctx context.Context, userId string, oldRefreshUUID string) bool
	BlacklistToken(ctx context.Context, userId string, oldRefreshUUID string) error
	StoreGoogleOauthState(ctx context.Context, state string) error
	CheckGoogleOauthStateExist(ctx context.Context, state string) bool
	StoreGoogleOauthTicket(ctx context.Context, ticket string, token model.Token) bool
	RedeemGoogleOauthTicket(ctx context.Context, ticket string) (model.Token, error)
}

type authRepositoryImpl struct {
	client *mongo.Client
	conf   *config.Config
	r      *redis.Client
}

func NewAuthRepository(client *mongo.Client, conf *config.Config, r *redis.Client) AuthRepository {
	return &authRepositoryImpl{
		client: client,
		conf:   conf,
		r:      r,
	}
}

func (a *authRepositoryImpl) StoreToken(ctx context.Context, token model.Token) bool {

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return false
	}

	err = a.r.Set(ctx, "session:"+token.UserId+":"+token.RefreshUUID, tokenBytes, 24*time.Hour).Err()
	if err != nil {
		return false
	}

	return true
}

func (a *authRepositoryImpl) CheckTokenExist(ctx context.Context, userId string, oldRefreshUUID string) bool {
	_, err := a.r.Get(ctx, "session:"+userId+":"+oldRefreshUUID).Result()
	if err != nil {
		return false
	}

	return true
}

func (a *authRepositoryImpl) DeleteToken(ctx context.Context, userId string, oldRefreshUUID string) bool {
	err := a.r.Del(ctx, "session:"+userId+":"+oldRefreshUUID).Err()
	if err != nil {
		return false
	}

	return true
}

func (a *authRepositoryImpl) CreateToken(ctx context.Context, userId, email string, user bson.M) model.Token {

	accessUUID := uuid.New().String()
	refreshUUID := uuid.New().String()

	atExpires := time.Now().Add(15 * time.Minute)
	rtExpires := time.Now().Add(24 * time.Hour)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":         userId,
		"email":       email,
		"access_uuid": accessUUID,
		"exp":         atExpires.Unix(),
	}).SignedString([]byte(a.conf.JwtKey))
	if err != nil {
		log.Fatalf("Error generating access token: %v", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":          userId,
		"email":        email,
		"refresh_uuid": refreshUUID,
		"exp":          rtExpires.Unix(),
	}).SignedString([]byte(a.conf.JwtKey))
	if err != nil {
		log.Fatalf("Error generating refresh token: %v", err)
	}

	token := model.Token{
		UserId: userId,
		UserData: bson.M{
			"username": user["username"],
			"avatar":   user["avatar"],
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessUUID:   accessUUID,
		RefreshUUID:  refreshUUID,
		AtExpires:    atExpires,
		RtExpires:    rtExpires,
	}

	return token
}

func (a *authRepositoryImpl) BlacklistToken(ctx context.Context, userId, oldRefreshUUID string) error {
	tokenRaw, err := a.r.Get(ctx, "session:"+userId+":"+oldRefreshUUID).Result()

	if err != nil {
		log.Println("Error getting token from Redis")
		return err
	}

	var token model.Token
	err = json.Unmarshal([]byte(tokenRaw), &token)

	if err != nil {
		log.Println("Error parsing token string")
		return err
	}

	accessToken, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.conf.JwtKey), nil
	})

	if err != nil {
		log.Println("Error invalid token")
	}

	accessClaims := accessToken.Claims.(jwt.MapClaims)

	accessExp, ok := accessClaims["exp"].(float64)
	if !ok {
		log.Println("Invalid exp claim in access token")
		return fmt.Errorf("invalid exp claim in access token")
	}
	accessTimeRemaining := time.Until(time.Unix(int64(accessExp), 0))

	if accessTimeRemaining <= 0 {
		return nil
	}

	key := "auth:blacklist:" + accessClaims["access_uuid"].(string)

	err = a.r.Set(ctx, key, "1", accessTimeRemaining).Err()

	if err != nil {
		log.Println("Error blacklisting access token")
		return err
	}

	return err
}

func (a *authRepositoryImpl) StoreGoogleOauthState(ctx context.Context, state string) error {

	err := a.r.Set(ctx, "google-auth:"+state, 1, 60*time.Second).Err()
	if err != nil {
		log.Println("Failed to store google oauth state: " + err.Error())
	}
	return err
}

func (a *authRepositoryImpl) CheckGoogleOauthStateExist(ctx context.Context, state string) bool {
	_, err := a.r.Get(ctx, "google-auth:"+state).Result()
	if err != nil {
		return false
	}
	err = a.r.Del(ctx, "google-auth:"+state).Err()
	if err != nil {
		return false
	}
	return true
}

func (a *authRepositoryImpl) StoreGoogleOauthTicket(ctx context.Context, ticket string, token model.Token) bool {

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return false
	}

	err = a.r.Set(ctx, "google-ticket:"+ticket, tokenBytes, 60*time.Second).Err()
	if err != nil {
		return false
	}

	return true
}

func (a *authRepositoryImpl) RedeemGoogleOauthTicket(ctx context.Context, ticket string) (model.Token, error) {
	tokenRaw, err := a.r.Get(ctx, "google-ticket:"+ticket).Result()
	if err != nil {
		return model.Token{}, errors.New("Failed to get token")
	}

	err = a.r.Del(ctx, "google-ticket:"+ticket).Err()
	if err != nil {
		return model.Token{}, errors.New("Failed to delete ticket")
	}

	var token model.Token
	err = json.Unmarshal([]byte(tokenRaw), &token)
	if err != nil {
		return model.Token{}, errors.New("Token is invalid")
	}

	return token, nil
}
