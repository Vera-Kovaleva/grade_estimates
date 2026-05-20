package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"grade_estimates/internal/auth"
	"grade_estimates/internal/config"
	"grade_estimates/internal/db"
	"grade_estimates/internal/estimate"
	"grade_estimates/internal/middleware"
	"grade_estimates/internal/semesters"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	log.Println("подключение к базе данных установлено")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		// Публичные маршруты (с защитой от перебора: 10 попыток в минуту)
		authHandler := auth.NewHandler(pool, cfg.JWTSecret)
		authRoutes := api.Group("/auth")
		authRoutes.Use(middleware.RateLimit(10, time.Minute))
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)

		// Защищённые маршруты
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			estimateHandler := estimate.NewHandler(pool)
			protected.POST("/estimate", estimateHandler.Estimate)
			protected.GET("/predictions", estimateHandler.GetHistory)
			protected.DELETE("/predictions", estimateHandler.ClearHistory)

			semestersHandler := semesters.NewHandler(pool)
			protected.GET("/semesters", semestersHandler.List)
			protected.POST("/semesters", semestersHandler.Add)
			protected.DELETE("/semesters/:id", semestersHandler.Delete)
		}
	}

	log.Printf("сервер запущен на порту %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
