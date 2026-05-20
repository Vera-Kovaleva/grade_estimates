package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *service
}

func NewHandler(db *pgxpool.Pool, jwtSecret string) *Handler {
	return &Handler{service: newService(db, jwtSecret)}
}

type credentialsRequest struct {
	Login    string `json:"login"    binding:"required,max=100"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "укажите логин и пароль"})
		return
	}

	if err := h.service.register(c.Request.Context(), req.Login, req.Password); err != nil {
		if errors.Is(err, ErrLoginTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "логин уже занят"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "пользователь зарегистрирован"})
}

// Login godoc
// POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "укажите логин и пароль"})
		return
	}

	token, err := h.service.login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
