package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"debitask/auth"
	"debitask/models"
	"debitask/store"
)

var JWTSecret string

// dummyHash is used to keep login response time constant when a user is not found,
// preventing timing-based user enumeration.
var dummyHash, _ = auth.HashPassword("dummy-timing-password-not-used")

func SetJWTSecret(secret string) {
	JWTSecret = secret
}

func Signup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	if !auth.ValidateEmail(req.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	if err := auth.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to process password", http.StatusInternalServerError)
		return
	}

	user, err := store.CreateUser(req.Email, passwordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, JWTSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	user, err := store.GetUserByEmail(req.Email)
	if err != nil {
		auth.VerifyPassword(dummyHash, req.Password) // constant-time: prevent user enumeration
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, JWTSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	})
}
