package models

import "time"

type Habit struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Days      int       `json:"days"`
	Penalty   float64   `json:"penalty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// populated on list requests
	CheckedInToday bool `json:"checked_in_today"`
}

type HabitLog struct {
	ID         string    `json:"id"`
	HabitID    string    `json:"habit_id"`
	UserID     string    `json:"user_id"`
	Date       time.Time `json:"date"`
	CheckedIn  bool      `json:"checked_in"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateHabitRequest struct {
	Name    string  `json:"name"`
	Days    int     `json:"days"`
	Penalty float64 `json:"penalty"`
}
