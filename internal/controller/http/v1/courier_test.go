package v1_test

import (
	"errors"
	"testing"

	v1 "github.com/almostinf/order_delivery_service/internal/controller/http/v1"
	"github.com/stretchr/testify/require"
)

func TestValidateCourier(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		in          v1.CreateCourierRequest
		expectedErr error
	}{
		{
			name: "success",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"10:00-20:00"},
			},
			expectedErr: nil,
		},
		{
			name: "wrong courier type",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"10:00-20:00"},
			},
			expectedErr: errors.New("invalid courier type format"),
		},
		{
			name: "wrong courier region",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, -2, 3},
				WorkingHours: []string{"10:00-20:00"},
			},
			expectedErr: errors.New("invalid courier region format"),
		},
		{
			name: "empty courier regions",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{},
				WorkingHours: []string{"10:00-20:00"},
			},
			expectedErr: errors.New("invalid courier regions format (regions are empty)"),
		},
		{
			name: "invalid time range format",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"10:00--20:00"},
			},
			expectedErr: errors.New("invalid working time: invalid time range format"),
		},
		{
			name: "invalid start time format",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"25:00-20:00"},
			},
			expectedErr: errors.New("invalid working time: invalid start time format"),
		},
		{
			name: "invalid end time format",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"20:00-24:00"},
			},
			expectedErr: errors.New("invalid working time: invalid end time format"),
		},
		{
			name: "the end time must not be less or equal than the start time",
			in: v1.CreateCourierRequest{
				CourierType:  "FOOT",
				Regions:      []int{1, 2, 3},
				WorkingHours: []string{"23:00-22:00"},
			},
			expectedErr: errors.New("invalid working time: the end time must not be less or equal than the start time"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v1.ValidateCourierRequest(tc.in)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.ErrorIs(t, tc.expectedErr, err)
			}
		})
	}
}
