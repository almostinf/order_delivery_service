package interfaces

import (
	"context"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/google/uuid"
)

type Courier interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.CourierResponse, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.CourierResponse, error)
	Create(ctx context.Context, courier *entity.Courier) (*entity.CourierResponse, error)
	GetMetaInfo(ctx context.Context, courierID uuid.UUID, startDate time.Time, endDate time.Time) (*entity.CourierMetaInfo, error)
	GetAssignments(ctx context.Context, date time.Time, courierID uuid.UUID, isAllCouriers bool) ([]*entity.CourierAssignment, error)
}
