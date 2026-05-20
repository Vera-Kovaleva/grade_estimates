package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User представляет пользователя из базы данных
type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

// create сохраняет нового пользователя в БД
func (r *repository) create(ctx context.Context, login, passwordHash string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (login, password_hash) VALUES ($1, $2)`,
		login, passwordHash,
	)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}
	return nil
}

// findByLogin возвращает пользователя по логину
func (r *repository) findByLogin(ctx context.Context, login string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, login, password_hash FROM users WHERE login = $1`,
		login,
	).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return user, nil
}
