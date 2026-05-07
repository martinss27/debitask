package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"debitask/models"
	"debitask/store"
)

func CreateHabit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Days < 1 || req.Days > 127 {
		http.Error(w, "days must be a bitmask between 1 and 127", http.StatusBadRequest)
		return
	}
	if req.Penalty < 0 {
		http.Error(w, "penalty must be a positive value", http.StatusBadRequest)
		return
	}

	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	habit, err := store.CreateHabit(userID, req.Name, req.Days, req.Penalty)
	if err != nil {
		log.Printf("CreateHabit error (user %s): %v\n", userID, err)
		http.Error(w, "failed to create habit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

func ListHabits(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	habits, err := store.GetHabitsByUser(userID, today)
	if err != nil {
		log.Printf("ListHabits error (user %s): %v\n", userID, err)
		http.Error(w, "failed to list habits", http.StatusInternalServerError)
		return
	}

	if habits == nil {
		habits = []models.Habit{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

func DeleteHabit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	deleted, err := store.DeleteHabit(id, userID)
	if err != nil {
		log.Printf("DeleteHabit error (user %s, habit %s): %v\n", userID, id, err)
		http.Error(w, "failed to delete habit", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "habit not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CheckInHabit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := store.CheckInHabit(id, userID, today); err != nil {
		log.Printf("CheckInHabit error (user %s, habit %s): %v\n", userID, id, err)
		if err.Error() == "habit not found" {
			http.Error(w, "habit not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to check in", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
