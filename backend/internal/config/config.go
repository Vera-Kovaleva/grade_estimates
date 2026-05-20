package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения из переменных окружения
type Config struct {
	ServerPort string
	DBUrl      string
	JWTSecret  string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string
}

// Load читает переменные окружения и возвращает конфигурацию
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, используются переменные окружения системы")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DBUrl:      getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/grade_estimates?sslmode=disable"),
		JWTSecret:  getEnv("JWT_SECRET", "change_me_in_production"),
		SMTPHost:   getEnv("SMTP_HOST", "sandbox.smtp.mailtrap.io"),
		SMTPPort:   getEnv("SMTP_PORT", "587"),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASS", ""),
		SMTPFrom:   getEnv("SMTP_FROM", "noreply@grade-estimates.local"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
