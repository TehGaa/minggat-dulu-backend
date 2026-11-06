package handler

import (
	"encoding/json"
	"minggat-dulu-backend/internal/booking/service"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	var body bson.M
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := body["userId"].(string)

	bookings, err := h.bookingService.GetUserBokingsByUserId(r.Context(), userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}
