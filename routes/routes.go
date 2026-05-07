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

	// Task routes
	mux.Handle("/api/tasks", middleware.AuthMiddleware(jwtSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateTask(w, r)
		case http.MethodGet:
			handlers.ListTasks(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/tasks/{id}", middleware.AuthMiddleware(jwtSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTask(w, r)
		case http.MethodPut:
			handlers.UpdateTask(w, r)
		case http.MethodDelete:
			handlers.DeleteTask(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
