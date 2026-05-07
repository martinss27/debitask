package handlers

import (
	"encoding/json"
	"net/http"

	"debitask/middleware"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	email := r.Context().Value(middleware.EmailKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
		"email":   email,
	})
}
