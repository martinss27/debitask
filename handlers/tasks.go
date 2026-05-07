package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"debitask/middleware"
	"debitask/models"
	"debitask/store"
)

func userIDFromContext(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return "", false
	}
	return userID, true
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if req.Deadline.IsZero() {
		http.Error(w, "deadline is required", http.StatusBadRequest)
		return
	}

	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	task, err := store.CreateTask(userID, req.Title, req.Description, req.Deadline)
	if err != nil {
		log.Printf("CreateTask error (user %s): %v\n", userID, err)
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	task, err := store.GetTaskByID(id, userID)
	if err != nil {
		log.Printf("GetTask error (user %s, task %s): %v\n", userID, id, err)
		http.Error(w, "failed to get task", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	tasks, err := store.GetTasksByUser(userID)
	if err != nil {
		log.Printf("ListTasks error (user %s): %v\n", userID, err)
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	id := r.PathValue("id")
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != nil && *req.Status != models.TaskStatusCompleted {
		http.Error(w, "status can only be set to completed", http.StatusBadRequest)
		return
	}

	task, err := store.UpdateTask(id, userID, req.Title, req.Description, req.Deadline, req.Status)
	if err != nil {
		log.Printf("UpdateTask error (user %s, task %s): %v\n", userID, id, err)
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	deleted, err := store.DeleteTask(id, userID)
	if err != nil {
		log.Printf("DeleteTask error (user %s, task %s): %v\n", userID, id, err)
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
