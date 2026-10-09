package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"minggat-dulu-backend/internal/pkg/config"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	DeleteEmail(ctx context.Context, userId, email string) error
	GetUserByEmail(ctx context.Context, email string) bson.M
	StoreUser(ctx context.Context, user bson.M) bool
	CreateUser(ctx context.Context, email, password string) (bson.M, error)
}

type UserRepositoryImpl struct {
	client *mongo.Client
	conf   *config.Config
	r      *redis.Client
}

func NewUserRepository(client *mongo.Client, conf *config.Config, r *redis.Client) UserRepository {
	return &UserRepositoryImpl{
		client: client,
		conf:   conf,
		r:      r,
	}
}

func (u *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) bson.M {

	userCollection := u.client.Database(u.conf.DBName).Collection("users")

	var result bson.M
	err := userCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&result)

	if err != nil {
		log.Println("Error finding user in MongoDB:", err)
		return nil
	}

	return result
}

func (u *UserRepositoryImpl) DeleteEmail(ctx context.Context, userId, email string) error {

	err := u.r.Del(ctx, "user:"+userId+":email").Err()
	if err != nil {
		log.Println("Error deleting email from Redis:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) StoreUser(ctx context.Context, user bson.M) bool {
	oid, ok := user["_id"].(bson.ObjectID)
	var userId string
	if ok {
		userId = oid.Hex()
	} else {
		userId = user["_id"].(string)
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user data:", err)
		return false
	}

	err = u.r.Set(ctx, "user:"+userId, userBytes, 24*time.Hour).Err()
	if err != nil {
		log.Println("Error storing user in Redis:", err)
		return false
	}

	return true
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, email, password string) (bson.M, error) {
	user := bson.M{
		"email":    email,
		"password": password,
	}
	res, err := u.client.Database(u.conf.DBName).Collection("users").InsertOne(ctx, user)
	if err != nil {
		log.Println("Error creating new user", err)
		return nil, err
	}
	user["_id"] = res.InsertedID

	return user, nil
}
