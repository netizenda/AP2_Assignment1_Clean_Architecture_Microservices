package app

import (
	"database/sql"
	"order-service/internal/client"
	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func NewApp(db *sql.DB, paymentURL string) *gin.Engine {
	repo := repository.NewPostgresOrderRepository(db)
	paymentClient := client.NewPaymentClient(paymentURL)
	uc := usecase.NewOrderUsecase(repo, paymentClient)
	handler := http.NewHandler(uc)

	r := gin.Default()
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/:id", handler.GetOrder)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)
	return r
}
