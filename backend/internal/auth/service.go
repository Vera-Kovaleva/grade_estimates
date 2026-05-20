package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginTaken         = errors.New("логин уже занят")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")
)

type service struct {
	repo      *repository
	jwtSecret string
}

func newService(db *pgxpool.Pool, jwtSecret string) *service {
	return &service{
		repo:      newRepository(db),
		jwtSecret: jwtSecret,
	}
}

// register создаёт нового пользователя
func (s *service) register(ctx context.Context, login, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хэширования пароля: %w", err)
	}

	if err := s.repo.create(ctx, login, string(hash)); err != nil {
		return ErrLoginTaken
	}
	return nil
}

// login проверяет учётные данные и возвращает JWT токен
func (s *service) login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.findByLogin(ctx, login)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("ошибка создания токена: %w", err)
	}

	return signed, nil
}
