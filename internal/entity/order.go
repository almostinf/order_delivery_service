package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderResponse
	DistributionDate time.Time `json:"distribution_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type OrderResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	CourierID     uuid.UUID `json:"courier_id"`
	Weight        float32   `json:"weight"`
	Regions       int       `json:"regions"`
	DeliveryHours []string  `json:"delivery_hours"`
	Cost          int       `json:"cost"`
	CompletedTime time.Time `json:"completed_time"`
}

type CompleteInfo struct {
	CourierID    uuid.UUID `json:"courier_id" binding:"required"`
	OrderID      uuid.UUID `json:"order_id" binding:"required"`
	CompleteTime time.Time `json:"complete_time" binding:"required"`
}

type OrdersGroup struct {
	GroupOrderID uuid.UUID       `json:"group_order_id" binding:"required"`
	Orders       []OrderResponse `json:"orders" binding:"required"`
}

type CourierAssignment struct {
	CourierID uuid.UUID     `json:"courier_id" binding:"required"`
	Orders    []OrdersGroup `json:"orders" binding:"required"`
}
