package app

import (
	"fmt"

	"github.com/almostinf/order_delivery_service/config"
	v1 "github.com/almostinf/order_delivery_service/internal/controller/http/v1"
	"github.com/almostinf/order_delivery_service/internal/infrastructure/repository"
	"github.com/almostinf/order_delivery_service/internal/usecase"
	"github.com/almostinf/order_delivery_service/pkg/logger"
	"github.com/almostinf/order_delivery_service/pkg/postgres"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	courierRepo := repository.NewCourierRepo(pg)
	orderRepo := repository.NewOrderRepo(pg)

	courierUseCase := usecase.NewCourierUseCase(courierRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	handler := gin.New()
	v1.NewRouter(handler, l, *courierUseCase, *orderUseCase)
	if err := handler.Run(":8080"); err != nil {
		l.Fatal(fmt.Errorf("app - Run - handler.New: %w", err))
	}
}
