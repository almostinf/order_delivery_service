package repository

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/pkg/postgres"
	"github.com/google/uuid"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

var _getAllOrdersSchema = `
	SELECT order_id, courier_id, weight, regions, delivery_hours, cost, completed_time
	FROM orders;
`

func (r *OrderRepo) GetAll(ctx context.Context, limit, offset int) ([]*entity.OrderResponse, error) {
	var orders []*entity.OrderResponse

	rows, err := r.Pool.Query(ctx, _getAllOrdersSchema)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		e := &entity.OrderResponse{}

		err = rows.Scan(&e.OrderID, &e.CourierID, &e.Weight, &e.Regions, &e.DeliveryHours, &e.Cost, &e.CompletedTime)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - GetAll - rows.Scan: %w", err)
		}

		orders = append(orders, e)
	}

	return orders, nil
}

var _getOrderSchema = `
	SELECT order_id, courier_id, weight, regions, delivery_hours, cost, completed_time
	FROM orders
	WHERE order_id = $1;
`

func (r *OrderRepo) Get(ctx context.Context, id uuid.UUID) (*entity.OrderResponse, error) {
	var order entity.OrderResponse

	rows, err := r.Pool.Query(ctx, _getOrderSchema, id)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&order.OrderID, &order.CourierID, &order.Weight, &order.Regions, &order.DeliveryHours, &order.Cost, &order.CompletedTime)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - Get - rows.Scan: %w", err)
		}
	} else {
		return nil, fmt.Errorf("OrderRepo - Get - no rows found")
	}

	return &order, nil
}

var _getFullOrder = `
	SELECT order_id, courier_id, weight, regions, delivery_hours, cost, completed_time, courier_id, distribution_date, created_at, updated_at
	FROM orders
	WHERE order_id = $1;
`

func (r *OrderRepo) getFullOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order

	rows, err := r.Pool.Query(ctx, _getFullOrder, id)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - getFullOrder - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&order.OrderID, &order.CourierID, &order.Weight, &order.Regions, &order.DeliveryHours, &order.Cost, &order.CompletedTime, &order.CourierID, &order.DistributionDate, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - getFullOrder - rows.Scan: %w", err)
		}
	} else {
		return nil, fmt.Errorf("OrderRepo - getFullOrder - no rows found")
	}

	return &order, nil
}

var _createOrderSchema = `
	INSERT INTO orders (order_id, weight, regions, delivery_hours, cost, completed_time, distribution_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING order_id;
`

func (r *OrderRepo) Create(ctx context.Context, order *entity.Order) (*entity.OrderResponse, error) {
	orderRes := &entity.OrderResponse{
		OrderID:       order.OrderID,
		Weight:        order.Weight,
		Regions:       order.Regions,
		DeliveryHours: order.DeliveryHours,
		Cost:          order.Cost,
		CompletedTime: order.CompletedTime,
	}

	err := r.Pool.QueryRow(ctx, _createOrderSchema, order.OrderID, order.Weight, order.Regions, order.DeliveryHours, order.Cost, order.CompletedTime, time.Time{}, order.CreatedAt, order.UpdatedAt).Scan(&orderRes.OrderID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Create - r.Pool.QueryRow: %w", err)
	}

	return orderRes, nil
}

var _setOrderCompletedTime = `
	UPDATE orders SET completed_time = $1 WHERE order_id = $2
`

func (r *OrderRepo) Complete(ctx context.Context, completeInfoReq []entity.CompleteInfo) ([]*entity.OrderResponse, error) {
	orders := make([]*entity.OrderResponse, 0)
	for _, completeInfo := range completeInfoReq {
		fullOrder, err := r.getFullOrder(ctx, completeInfo.OrderID)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - Complete - r.Get: %w", err)
		}

		_, err = r.Pool.Exec(ctx, _checkIfCourierExists, completeInfo.CourierID)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - Complete - r.Pool.Exec(_checkIfCourierExists): %w", err)
		}

		if fullOrder.CourierID != completeInfo.CourierID {
			return nil, fmt.Errorf("OrderRepo - Complete - courierID and courierIDResponse don't match")
		}

		if !(fullOrder.DistributionDate.Day() == completeInfo.CompleteTime.Day() && fullOrder.DistributionDate.Month() == completeInfo.CompleteTime.Month() &&
			fullOrder.DistributionDate.Year() == completeInfo.CompleteTime.Year()) {
			return nil, fmt.Errorf("OrderRepo - Complete - complete_time and distribution_date don't match")
		}

		if fullOrder.CompletedTime.IsZero() {
			fullOrder.CompletedTime = completeInfo.CompleteTime
		}

		_, err = r.Pool.Exec(ctx, _setOrderCompletedTime, completeInfo.CompleteTime, completeInfo.OrderID)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - Complete - r.Pool.Exec(_setOrderCompletedTime): %w", err)
		}

		order := &entity.OrderResponse{
			OrderID:       fullOrder.OrderID,
			CourierID:     fullOrder.CourierID,
			Weight:        fullOrder.Weight,
			Regions:       fullOrder.Regions,
			DeliveryHours: fullOrder.DeliveryHours,
			Cost:          fullOrder.Cost,
			CompletedTime: fullOrder.CompletedTime,
		}

		orders = append(orders, order)
	}

	return orders, nil
}

var _setOrderCourierID = `
	UPDATE orders SET courier_id = $1 WHERE order_id = $2
`

func (r *OrderRepo) SetCourierID(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) (*entity.OrderResponse, error) {
	order, err := r.Get(ctx, orderID)

	if err != nil {
		return nil, fmt.Errorf("OrderRepo - SetCourierID - r.Get: %w", err)
	}

	_, err = r.Pool.Exec(ctx, _checkIfCourierExists, courierID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - SetCourierID - r.Pool.Exec(_checkIfCourierExists): %w", err)
	}

	order.CourierID = courierID

	_, err = r.Pool.Exec(ctx, _setOrderCourierID, courierID, orderID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - SetCourierID - r.Pool.Exec(_setOrderCourierID): %w", err)
	}

	orderRes := &entity.OrderResponse{
		OrderID:       order.OrderID,
		CourierID:     order.CourierID,
		Weight:        order.Weight,
		Regions:       order.Regions,
		DeliveryHours: order.DeliveryHours,
		Cost:          order.Cost,
		CompletedTime: order.CompletedTime,
	}

	return orderRes, nil
}

var _getOrdersLower40kgAndNotDistributed = `
	SELECT order_id, courier_id, weight, regions, delivery_hours, cost, completed_time
	FROM orders
	WHERE weight < 40 AND distribution_date = '0001-01-01 00:00:00';
`

func (r *OrderRepo) getOrdersForAssign(ctx context.Context) ([]*entity.OrderResponse, error) {
	var orders []*entity.OrderResponse

	rows, err := r.Pool.Query(ctx, _getOrdersLower40kgAndNotDistributed)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - getOrdersForAssign - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.OrderResponse
		err := rows.Scan(&order.OrderID, &order.CourierID, &order.Weight, &order.Regions, &order.DeliveryHours, &order.Cost, &order.CompletedTime)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - getOrdersForAssign - rows.Scan: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

type timeRange struct {
	start time.Time
	end   time.Time
}

func parseTime(t1 string, t2 string) (timeRange, timeRange) {
	var tr1, tr2 timeRange

	layout := "15:04"
	parts := strings.Split(t1, "-")
	startTime, _ := time.Parse(layout, parts[0]) // We can't get an error cause we validate it in controller
	endTime, _ := time.Parse(layout, parts[1])

	tr1.start, tr1.end = startTime, endTime

	parts = strings.Split(t2, "-")
	startTime, _ = time.Parse(layout, parts[0])
	endTime, _ = time.Parse(layout, parts[1])

	tr2.start, tr2.end = startTime, endTime

	return tr1, tr2
}

func checkTimeOverlap(t1 timeRange, t2 timeRange, overlap time.Duration) bool {
	if t1.end.Before(t2.start) || t2.end.Before(t1.start) {
		return false
	}

	overlapStart := t1.start
	if t2.start.After(overlapStart) {
		overlapStart = t2.start
	}

	overlapEnd := t1.end
	if t2.end.Before(overlapEnd) {
		overlapEnd = t2.end
	}

	return overlapEnd.Sub(overlapStart) >= overlap
}

func containRegion(regions []int, target int) bool {
	for _, region := range regions {
		if region == target {
			return true
		}
	}

	return false
}

type ordersGroupWithLeftTime struct {
	orders   []entity.OrdersGroup
	leftTime int
}

func initOrdersGroup(assignments map[uuid.UUID]map[int]ordersGroupWithLeftTime, courierID uuid.UUID, order entity.OrderResponse, leftTime int) {
	assignments[courierID][order.Regions] = ordersGroupWithLeftTime{
		orders:   make([]entity.OrdersGroup, 0),
		leftTime: leftTime,
	}

	orders := assignments[courierID][order.Regions].orders
	orderG := entity.OrdersGroup{
		GroupOrderID: uuid.New(),
		Orders:       make([]entity.OrderResponse, 0),
	}

	orderG.Orders = append(orderG.Orders, order)
	orders = append(orders, orderG)

	assignments[courierID][order.Regions] = ordersGroupWithLeftTime{
		orders:   orders,
		leftTime: leftTime,
	}
}

func findAndSetFreeCourier(couriers []*entity.CourierResponse, order *entity.OrderResponse, assignments map[uuid.UUID]map[int]ordersGroupWithLeftTime) bool {
	for _, courier := range couriers {
		if containRegion(courier.Regions, order.Regions) {
			var overlap time.Duration
			var maxRegions int
			var maxCount int
			var nextDeliveryTime int

			switch courier.CourierType {
			case "FOOT":
				overlap = 25 * time.Minute
				maxRegions = 1
				maxCount = 2
				nextDeliveryTime = 10
			case "BIKE":
				overlap = 12 * time.Minute
				maxRegions = 2
				maxCount = 4
				nextDeliveryTime = 8
			case "AUTO":
				overlap = 8 * time.Minute
				maxRegions = 3
				maxCount = 7
				nextDeliveryTime = 4
			}

			for i := 0; i < len(courier.WorkingHours); i++ {
				for j := 0; j < len(order.DeliveryHours); j++ {
					courierTime, orderTime := parseTime(courier.WorkingHours[i], order.DeliveryHours[j])
					if checkTimeOverlap(courierTime, orderTime, overlap) {
						if reg, ok := assignments[courier.CourierID]; ok {
							if orderGroup, ok := reg[order.Regions]; ok {
								for _, orderG := range orderGroup.orders {
									if len(orderG.Orders) < maxCount {
										if orderGroup.leftTime > nextDeliveryTime {
											orderG.Orders = append(orderG.Orders, *order)
											assignments[courier.CourierID][order.Regions] = ordersGroupWithLeftTime{
												orders:   orderGroup.orders,
												leftTime: orderGroup.leftTime - nextDeliveryTime,
											}
											return true
										}
									}
								}

								if orderGroup.leftTime > int(overlap.Minutes()) {
									initOrdersGroup(assignments, courier.CourierID, *order, orderGroup.leftTime-int(overlap.Minutes()))
									return true
								}
								log.Println("doesn't have enough time", orderGroup.leftTime, int(overlap.Minutes()))
							} else {
								if len(reg) <= maxRegions {
									// reg[order.Regions] = ordersGroupWithLeftTime{
									// 	orders:   make([]entity.OrdersGroup, 0),
									// 	leftTime: int(courierTime.end.Sub(courierTime.start).Minutes()),
									// }

									initOrdersGroup(assignments, courier.CourierID, *order, int(courierTime.end.Sub(courierTime.start).Minutes())-int(overlap.Minutes()))
									return true
								}
								log.Println("Got max regions limit")
							}
						} else {
							assignments[courier.CourierID] = make(map[int]ordersGroupWithLeftTime)
							initOrdersGroup(assignments, courier.CourierID, *order, int(courierTime.end.Sub(courierTime.start).Minutes())-int(overlap.Minutes()))
							return true
						}
					} else {
						log.Println("Time is not overlap: ", courier, order)
					}
				}
			}
		}
	}

	return false
}

var _getCouriersWithGivenType = `
	SELECT courier_id, courier_type, regions, working_hours
	FROM couriers
	WHERE courier_type = $1;
`

func (r *OrderRepo) getCouriersWithGivenType(ctx context.Context, courierType string) ([]*entity.CourierResponse, error) {
	var couriers []*entity.CourierResponse

	rows, err := r.Pool.Query(ctx, _getCouriersWithGivenType, courierType)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - getCouriersWithGivenType - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var courier entity.CourierResponse
		err = rows.Scan(&courier.CourierID, &courier.CourierType, &courier.Regions, &courier.WorkingHours)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - getCouriersWithGivenType - rows.Scan: %w", err)
		}

		couriers = append(couriers, &courier)
	}

	return couriers, nil
}

var _updateDistributionDateAndCourierIDInOrder = `
	UPDATE orders
	SET distribution_date = $1,
		courier_id = $2
	WHERE order_id = $3;
`

func (r *OrderRepo) Assign(ctx context.Context, date time.Time) ([]*entity.CourierAssignment, error) {
	couriersAssignment := make([]*entity.CourierAssignment, 0)
	assignments := make(map[uuid.UUID]map[int]ordersGroupWithLeftTime) // mapping: courier_id -> region -> ordersGroupWithLeftTime

	orders, err := r.getOrdersForAssign(ctx)
	log.Println("len of orders for assign: ", len(orders))
	if err != nil {
		return nil, err
	}

	footCouriers, err := r.getCouriersWithGivenType(ctx, "FOOT")
	log.Println("len of foot couriers for assign: ", len(footCouriers))
	if err != nil {
		return nil, err
	}

	bikeCouriers, err := r.getCouriersWithGivenType(ctx, "BIKE")
	log.Println("len of bike couriers for assign: ", len(bikeCouriers))
	if err != nil {
		return nil, err
	}

	autoCouriers, err := r.getCouriersWithGivenType(ctx, "AUTO")
	log.Println("len of auto couriers for assign: ", len(autoCouriers))
	if err != nil {
		return nil, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Weight > orders[j].Weight
	})

	for _, order := range orders {
		if order.Weight > 20 {
			findAndSetFreeCourier(autoCouriers, order, assignments)
		} else if order.Weight > 10 {
			if !findAndSetFreeCourier(autoCouriers, order, assignments) {
				findAndSetFreeCourier(bikeCouriers, order, assignments)
			}
		} else {
			if !findAndSetFreeCourier(autoCouriers, order, assignments) {
				if !findAndSetFreeCourier(bikeCouriers, order, assignments) {
					findAndSetFreeCourier(footCouriers, order, assignments)
				}
			}
		}
		log.Println(assignments)
	}

	for courierID, regions := range assignments {
		assignment := entity.CourierAssignment{
			CourierID: courierID,
			Orders:    make([]entity.OrdersGroup, 0),
		}

		for _, orderGroup := range regions {
			for _, orderG := range orderGroup.orders {
				for _, order := range orderG.Orders {
					_, err = r.Pool.Exec(ctx, _updateDistributionDateAndCourierIDInOrder, date, courierID, order.OrderID)
					if err != nil {
						return nil, fmt.Errorf("OrderRepo - Assign - r.Pool.Exec(_updateDistributionDateAndCourierIDInOrder): %w", err)
					}
				}
			}
			assignment.Orders = append(assignment.Orders, orderGroup.orders...)
		}

		couriersAssignment = append(couriersAssignment, &assignment)
	}

	return couriersAssignment, nil
}
