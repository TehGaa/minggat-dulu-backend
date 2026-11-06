package repository

import (
	"context"
	"minggat-dulu-backend/internal/pkg/config"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BookingRepository interface {
	GetUserBokingsByUserId(ctx context.Context, userId string) ([]bson.M, error)
}

type bookingRepositoryImpl struct {
	client *mongo.Client
	conf   *config.Config
	r      *redis.Client
}

func NewBookingService(client *mongo.Client, conf *config.Config, r *redis.Client) BookingRepository {
	return &bookingRepositoryImpl{
		client: client,
		conf:   conf,
		r:      r,
	}
}

func (b *bookingRepositoryImpl) GetUserBokingsByUserId(ctx context.Context, userId string) ([]bson.M, error) {
	userObjectID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	cursor, err := b.client.Database(b.conf.DBName).Collection("bookings").Find(ctx, bson.M{
		"user_id": userObjectID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []bson.M
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
