package semesters

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// errNotFound возвращается когда запись не найдена или не принадлежит пользователю
var errNotFound = errors.New("запись не найдена")

type service struct {
	repo *repository
}

func newService(db *pgxpool.Pool) *service {
	return &service{repo: newRepository(db)}
}

// add создаёт новую запись семестра
func (s *service) add(ctx context.Context, userID int, subject string, grade int, params map[string]float64) (int, error) {
	return s.repo.save(ctx, userID, subject, grade, params)
}

// list возвращает все записи семестров пользователя
func (s *service) list(ctx context.Context, userID int) ([]Semester, error) {
	return s.repo.findAll(ctx, userID)
}

// remove удаляет запись по id
func (s *service) remove(ctx context.Context, userID, id int) error {
	return s.repo.delete(ctx, userID, id)
}
