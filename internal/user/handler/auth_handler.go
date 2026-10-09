package handler

import (
	"encoding/json"
	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/service"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body bson.M
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	email := body["email"]
	password := body["password"]
	if email == nil || password == nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), email.(string), password.(string))
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	var body bson.M
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	refreshToken := body["refresh_token"]

	if refreshToken == nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	token := h.authService.Refresh(r.Context(), refreshToken.(string))

	if token.AccessToken == "" {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
