package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Postgres struct {
	db *sqlx.DB
}

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Create(ctx context.Context, email string, passwordHash string) (*User, error) {
	query := "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at"
	row := p.db.QueryRowContext(ctx, query, email, passwordHash)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	return &user, err
}
