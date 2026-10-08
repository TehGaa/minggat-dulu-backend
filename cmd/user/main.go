package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/pkg/database"
	customMiddleware "minggat-dulu-backend/internal/pkg/middleware"
	"minggat-dulu-backend/internal/user/handler"
	"minggat-dulu-backend/internal/user/repository"
	"minggat-dulu-backend/internal/user/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("User Service starting on :8082...")

	config := config.GetConfig()

	mongodbClient := database.GetMongoClient(config.MongoURI)
	redisClient := database.GetRedisClient(config.RedisAddr, config.RedisPassword, config.RedisDB)
	userRepository := repository.NewUserRepository(mongodbClient, config, redisClient)
	userService := service.NewUserService(userRepository)
	authRepository := repository.NewAuthRepository(mongodbClient, config, redisClient)
	authService := service.NewAuthService(authRepository, userRepository)

	userHandler := handler.NewUserHandler(authService, userService)

	defer func() {
		if err := mongodbClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	r := chi.NewRouter()

	customMiddleware.InitMiddleware()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User Service OK"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				userHandler.Login(w, r)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.AuthMiddleware)
			r.Route("/user", func(r chi.Router) {
				r.Get("/", userHandler.GetUserByEmail)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8082", r))
}
