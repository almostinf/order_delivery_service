package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	mocks "github.com/almostinf/order_delivery_service/internal/mocks/repo"
	"github.com/almostinf/order_delivery_service/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func courier(t *testing.T) (*usecase.CourierUseCase, *mocks.MockCourier) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockCourier(mockCtrl)
	courier := usecase.NewCourierUseCase(repo)

	return courier, repo
}

func TestGetCourier(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	ctx := context.Background()
	var id uuid.UUID
	repoErr := errors.New("some error")

	courierResponse := &entity.CourierResponse{}

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockCourier)
		res   *entity.CourierResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().Get(ctx, id).Return(courierResponse, nil).Times(1)
			},
			res:   courierResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().Get(ctx, id).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courier, repo := courier(t)

			tc.mock(repo)

			res, err := courier.Get(tc.args.ctx, tc.args.id)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAllCouriers(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}

	ctx := context.Background()
	limit, offset := 10, 10
	repoErr := errors.New("some error")

	couriersResponse := []*entity.CourierResponse(nil)

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockCourier)
		res   []*entity.CourierResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				limit:  limit,
				offset: offset,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetAll(ctx, limit, offset).Return(couriersResponse, nil).Times(1)
			},
			res:   couriersResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:    ctx,
				limit:  limit,
				offset: offset,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetAll(ctx, limit, offset).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courier, repo := courier(t)

			tc.mock(repo)

			res, err := courier.GetAll(tc.args.ctx, tc.args.limit, tc.args.offset)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateCourier(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		courierReq *entity.Courier
	}

	ctx := context.Background()
	courierResponse := &entity.CourierResponse{}
	courierReq := &entity.Courier{}

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockCourier)
		res   *entity.CourierResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:        ctx,
				courierReq: courierReq,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().Create(ctx, courierReq).Return(courierResponse, nil).Times(1)
			},
			res:   courierResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:        ctx,
				courierReq: courierReq,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().Create(ctx, courierReq).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courier, repo := courier(t)

			tc.mock(repo)

			res, err := courier.Create(tc.args.ctx, tc.args.courierReq)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetMetaInfoFromCourier(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		courierID uuid.UUID
		startDate time.Time
		endDate   time.Time
	}

	ctx := context.Background()
	var courierID uuid.UUID
	startDate := time.Now()
	endDate := time.Now()
	courierMetaResponse := &entity.CourierMetaInfo{}

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockCourier)
		res   *entity.CourierMetaInfo
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:       ctx,
				courierID: courierID,
				startDate: startDate,
				endDate:   endDate,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetMetaInfo(ctx, courierID, startDate, endDate).Return(courierMetaResponse, nil).Times(1)
			},
			res:   courierMetaResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:       ctx,
				courierID: courierID,
				startDate: startDate,
				endDate:   endDate,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetMetaInfo(ctx, courierID, startDate, endDate).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courier, repo := courier(t)

			tc.mock(repo)

			res, err := courier.GetMetaInfo(tc.args.ctx, tc.args.courierID, tc.args.startDate, tc.args.endDate)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAssignments(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		date          time.Time
		courierID     uuid.UUID
		isAllCouriers bool
	}

	ctx := context.Background()
	date := time.Time{}
	var courierID uuid.UUID
	isAllCouriers := false

	var courierAssignments []*entity.CourierAssignment

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockCourier)
		res   []*entity.CourierAssignment
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:           ctx,
				date:          date,
				courierID:     courierID,
				isAllCouriers: isAllCouriers,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetAssignments(ctx, date, courierID, isAllCouriers).Return(courierAssignments, nil).Times(1)
			},
			res:   courierAssignments,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:           ctx,
				date:          date,
				courierID:     courierID,
				isAllCouriers: isAllCouriers,
			},
			mock: func(repo *mocks.MockCourier) {
				repo.EXPECT().GetAssignments(ctx, date, courierID, isAllCouriers).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courier, repo := courier(t)

			tc.mock(repo)

			res, err := courier.GetAssignments(tc.args.ctx, tc.args.date, tc.args.courierID, tc.args.isAllCouriers)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
