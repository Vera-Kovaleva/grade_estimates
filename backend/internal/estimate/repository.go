package estimate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Prediction представляет одну запись из истории расчётов
type Prediction struct {
	ID             int                `json:"id"`
	Parameters     map[string]float64 `json:"parameters"`
	PredictedGrade int                `json:"predicted_grade"`
	CreatedAt      time.Time          `json:"created_at"`
}

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

// save сохраняет результат расчёта в БД и возвращает id записи
func (r *repository) save(ctx context.Context, userID int, params map[string]float64, grade int) (int, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return 0, fmt.Errorf("ошибка сериализации параметров: %w", err)
	}

	var id int
	err = r.db.QueryRow(ctx,
		`INSERT INTO predictions (user_id, parameters, predicted_grade) VALUES ($1, $2, $3) RETURNING id`,
		userID, data, grade,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось сохранить расчёт: %w", err)
	}
	return id, nil
}

// findAll возвращает всю историю расчётов пользователя
func (r *repository) findAll(ctx context.Context, userID int) ([]Prediction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, parameters, predicted_grade, created_at
		 FROM predictions WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить историю: %w", err)
	}
	defer rows.Close()

	var predictions []Prediction
	for rows.Next() {
		var p Prediction
		var rawParams []byte
		if err := rows.Scan(&p.ID, &rawParams, &p.PredictedGrade, &p.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rawParams, &p.Parameters); err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}
	return predictions, nil
}

// deleteAll удаляет всю историю расчётов пользователя
func (r *repository) deleteAll(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM predictions WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("не удалось очистить историю: %w", err)
	}
	return nil
}
