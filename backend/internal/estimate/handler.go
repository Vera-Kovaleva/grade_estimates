package estimate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *service
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{service: newService(db)}
}

type estimateRequest struct {
	Parameters map[string]float64 `json:"parameters" binding:"required"`
}

// Estimate godoc
// POST /api/v1/estimate
func (h *Handler) Estimate(c *gin.Context) {
	var req estimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "укажите параметры расчёта"})
		return
	}

	userID := c.GetInt("user_id")

	id, grade, err := h.service.calculate(c.Request.Context(), userID, req.Parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось выполнить расчёт"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "grade": grade})
}

// GetHistory godoc
// GET /api/v1/predictions
func (h *Handler) GetHistory(c *gin.Context) {
	userID := c.GetInt("user_id")

	predictions, err := h.service.history(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить историю"})
		return
	}

	// возвращаем пустой массив вместо null если истории нет
	if predictions == nil {
		predictions = []Prediction{}
	}

	c.JSON(http.StatusOK, gin.H{"predictions": predictions})
}

// ClearHistory godoc
// DELETE /api/v1/predictions
func (h *Handler) ClearHistory(c *gin.Context) {
	userID := c.GetInt("user_id")

	if err := h.service.clearHistory(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось очистить историю"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "история очищена"})
}
