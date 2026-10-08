package service

import (
	"context"
	"log"

	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) model.Token
}

type authServiceImpl struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
}

func NewAuthService(
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
) AuthService {
	return &authServiceImpl{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (a *authServiceImpl) Login(ctx context.Context, email string, password string) model.Token {
	res := a.userRepository.GetUserByEmail(ctx, email)

	if res == nil {
		return model.Token{}
	}

	// err := bcrypt.CompareHashAndPassword([]byte(res["password"].(string)), []byte(password))

	// if err != nil {
	// 	return model.Token{}
	// }
	oid, ok := res["_id"].(bson.ObjectID)
	var userId string
	if ok {
		userId = oid.Hex()
	} else {
		userId = res["_id"].(string) // fallback in case it's actually a string
	}

	err := a.userRepository.StoreEmail(ctx, userId, email)
	if err != nil {
		return model.Token{}
	}

	token := a.authRepository.CreateToken(ctx, userId, email)
	if !a.authRepository.StoreToken(ctx, token) {
		err = a.userRepository.DeleteEmail(ctx, userId, email)
		if err != nil {
			log.Println("Error deleting email from Redis:", err)
		}
		return model.Token{}
	}

	return token
}
