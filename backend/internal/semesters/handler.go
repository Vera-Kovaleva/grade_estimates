package semesters

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *service
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{service: newService(db)}
}

type addSemesterRequest struct {
	Subject     string             `json:"subject" binding:"required,min=1,max=255"`
	ActualGrade int                `json:"actual_grade" binding:"required,min=2,max=5"`
	Parameters  map[string]float64 `json:"parameters" binding:"required"`
}

// List godoc
// GET /api/v1/semesters
func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt("user_id")

	semesters, err := h.service.list(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить записи"})
		return
	}

	if semesters == nil {
		semesters = []Semester{}
	}

	c.JSON(http.StatusOK, gin.H{"semesters": semesters})
}

// Add godoc
// POST /api/v1/semesters
func (h *Handler) Add(c *gin.Context) {
	var req addSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "укажите предмет, оценку (2–5) и параметры"})
		return
	}

	userID := c.GetInt("user_id")

	id, err := h.service.add(c.Request.Context(), userID, req.Subject, req.ActualGrade, req.Parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить запись"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Delete godoc
// DELETE /api/v1/semesters/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	userID := c.GetInt("user_id")

	if err := h.service.remove(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "запись не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить запись"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "запись удалена"})
}
