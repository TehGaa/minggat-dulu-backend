package service

import (
	"context"
	"encoding/json"
	"minggat-dulu-backend/internal/user/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) bson.M
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (u *userServiceImpl) GetUserByEmail(ctx context.Context, email string) bson.M {
	user := u.userRepository.GetUserByEmail(ctx, email)

	data, err := json.Marshal(user)
	if err != nil {
		return nil
	}

	var userData bson.M
	err = json.Unmarshal(data, &userData)
	if err != nil {
		return nil
	}

	return userData
}
