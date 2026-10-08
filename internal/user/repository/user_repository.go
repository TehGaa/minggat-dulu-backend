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
	StoreEmail(ctx context.Context, userId, email string) error
	DeleteEmail(ctx context.Context, userId, email string) error
	GetUserByEmail(ctx context.Context, email string) bson.M
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

	oid, ok := result["_id"].(bson.ObjectID)
	var userId string
	if ok {
		userId = oid.Hex()
	} else {
		userId = result["_id"].(string) // fallback in case it's actually a string
	}

	resultBytes, err := json.Marshal(result)
	if err == nil {
		err = u.r.Set(ctx, "user:"+userId, resultBytes, 24*time.Hour).Err()
		if err != nil {
			log.Println("Error storing user data in Redis:", err)
		}
	} else {
		log.Println("Error marshaling user data for Redis:", err)
	}

	return result
}

func (u *UserRepositoryImpl) StoreEmail(ctx context.Context, userId, email string) error {

	err := u.r.Set(ctx, "user:"+userId+":email", email, 24*time.Hour).Err()
	if err != nil {
		log.Println("Error storing email in Redis:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) DeleteEmail(ctx context.Context, userId, email string) error {

	err := u.r.Del(ctx, "user:"+userId+":email").Err()
	if err != nil {
		log.Println("Error deleting email from Redis:", err)
		return err
	}

	return nil
}
