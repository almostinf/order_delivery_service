package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/internal/infrastructure/interfaces"
	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo interfaces.Order
}

func NewOrderUseCase(r interfaces.Order) *OrderUseCase {
	return &OrderUseCase{r}
}

func (uc *OrderUseCase) Get(ctx context.Context, id uuid.UUID) (*entity.OrderResponse, error) {
	order, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Get - uc.repo.Get: %w", err)
	}

	return order, nil
}

func (uc *OrderUseCase) GetAll(ctx context.Context, limit, offset int) ([]*entity.OrderResponse, error) {
	orders, err := uc.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - GetAll - uc.repo.GetAll: %w", err)
	}

	return orders, nil
}

func (uc *OrderUseCase) Create(ctx context.Context, order *entity.Order) (*entity.OrderResponse, error) {
	orderRes, err := uc.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Create - uc.repo.Create: %w", err)
	}

	return orderRes, nil
}

func (uc *OrderUseCase) Complete(ctx context.Context, completeInfoReq []entity.CompleteInfo) ([]*entity.OrderResponse, error) {
	orderRes, err := uc.repo.Complete(ctx, completeInfoReq)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Complete - uc.repo.Complete: %w", err)
	}

	return orderRes, nil
}

func (uc *OrderUseCase) SetCourierID(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) (*entity.OrderResponse, error) {
	orderRes, err := uc.repo.SetCourierID(ctx, orderID, courierID)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - SetCourierID - uc.repo.SetCourierID: %w", err)
	}

	return orderRes, nil
}

func (uc *OrderUseCase) Assign(ctx context.Context, date time.Time) ([]*entity.CourierAssignment, error) {
	courierAssignments, err := uc.repo.Assign(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Assign - uc.repo.Assign: %w", err)
	}

	return courierAssignments, nil
}
