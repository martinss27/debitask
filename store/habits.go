package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"debitask/models"
)

func CreateHabit(userID, name string, days int, penalty float64) (*models.Habit, error) {
	habit := &models.Habit{}
	err := DB.QueryRow(`
		INSERT INTO habits (user_id, name, days, penalty)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, days, penalty, created_at, updated_at`,
		userID, name, days, penalty,
	).Scan(&habit.ID, &habit.UserID, &habit.Name, &habit.Days, &habit.Penalty, &habit.CreatedAt, &habit.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}
	return habit, nil
}

func GetHabitsByUser(userID string, today time.Time) ([]models.Habit, error) {
	rows, err := DB.Query(`
		SELECT h.id, h.user_id, h.name, h.days, h.penalty, h.created_at, h.updated_at,
		       COALESCE(hl.checked_in, false) AS checked_in_today
		FROM habits h
		LEFT JOIN habit_logs hl ON hl.habit_id = h.id AND hl.date = $2
		WHERE h.user_id = $1
		ORDER BY h.created_at ASC`,
		userID, today.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}
	defer rows.Close()

	var habits []models.Habit
	for rows.Next() {
		var h models.Habit
		if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &h.Days, &h.Penalty, &h.CreatedAt, &h.UpdatedAt, &h.CheckedInToday); err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habits = append(habits, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate habits: %w", err)
	}
	return habits, nil
}

func DeleteHabit(id, userID string) (bool, error) {
	result, err := DB.Exec(`DELETE FROM habits WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return false, fmt.Errorf("failed to delete habit: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get affected rows: %w", err)
	}
	return rows > 0, nil
}

// CheckInHabit records a check-in for today. Uses ON CONFLICT to handle duplicate check-ins gracefully.
func CheckInHabit(habitID, userID string, today time.Time) error {
	var ownerID string
	err := DB.QueryRow(`SELECT user_id FROM habits WHERE id = $1`, habitID).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("habit not found")
	}
	if err != nil {
		return fmt.Errorf("failed to verify habit: %w", err)
	}
	if ownerID != userID {
		return fmt.Errorf("habit not found")
	}

	_, err = DB.Exec(`
		INSERT INTO habit_logs (habit_id, user_id, date, checked_in)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (habit_id, date) DO UPDATE SET checked_in = true`,
		habitID, userID, today.Format("2006-01-02"),
	)
	if err != nil {
		return fmt.Errorf("failed to check in: %w", err)
	}
	return nil
}

func GetHabitLogs(habitID, userID string, from, to time.Time) ([]models.HabitLog, error) {
	rows, err := DB.Query(`
		SELECT id, habit_id, user_id, date, checked_in, created_at
		FROM habit_logs
		WHERE habit_id = $1 AND user_id = $2 AND date >= $3 AND date <= $4
		ORDER BY date ASC`,
		habitID, userID, from.Format("2006-01-02"), to.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.HabitLog
	for rows.Next() {
		var l models.HabitLog
		if err := rows.Scan(&l.ID, &l.HabitID, &l.UserID, &l.Date, &l.CheckedIn, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit log: %w", err)
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate habit logs: %w", err)
	}
	return logs, nil
}

// MarkMissedHabits runs at the start of each day and records missed check-ins for yesterday.
func MarkMissedHabits(yesterday time.Time) (int64, error) {
	// bit position for yesterday's weekday (0=Monday)
	weekday := int(yesterday.Weekday()+6) % 7
	bit := 1 << weekday

	result, err := DB.Exec(`
		INSERT INTO habit_logs (habit_id, user_id, date, checked_in)
		SELECT h.id, h.user_id, $1, false
		FROM habits h
		WHERE (h.days & $2) > 0
		  AND NOT EXISTS (
		    SELECT 1 FROM habit_logs hl
		    WHERE hl.habit_id = h.id AND hl.date = $1
		  )
		ON CONFLICT (habit_id, date) DO NOTHING`,
		yesterday.Format("2006-01-02"), bit,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to mark missed habits: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}
	return count, nil
}
