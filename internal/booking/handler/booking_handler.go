package handler

import (
	"net/http"
)

type BookingHandler struct {
}

func NewBookingHandler() *BookingHandler {
	return &BookingHandler{}
}

func (h *BookingHandler) getUserBookings(w http.ResponseWriter, r *http.Request) {
	//TODO
}
