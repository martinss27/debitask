package main

import (
	"log"
	"net/http"

	"debitask/config"
	"debitask/handlers"
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

	mux := http.NewServeMux()
	routes.Register(mux, cfg.JWTSecret)

	log.Printf("server starting on http://localhost%s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("server error: %v\n", err)
	}
}
