package store

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"debitask/models"
)

func CreateUser(email, passwordHash string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id::text, email, password_hash, created_at, updated_at",
		email, passwordHash,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(
		"SELECT id::text, email, password_hash, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(
		"SELECT id::text, email, password_hash, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
