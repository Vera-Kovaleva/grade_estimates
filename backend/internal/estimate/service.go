package estimate

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	repo *repository
}

func newService(db *pgxpool.Pool) *service {
	return &service{repo: newRepository(db)}
}

// calculate вычисляет оценку по входным параметрам на основе формулы GPA.
// GPA = 2.0 + 0.12S − 0.18A + 0.10T + 0.10P + 0.05E + 0.04Sp + 0.03M + 0.03V + 0.07Ed + 0.08H
// Где:
//
//	S  — учебное время (study_hours):              0–4 балла
//	A  — прогулы (absences):                       0–4 балла
//	T  — репетитор (tutor):                        0 или 2 балла
//	P  — поддержка родителей (parental):           0–4 балла
//	E  — внеклассная деятельность (extracurricular): 0 или 1 балл
//	Sp — спорт (sports):                           0 или 1 балл
//	M  — музыка (music):                           0 или 1 балл
//	V  — волонтерство (volunteering):              0 или 1 балл
//	Ed — образование родителей (parent_education): 0–4 балла
//	H  — выполнение д/з (homework):                0–4 балла
func (s *service) calculate(ctx context.Context, userID int, params map[string]float64) (int, int, error) {
	gpa := computeGPA(params)
	grade := gpaToGrade(gpa)

	id, err := s.repo.save(ctx, userID, params, grade)
	if err != nil {
		return 0, 0, err
	}
	return id, grade, nil
}

// history возвращает историю расчётов пользователя
func (s *service) history(ctx context.Context, userID int) ([]Prediction, error) {
	return s.repo.findAll(ctx, userID)
}

// clearHistory удаляет всю историю расчётов пользователя
func (s *service) clearHistory(ctx context.Context, userID int) error {
	return s.repo.deleteAll(ctx, userID)
}

// computeGPA вычисляет GPA по формуле из технического задания.
func computeGPA(params map[string]float64) float64 {
	S := params["study_hours"]
	A := params["absences"]
	T := params["tutor"]
	P := params["parental"]
	E := params["extracurricular"]
	Sp := params["sports"]
	M := params["music"]
	V := params["volunteering"]
	Ed := params["parent_education"]
	H := params["homework"]

	return 2.0 + 0.12*S - 0.18*A + 0.10*T + 0.10*P + 0.05*E + 0.04*Sp + 0.03*M + 0.03*V + 0.07*Ed + 0.08*H
}

// gpaToGrade переводит значение GPA в школьную оценку (2–5).
func gpaToGrade(gpa float64) int {
	switch {
	case gpa >= 3.5:
		return 5
	case gpa >= 3.0:
		return 4
	case gpa >= 2.5:
		return 3
	default:
		return 2
	}
}
