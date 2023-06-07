package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/internal/infrastructure/interfaces"
	"github.com/google/uuid"
)

type CourierUseCase struct {
	repo interfaces.Courier
}

func NewCourierUseCase(r interfaces.Courier) *CourierUseCase {
	return &CourierUseCase{r}
}

func (uc *CourierUseCase) Get(ctx context.Context, id uuid.UUID) (*entity.CourierResponse, error) {
	courier, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CourierUseCase - Get - uc.repo.Get: %w", err)
	}

	return courier, nil
}

func (uc *CourierUseCase) GetAll(ctx context.Context, limit, offset int) ([]*entity.CourierResponse, error) {
	couriers, err := uc.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("CourierUseCase - GetAll - uc.repo.GetAll: %w", err)
	}

	return couriers, nil
}

func (uc *CourierUseCase) Create(ctx context.Context, courier *entity.Courier) (*entity.CourierResponse, error) {
	courierRes, err := uc.repo.Create(ctx, courier)
	if err != nil {
		return nil, fmt.Errorf("CourierUseCase - Create - uc.repo.Create: %w", err)
	}

	return courierRes, nil
}

func (uc *CourierUseCase) GetMetaInfo(ctx context.Context, courierID uuid.UUID, startDate time.Time, endDate time.Time) (*entity.CourierMetaInfo, error) {
	courierMetaInfo, err := uc.repo.GetMetaInfo(ctx, courierID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("CourierUseCase - GetMetaInfo - uc.repo.GetMetaInfo: %w", err)
	}

	return courierMetaInfo, nil
}

func (uc *CourierUseCase) GetAssignments(ctx context.Context, date time.Time, courierID uuid.UUID, isAllCouriers bool) ([]*entity.CourierAssignment, error) {
	courierAssignments, err := uc.repo.GetAssignments(ctx, date, courierID, isAllCouriers)
	if err != nil {
		return nil, fmt.Errorf("CourierUseCase - GetAssignments - uc.repo.GetAssignments: %w", err)
	}

	return courierAssignments, nil
}
