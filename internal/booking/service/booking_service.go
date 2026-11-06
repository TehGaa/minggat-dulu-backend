package service

import (
	"context"
	"minggat-dulu-backend/internal/booking/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingService interface {
	GetUserBokingsByUserId(ctx context.Context, userId string) ([]bson.M, error)
}

type bookingServiceImpl struct {
	bookingRepository repository.BookingRepository
}

func NewBookingService(bookingRepository repository.BookingRepository) BookingService {
	return &bookingServiceImpl{
		bookingRepository: bookingRepository,
	}
}

func (s *bookingServiceImpl) GetUserBokingsByUserId(ctx context.Context, userId string) ([]bson.M, error) {
	bookings, err := s.bookingRepository.GetUserBokingsByUserId(ctx, userId)
	return bookings, err
}
