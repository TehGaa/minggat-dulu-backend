package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"minggat-dulu-backend/internal/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	userHandler "minggat-dulu-backend/internal/user/handler"
)

func main() {
	fmt.Println("User Service starting on :8082...")

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DBNAME")
	if uri == "" {
		uri = "mongodb://root:password@localhost:27017"
	}
	if dbName == "" {
		dbName = "minggat_dulu"
	}

	client := database.GetMongoClient(uri)
	userHandler := userHandler.NewUserHandler(client, dbName)

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
