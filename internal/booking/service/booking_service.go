package service

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingService interface {
	getUserBokingsByUserId(userId string) (bson.M, error)
}

type bookingServiceImpl struct {
}

func NewBookingService() BookingService {
	return &bookingServiceImpl{}
}

func (s *bookingServiceImpl) getUserBokingsByUserId(userId string) (bson.M, error) {
	return bson.M{}, nil
}
