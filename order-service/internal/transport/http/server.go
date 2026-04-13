package http

import (
	"log"
	"net/http"

	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	usecase *usecase.OrderUsecase
}

func NewServer(uc *usecase.OrderUsecase) *Server {
	return &Server{usecase: uc}
}

func (s *Server) Run(addr string) {
	r := gin.Default()

	r.POST("/orders", s.CreateOrder)
	r.GET("/orders/:id", s.GetOrder)
	r.PATCH("/orders/:id/cancel", s.CancelOrder)

	log.Printf("🚀 Order REST Server started on %s", addr)
	r.Run(addr)
}

func (s *Server) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id"`
		ItemName   string `json:"item_name"`
		Amount     int64  `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idempotencyKey := uuid.New().String()

	order, err := s.usecase.CreateOrder(c.Request.Context(), req.CustomerID, req.ItemName, req.Amount, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (s *Server) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := s.usecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (s *Server) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	if err := s.usecase.CancelOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
