package v1

import (
	"net/http"
	"os"
	"time"

	_ "github.com/almostinf/order_delivery_service/docs"
	"github.com/almostinf/order_delivery_service/internal/usecase"
	"github.com/almostinf/order_delivery_service/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func NewRateLimiterMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// NewRouter -.
// Swagger spec:
// @title       Order Delivery Service API
// @description Order and Courier services
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, c usecase.CourierUseCase, o usecase.OrderUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	isRateLimit, ok := os.LookupEnv("RATE_LIMITER")
	if !ok {
		l.Fatal("RATE_LIMITER variable is not set")
	}

	if isRateLimit == "ENABLE" { // I disable rate limiter for integration tests
		handler.Use(NewRateLimiterMiddleware(rate.Every(time.Second), 10))
	}

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	h := handler.Group("/v1")
	{
		newCourierRoutes(h, c, l)
		newOrderRoutes(h, o, l)
	}
}
