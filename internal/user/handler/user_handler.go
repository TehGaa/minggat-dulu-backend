package handler

import (
	"encoding/json"
	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/service"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewUserHandler(authService service.AuthService, userService service.UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var loginRequest model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res := h.authService.Login(
		r.Context(),
		loginRequest.Email,
		loginRequest.Password,
	)

	if res.AccessToken == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {

	var body bson.M
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := body["email"].(string)

	res := h.userService.GetUserByEmail(r.Context(), email)

	if res == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
