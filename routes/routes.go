package routes

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"debitask/handlers"
	"debitask/middleware"
)

func Register(mux *http.ServeMux, jwtSecret string) {
	// 5 login attempts per minute per IP, burst of 3
	loginLimiter := middleware.NewRateLimiter(rate.Every(60*time.Second/5), 3)

	// Public routes
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/signup", handlers.Signup)
	mux.Handle("/api/login", loginLimiter.Limit(http.HandlerFunc(handlers.Login)))

	// Protected routes
	mux.Handle("/api/profile", middleware.AuthMiddleware(jwtSecret, http.HandlerFunc(handlers.Profile)))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
