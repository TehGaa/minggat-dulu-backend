package repository

import "go.mongodb.org/mongo-driver/v2/bson"

type BookingRepository interface {
}

type bookingRepositoryImpl struct {
}

func NewBookingService() BookingRepository {
	return &bookingRepositoryImpl{}
}

func (b *bookingRepositoryImpl) getUserBokingsByUserId(userId string) (bson.M, error) {
	return bson.M{}, nil
}
