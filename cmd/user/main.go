package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/pkg/database"
	userHandler "minggat-dulu-backend/internal/user/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("User Service starting on :8082...")

	config := config.GetConfig()

	client := database.GetMongoClient(config.MongoURI)
	userHandler := userHandler.NewUserHandler(client, config.DBName)

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User Service OK"))
		})

		//TODO: add auth apis
		r.Route("/auth", func(r chi.Router) {

		})

		r.Group(func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Get("/", userHandler.GetUsers)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8082", r))
}
