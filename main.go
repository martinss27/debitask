package main

import (
	"log"
	"net/http"
	"time"

	"debitask/config"
	"debitask/handlers"
	"debitask/jobs"
	"debitask/routes"
	"debitask/store"
)

func main() {
	cfg := config.Load()

	if err := store.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("database connection failed: %v\n", err)
	}
	defer store.Close()

	handlers.SetJWTSecret(cfg.JWTSecret)

	jobs.StartOverdueChecker(24 * time.Hour)
	jobs.StartHabitMissChecker(24 * time.Hour)

	mux := http.NewServeMux()
	routes.Register(mux, cfg.JWTSecret)

	log.Printf("server starting on http://localhost%s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("server error: %v\n", err)
	}
}
