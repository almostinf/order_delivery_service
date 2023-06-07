package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/pkg/postgres"
	"github.com/google/uuid"
)

type CourierRepo struct {
	*postgres.Postgres
}

func NewCourierRepo(pg *postgres.Postgres) *CourierRepo {
	return &CourierRepo{pg}
}

var _getAllSchema = `
	SELECT courier_id, courier_type, regions, working_hours
	FROM couriers;
`

func (r *CourierRepo) GetAll(ctx context.Context, limit, offset int) ([]*entity.CourierResponse, error) {
	var couriers []*entity.CourierResponse

	rows, err := r.Pool.Query(ctx, _getAllSchema)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - GetAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		e := &entity.CourierResponse{}

		err = rows.Scan(&e.CourierID, &e.CourierType, &e.Regions, &e.WorkingHours)
		if err != nil {
			return nil, fmt.Errorf("CourierRepo - GetAll - rows.Scan: %w", err)
		}

		couriers = append(couriers, e)
	}

	return couriers, nil
}

var _getSchema = `
	SELECT courier_id, courier_type, regions, working_hours
	FROM couriers
	WHERE courier_id = $1;
`

func (r *CourierRepo) Get(ctx context.Context, id uuid.UUID) (*entity.CourierResponse, error) {
	var courier entity.CourierResponse

	rows, err := r.Pool.Query(ctx, _getSchema, id)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&courier.CourierID, &courier.CourierType, &courier.Regions, &courier.WorkingHours)
		if err != nil {
			return nil, fmt.Errorf("CourierRepo - Get - rows.Scan: %w", err)
		}
	} else {
		return nil, fmt.Errorf("CourierRepo - Get - no rows found")
	}

	return &courier, nil
}

var _createSchema = `
	INSERT INTO couriers (courier_id, courier_type, regions, working_hours, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING courier_id;
`

func (r *CourierRepo) Create(ctx context.Context, courier *entity.Courier) (*entity.CourierResponse, error) {
	courierRes := &entity.CourierResponse{
		CourierID:    courier.CourierID,
		CourierType:  courier.CourierType,
		Regions:      courier.Regions,
		WorkingHours: courier.WorkingHours,
	}

	err := r.Pool.QueryRow(ctx, _createSchema, courier.CourierID, courier.CourierType, courier.Regions, courier.WorkingHours, courier.CreatedAt, courier.UpdatedAt).Scan(&courierRes.CourierID)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - Create - r.Pool.QueryRow: %w", err)
	}

	return courierRes, nil
}

var _checkIfCourierExists = `
	SELECT EXISTS(SELECT 1 FROM couriers WHERE courier_id = $1)
`

var _getTotalCostOfOrdersBetweenDuration = `
	SELECT cost FROM orders
	WHERE courier_id = $1
	AND completed_time BETWEEN $2 AND $3
`

func getTotalCost(costs []int) int {
	totalCost := 0
	for _, cost := range costs {
		totalCost += cost
	}
	return totalCost
}

func (r *CourierRepo) GetMetaInfo(ctx context.Context, courierID uuid.UUID, startDate time.Time, endDate time.Time) (*entity.CourierMetaInfo, error) {
	courier, err := r.Get(ctx, courierID)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - GetMetaInfo - r.Get: %w", err)
	}

	courierMetaInfo := entity.CourierMetaInfo{
		CourierResponse: *courier,
		Rating:          0,
		Earnings:        0,
	}

	rows, err := r.Pool.Query(ctx, _getTotalCostOfOrdersBetweenDuration, courierID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - GetMetaInfo - r.Pool.Query(_getTotalCostOfOrdersBetweenDuration): %w", err)
	}
	defer rows.Close()

	costs := make([]int, 0)
	for rows.Next() {
		var cost int
		err = rows.Scan(&cost)
		if err != nil {
			return nil, fmt.Errorf("CourierRepo - GetMetaInfo - rows.Scan: %w", err)
		}

		costs = append(costs, cost)
	}

	var coefCost, coefRating int
	switch courier.CourierType {
	case "FOOT":
		coefCost = 2
		coefRating = 3
	case "BIKE":
		coefCost = 3
		coefRating = 2
	case "AUTO":
		coefCost = 4
		coefRating = 1
	}

	courierMetaInfo.Earnings = coefCost * getTotalCost(costs)

	dur := int(endDate.Sub(startDate).Hours())
	log.Println(dur, coefRating, len(costs))

	courierMetaInfo.Rating = len(costs) / dur * coefRating

	return &courierMetaInfo, nil
}

var _getAllOrdersWithGivenDistributionDate = `
	SELECT order_id, weight, regions, delivery_hours, cost, completed_time
	FROM orders
	WHERE distribution_date = $1;
`

func (r *CourierRepo) getDistributedOrders(ctx context.Context, date time.Time) ([]*entity.OrderResponse, error) {
	var orders []*entity.OrderResponse

	rows, err := r.Pool.Query(ctx, _getAllOrdersWithGivenDistributionDate, date)
	if err != nil {
		return nil, fmt.Errorf("CourierRepo - getDistributedOrders - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.OrderResponse
		err := rows.Scan(&order.OrderID, &order.Weight, &order.Regions, &order.DeliveryHours, &order.Cost, &order.CompletedTime)
		if err != nil {
			return nil, fmt.Errorf("CourierRepo - getDistributedOrders - rows.Scan: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func newOrderGroup(assignments map[uuid.UUID]map[int][]entity.OrdersGroup, courierID uuid.UUID, region int, order *entity.OrderResponse) {
	ordersGroups := assignments[courierID][order.Regions]

	newOrderGroup := entity.OrdersGroup{
		GroupOrderID: uuid.New(),
		Orders:       make([]entity.OrderResponse, 0),
	}

	newOrderGroup.Orders = append(newOrderGroup.Orders, *order)
	ordersGroups = append(ordersGroups, newOrderGroup)

	assignments[courierID][order.Regions] = ordersGroups
}

func (r *CourierRepo) setOrderInAssignments(ctx context.Context, assignments map[uuid.UUID]map[int][]entity.OrdersGroup, courierID uuid.UUID, order *entity.OrderResponse) error {
	courier, err := r.Get(ctx, courierID)
	if err != nil {
		return err
	}

	var maxCount int
	switch courier.CourierType {
	case "FOOT":
		maxCount = 2
	case "BIKE":
		maxCount = 4
	case "AUTO":
		maxCount = 7
	}

	if reg, ok := assignments[courierID]; ok {
		for region, ordersGroup := range reg {
			if region == order.Regions {
				for _, orderGroup := range ordersGroup {
					if len(orderGroup.Orders) < maxCount {
						orderGroup.Orders = append(orderGroup.Orders, *order)
						assignments[courierID][order.Regions] = ordersGroup
						return nil
					}
				}

				newOrderGroup(assignments, courierID, order.Regions, order)
			}
		}
	} else {
		assignments[courierID] = make(map[int][]entity.OrdersGroup)
		assignments[courierID][order.Regions] = make([]entity.OrdersGroup, 0)

		newOrderGroup(assignments, courierID, order.Regions, order)
	}

	return nil
}

func (r *CourierRepo) GetAssignments(ctx context.Context, date time.Time, courierID uuid.UUID, isAllCouriers bool) ([]*entity.CourierAssignment, error) {
	orderRepo := NewOrderRepo(r.Postgres)
	couriersAssignment := make([]*entity.CourierAssignment, 0)
	assignments := make(map[uuid.UUID]map[int][]entity.OrdersGroup) // mapping: courier_id -> region -> slice of orders group
	distributedOrders, err := r.getDistributedOrders(ctx, date)
	log.Println("Len of distributed orders: ", len(distributedOrders))

	if err != nil {
		return nil, err
	}

	for _, distributedOrder := range distributedOrders {
		order, err := orderRepo.Get(ctx, distributedOrder.OrderID)
		if err != nil {
			return nil, err
		}

		if isAllCouriers {
			if err := r.setOrderInAssignments(ctx, assignments, order.CourierID, distributedOrder); err != nil {
				return nil, err
			}
		} else {
			if order.CourierID == courierID {
				if err := r.setOrderInAssignments(ctx, assignments, courierID, distributedOrder); err != nil {
					return nil, err
				}
			}
		}
	}

	for courierID, regions := range assignments {
		assignment := entity.CourierAssignment{
			CourierID: courierID,
			Orders:    make([]entity.OrdersGroup, 0),
		}

		for _, orderGroup := range regions {
			assignment.Orders = append(assignment.Orders, orderGroup...)
		}

		couriersAssignment = append(couriersAssignment, &assignment)
	}

	return couriersAssignment, nil
}
