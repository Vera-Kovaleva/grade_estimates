package semesters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Semester представляет одну запись из журнала прошлых семестров
type Semester struct {
	ID          int                `json:"id"`
	Subject     string             `json:"subject"`
	ActualGrade int                `json:"actual_grade"`
	Parameters  map[string]float64 `json:"parameters"`
	CreatedAt   time.Time          `json:"created_at"`
}

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

// save сохраняет новую запись семестра и возвращает id
func (r *repository) save(ctx context.Context, userID int, subject string, grade int, params map[string]float64) (int, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return 0, fmt.Errorf("ошибка сериализации параметров: %w", err)
	}

	var id int
	err = r.db.QueryRow(ctx,
		`INSERT INTO semesters (user_id, subject, actual_grade, parameters) VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, subject, grade, data,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось сохранить запись: %w", err)
	}
	return id, nil
}

// findAll возвращает все записи семестров пользователя
func (r *repository) findAll(ctx context.Context, userID int) ([]Semester, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, subject, actual_grade, parameters, created_at
		 FROM semesters WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить записи: %w", err)
	}
	defer rows.Close()

	var semesters []Semester
	for rows.Next() {
		var s Semester
		var rawParams []byte
		if err := rows.Scan(&s.ID, &s.Subject, &s.ActualGrade, &rawParams, &s.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rawParams, &s.Parameters); err != nil {
			return nil, err
		}
		semesters = append(semesters, s)
	}
	return semesters, nil
}

// delete удаляет запись по id, проверяя что она принадлежит пользователю
func (r *repository) delete(ctx context.Context, userID, id int) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM semesters WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("не удалось удалить запись: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errNotFound
	}
	return nil
}
