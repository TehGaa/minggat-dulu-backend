package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserHandler struct {
	client  *mongo.Client
	dbName  string
}

func NewUserHandler(client *mongo.Client, dbName string) *UserHandler {
	return &UserHandler{
		client:  client,
		dbName:  dbName,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	userCollection := h.client.Database(h.dbName).Collection("users")

	// Create a new context from the request context with a 10 second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var result []bson.M

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
