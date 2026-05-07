package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"debitask/models"
)

func CreateTask(userID, title, description string, deadline time.Time) (*models.Task, error) {
	task := &models.Task{}
	desc := sql.NullString{String: description, Valid: description != ""}
	err := DB.QueryRow(`
		INSERT INTO tasks (user_id, title, description, deadline)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, COALESCE(description, ''), deadline, status, created_at, updated_at`,
		userID, title, desc, deadline,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Deadline, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func GetTasksByUser(userID string) ([]models.Task, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, title, COALESCE(description, ''), deadline, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY deadline ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Deadline, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}
	return tasks, nil
}

func GetTaskByID(id, userID string) (*models.Task, error) {
	task := &models.Task{}
	err := DB.QueryRow(`
		SELECT id, user_id, title, COALESCE(description, ''), deadline, status, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Deadline, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func UpdateTask(id, userID string, title, description *string, deadline *time.Time, status *models.TaskStatus) (*models.Task, error) {
	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}
	task := &models.Task{}
	err := DB.QueryRow(`
		UPDATE tasks SET
			title       = COALESCE($3, title),
			description = COALESCE($4, description),
			deadline    = COALESCE($5, deadline),
			status      = COALESCE($6, status),
			updated_at  = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, COALESCE(description, ''), deadline, status, created_at, updated_at`,
		id, userID, title, desc, deadline, status,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Deadline, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return task, nil
}

func DeleteTask(id, userID string) (bool, error) {
	result, err := DB.Exec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return false, fmt.Errorf("failed to delete task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get affected rows: %w", err)
	}
	return rows > 0, nil
}

func MarkOverdueTasks() (int64, error) {
	result, err := DB.Exec(`
		UPDATE tasks SET status = 'overdue', updated_at = now()
		WHERE status = 'pending' AND deadline < now()`)
	if err != nil {
		return 0, fmt.Errorf("failed to mark overdue tasks: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}
	return count, nil
}

