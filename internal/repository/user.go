package repository

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}
