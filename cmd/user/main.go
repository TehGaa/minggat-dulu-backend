package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("User Service starting on :8081...")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User Service OK"))
		})

		r.Route("/auth", func(r chi.Router) {

		})

		r.Group(func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {

			})
		})
	})

	log.Fatal(http.ListenAndServe(":8081", r))
}
