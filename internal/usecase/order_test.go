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

func order(t *testing.T) (*usecase.OrderUseCase, *mocks.MockOrder) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockOrder(mockCtrl)
	order := usecase.NewOrderUseCase(repo)

	return order, repo
}

func TestGetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	ctx := context.Background()
	var id uuid.UUID
	orderResponse := &entity.OrderResponse{}

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   *entity.OrderResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Get(ctx, id).Return(orderResponse, nil).Times(1)
			},
			res:   orderResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mock: func(repo *mocks.MockOrder) {
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

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.Get(tc.args.ctx, tc.args.id)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}

	ctx := context.Background()
	limit, offset := 10, 10
	repoErr := errors.New("some error")

	ordersResponse := []*entity.OrderResponse(nil)

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   []*entity.OrderResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				limit:  limit,
				offset: offset,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().GetAll(ctx, limit, offset).Return(ordersResponse, nil).Times(1)
			},
			res:   ordersResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:    ctx,
				limit:  limit,
				offset: offset,
			},
			mock: func(repo *mocks.MockOrder) {
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

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.GetAll(tc.args.ctx, tc.args.limit, tc.args.offset)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		orderReq *entity.Order
	}

	ctx := context.Background()
	orderResponse := &entity.OrderResponse{}
	orderReq := &entity.Order{}

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   *entity.OrderResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      ctx,
				orderReq: orderReq,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Create(ctx, orderReq).Return(orderResponse, nil).Times(1)
			},
			res:   orderResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:      ctx,
				orderReq: orderReq,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Create(ctx, orderReq).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.Create(tc.args.ctx, tc.args.orderReq)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCompleteOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx             context.Context
		completeInfoReq []entity.CompleteInfo
	}

	ctx := context.Background()
	var completeInfoReq []entity.CompleteInfo
	var orderResponse []*entity.OrderResponse

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   []*entity.OrderResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:             ctx,
				completeInfoReq: completeInfoReq,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Complete(ctx, completeInfoReq).Return(orderResponse, nil).Times(1)
			},
			res:   orderResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:             ctx,
				completeInfoReq: completeInfoReq,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Complete(ctx, completeInfoReq).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.Complete(tc.args.ctx, tc.args.completeInfoReq)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSetCourierID(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		orderID   uuid.UUID
		courierID uuid.UUID
	}

	ctx := context.Background()
	var orderID, courierID uuid.UUID
	orderResponse := &entity.OrderResponse{}

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   *entity.OrderResponse
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:       ctx,
				orderID:   orderID,
				courierID: courierID,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().SetCourierID(ctx, orderID, courierID).Return(orderResponse, nil).Times(1)
			},
			res:   orderResponse,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:       ctx,
				orderID:   orderID,
				courierID: courierID,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().SetCourierID(ctx, orderID, courierID).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.SetCourierID(tc.args.ctx, tc.args.orderID, tc.args.courierID)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAssign(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		date time.Time
	}

	ctx := context.Background()
	date := time.Time{}
	var courierAssignments []*entity.CourierAssignment

	repoErr := errors.New("some error")

	testcases := []struct {
		name  string
		args  args
		mock  func(repo *mocks.MockOrder)
		res   []*entity.CourierAssignment
		isErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:  ctx,
				date: date,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Assign(ctx, date).Return(courierAssignments, nil).Times(1)
			},
			res:   courierAssignments,
			isErr: false,
		},
		{
			name: "repo error",
			args: args{
				ctx:  ctx,
				date: date,
			},
			mock: func(repo *mocks.MockOrder) {
				repo.EXPECT().Assign(ctx, date).Return(nil, repoErr).Times(1)
			},
			res:   nil,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			order, repo := order(t)

			tc.mock(repo)

			res, err := order.Assign(tc.args.ctx, tc.args.date)

			require.Equal(t, res, tc.res)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
