package handler

import (
	"encoding/json"
	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/service"
	"net/http"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	authService service.AuthService
	conf        *config.Config
}

func NewAuthHandler(authService service.AuthService, conf *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		conf:        conf,
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

	err = h.authService.Logout(r.Context(), refreshToken.(string))
	if err != nil {
		http.Error(w, "Failed to logout: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Logout successful"})
}

func (h *AuthHandler) LoginByGoogle(w http.ResponseWriter, r *http.Request) {
	u, err := h.authService.LoginByGoogle(r.Context())
	if err != nil {
		http.Error(w, "Failed to open google login window: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) LoginByGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	ticket, err := h.authService.LoginByGoogleCallback(r.Context(), state, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	callbackURL := strings.TrimRight(h.conf.OriginUrl, "/") +
		"/api/auth/google/callback?ticket=" + url.QueryEscape(ticket)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")

	http.Redirect(w, r, callbackURL, http.StatusSeeOther)
}

func (h *AuthHandler) RedeemGoogleOauthTicket(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket := body["ticket"]
	if ticket == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.RedeemGoogleOauthTicket(r.Context(), ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
