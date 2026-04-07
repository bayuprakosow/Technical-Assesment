package repository

import (
	"context"
	"database/sql"
	"time"
)

type Finding struct {
	ID        int64
	UserID    int64
	Title     string
	Severity  string
	Status    string
	CreatedAt time.Time
}

type FindingRepository struct {
	db *sql.DB
}

func NewFindingRepository(db *sql.DB) *FindingRepository {
	return &FindingRepository{db: db}
}

func (r *FindingRepository) ListByUser(ctx context.Context, userID int64) ([]Finding, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, severity, status, created_at FROM findings WHERE user_id = $1 ORDER BY id DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Finding
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ID, &f.UserID, &f.Title, &f.Severity, &f.Status, &f.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *FindingRepository) Create(ctx context.Context, userID int64, title, severity, status string) (*Finding, error) {
	if status == "" {
		status = "open"
	}
	var f Finding
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO findings (user_id, title, severity, status) VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, severity, status, created_at`,
		userID, title, severity, status,
	).Scan(&f.ID, &f.UserID, &f.Title, &f.Severity, &f.Status, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
