package v1_test

import (
	"errors"
	"testing"

	v1 "github.com/almostinf/order_delivery_service/internal/controller/http/v1"
	"github.com/stretchr/testify/require"
)

func TestValidateOrder(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		in          v1.CreateOrderRequest
		expectedErr error
	}{
		{
			name: "success",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"10:00-20:00"},
				Cost:          100,
			},
			expectedErr: nil,
		},
		{
			name: "wrong weight",
			in: v1.CreateOrderRequest{
				Weight:        -1,
				Regions:       59,
				DeliveryHours: []string{"10:00-20:00"},
				Cost:          100,
			},
			expectedErr: errors.New("the order weight can't be less than zero"),
		},
		{
			name: "wrong regions",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       -59,
				DeliveryHours: []string{"10:00-20:00"},
				Cost:          100,
			},
			expectedErr: errors.New("the order region can't be less than zero"),
		},
		{
			name: "wrong delivery hours format",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"10:00---20:00"},
				Cost:          100,
			},
			expectedErr: errors.New("invalid delivery hours: invalid time range format"),
		},
		{
			name: "invalid start time format",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"30:00-20:00"},
				Cost:          100,
			},
			expectedErr: errors.New("invalid delivery hours: invalid start time format"),
		},
		{
			name: "invalid end time format",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"20:00-40:00"},
				Cost:          100,
			},
			expectedErr: errors.New("invalid delivery hours: invalid end time format"),
		},
		{
			name: "the end time must not be less or equal than the start time",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"15:00-12:00"},
				Cost:          100,
			},
			expectedErr: errors.New("invalid delivery hours: the end time must not be less or equal than the start time"),
		},
		{
			name: "wrong cost",
			in: v1.CreateOrderRequest{
				Weight:        25,
				Regions:       59,
				DeliveryHours: []string{"10:00-20:00"},
				Cost:          -100,
			},
			expectedErr: errors.New("the order cost can't be less than zero"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v1.ValidateOrderRequest(tc.in)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.ErrorIs(t, tc.expectedErr, err)
			}
		})
	}
}
