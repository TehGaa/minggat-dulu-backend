package main

import (
	"fmt"
	"log"
	"net/http"

	customMiddleware "minggat-dulu-backend/internal/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("Booking Service starting on :8083...")

	// bookingHandler := handler.NewBookingHandler()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Booking Service OK"))
		})

		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.AuthMiddleware)
			r.Use(customMiddleware.ApiMiddleware)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// bookingHandler.getUserBookings(w, r)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8083", r))
}
