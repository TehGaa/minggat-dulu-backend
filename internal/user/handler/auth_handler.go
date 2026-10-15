package handler

import (
	"encoding/json"
	"minggat-dulu-backend/internal/pkg/config"
	"minggat-dulu-backend/internal/user/model"
	"minggat-dulu-backend/internal/user/service"
	"net/http"
	"net/url"

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

	token, err := h.authService.LoginByGoogleCallback(r.Context(), state, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	})

	fragment := url.Values{}
	fragment.Set("access_token", token.AccessToken)
	fragment.Set("refresh_token", token.RefreshToken)

	http.Redirect(w, r, h.conf.OriginUrl+"#"+fragment.Encode(), http.StatusSeeOther)
}
