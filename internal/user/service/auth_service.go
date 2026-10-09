package service

import (
	"context"
	"errors"
	"log"

	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) model.Token
	Register(ctx context.Context, email, password string) (model.Token, error)
	Refresh(ctx context.Context, refreshToken string) model.Token
}

type authServiceImpl struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
	conf           *config.Config
}

func NewAuthService(
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
	conf *config.Config,
) AuthService {
	return &authServiceImpl{
		authRepository: authRepository,
		userRepository: userRepository,
		conf:           conf,
	}
}

func (a *authServiceImpl) Login(ctx context.Context, email string, password string) model.Token {
	res := a.userRepository.GetUserByEmail(ctx, email)

	if res == nil {
		return model.Token{}
	}

	err := bcrypt.CompareHashAndPassword([]byte(res["password"].(string)), []byte(password))

	if err != nil {
		return model.Token{}
	}

	oid, ok := res["_id"].(bson.ObjectID)
	var userId string
	if ok {
		userId = oid.Hex()
	} else {
		userId = res["_id"].(string)
	}

	token := a.authRepository.CreateToken(ctx, userId, email)
	if token.AccessToken == "" {
		return model.Token{}
	}

	if !a.authRepository.StoreToken(ctx, token) {
		return model.Token{}
	}

	if !a.userRepository.StoreUser(ctx, res) {
		a.authRepository.DeleteToken(ctx, userId, email)
		return model.Token{}
	}

	return token
}

func (a *authServiceImpl) Register(ctx context.Context, email, password string) (model.Token, error) {
	user := a.userRepository.GetUserByEmail(ctx, email)

	if user != nil {
		return model.Token{}, errors.New("User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Token{}, errors.New("Error creating hashed password")
	}

	user, err = a.userRepository.CreateUser(ctx, email, string(hashedPassword))
	if err != nil {
		return model.Token{}, err
	}

	oid, ok := user["_id"].(bson.ObjectID)
	var userId string
	if ok {
		userId = oid.Hex()
	} else {
		userIdStr, ok := user["_id"].(string)
		if !ok {
			return model.Token{}, errors.New("Invalid or missing _id in created user")
		}
		userId = userIdStr
	}

	token := a.authRepository.CreateToken(ctx, userId, email)
	if token.AccessToken == "" {
		return model.Token{}, errors.New("Error creating token")
	}

	if !a.authRepository.StoreToken(ctx, token) {
		return model.Token{}, errors.New("Error storing token into Redis")
	}

	if !a.userRepository.StoreUser(ctx, user) {
		return model.Token{}, errors.New("Error storing user data into Redis")
	}

	return token, nil
}

func (h *authServiceImpl) Refresh(ctx context.Context, refreshToken string) model.Token {
	oldRefreshToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.conf.JwtKey), nil
	})
	if err != nil {
		log.Println("Invalid refresh token")
		return model.Token{}
	}

	claims := oldRefreshToken.Claims.(jwt.MapClaims)

	userId, ok := claims["_id"].(string)
	if !ok {
		log.Println("Invalid _id in claims")
		return model.Token{}
	}

	oldRefreshUUID, ok := claims["refresh_uuid"].(string)
	if !ok {
		log.Println("Invalid refresh_uuid in claims")
		return model.Token{}
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		log.Println("Invalid email in claims")
		return model.Token{}
	}

	if !h.authRepository.CheckTokenExist(ctx, userId, oldRefreshUUID) {
		log.Println("Error getting token from Redis")
		return model.Token{}
	}

	err = h.authRepository.BlacklistToken(ctx, userId, oldRefreshUUID)

	if err != nil {
		log.Println("Error blacklisting token in Redis")
		return model.Token{}
	}

	newToken := h.authRepository.CreateToken(ctx, userId, userEmail)

	if newToken.AccessToken == "" {
		log.Println("Error creating token")
		return model.Token{}
	}

	if !h.authRepository.DeleteToken(ctx, userId, oldRefreshUUID) {
		log.Println("Error deleting token in Redis")
		return model.Token{}
	}

	if !h.authRepository.StoreToken(ctx, newToken) {
		log.Println("Error storing token in Redis")
		return model.Token{}
	}

	res := h.userRepository.GetUserByEmail(ctx, claims["email"].(string))

	if res == nil {
		return model.Token{}
	}

	if !h.userRepository.StoreUser(ctx, res) {
		log.Println("Error storing user in Redis")
		if !h.authRepository.DeleteToken(ctx, userId, claims["email"].(string)) {
			log.Println("Error delete token in Redis")
		}
		return model.Token{}
	}

	return newToken
}
