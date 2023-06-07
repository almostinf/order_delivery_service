package interfaces

import (
	"context"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/google/uuid"
)

type Order interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.OrderResponse, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.OrderResponse, error)
	Create(ctx context.Context, courier *entity.Order) (*entity.OrderResponse, error)
	Complete(ctx context.Context, completeInfoReq []entity.CompleteInfo) ([]*entity.OrderResponse, error)
	SetCourierID(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) (*entity.OrderResponse, error)
	Assign(ctx context.Context, date time.Time) ([]*entity.CourierAssignment, error)
}
