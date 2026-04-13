package http

import (
	"errors"
	"net/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.OrderUsecase
}

func NewHandler(uc *usecase.OrderUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		ItemName   string `json:"item_name" binding:"required"`
		Amount     int64  `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idempKey := c.GetHeader("X-Idempotency-Key")
	order, err := h.usecase.CreateOrder(c.Request.Context(), req.CustomerID, req.ItemName, req.Amount, idempKey)

	if errors.Is(err, usecase.ErrPaymentUnavailable) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Payment service is unavailable"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.usecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	err := h.usecase.CancelOrder(c.Request.Context(), id)
	if errors.Is(err, usecase.ErrOrderNotFound) || errors.Is(err, usecase.ErrCannotCancel) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
