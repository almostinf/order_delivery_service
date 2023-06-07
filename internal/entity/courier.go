package entity

import (
	"time"

	"github.com/google/uuid"
)

type Courier struct {
	CourierResponse
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CourierResponse struct {
	CourierID    uuid.UUID `json:"courier_id"`
	CourierType  string    `json:"courier_type"`
	Regions      []int     `json:"regions"`
	WorkingHours []string  `json:"working_hours"`
}

type CourierMetaInfo struct {
	CourierResponse
	Rating   int `json:"rating"`
	Earnings int `json:"earnings"`
}
