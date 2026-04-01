package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"payment-service/internal/repository"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

func NewApp(db *sql.DB) *gin.Engine {
	repo := repository.NewPostgresPaymentRepository(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := http.NewHandler(uc)

	r := gin.Default()
	r.POST("/payments", handler.CreatePayment)
	r.GET("/payments/:order_id", handler.GetPaymentByOrderID)
	return r
}
