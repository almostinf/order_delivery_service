package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/internal/usecase"
	"github.com/almostinf/order_delivery_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type orderRoutes struct {
	uc usecase.OrderUseCase
	l  logger.Interface
}

func newOrderRoutes(handler *gin.RouterGroup, uc usecase.OrderUseCase, l logger.Interface) {
	r := &orderRoutes{uc, l}

	h := handler.Group("/orders")
	{
		h.GET("/", r.getAll)
		h.GET("/:order_id", r.get)
		h.POST("/", r.create)
		h.POST("/complete", r.complete)
		h.PUT("/set_courier", r.setCourierID)
		h.POST("/assign", r.assign)
	}
}

// @Summary     Get All Orders
// @Description Get All Orders from Postgres
// @ID          get-all-orders
// @Tags  	    orders
// @Produce     json
// @Param       limit query int false "Limit the number of results (default: 1)"
// @Param       offset query int false "Offset the list of results (default: 0)"
// @Success     200 {array} entity.OrderResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/ [get]
func (r *orderRoutes) getAll(c *gin.Context) {
	limit, offset := 1, 0
	if limitStr, ok := c.GetQuery("limit"); ok {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			r.l.Error(err, "http - v1 - order - getAll - strconv.Atoi")
			errorResponse(c, http.StatusBadRequest, "failed conversation offset to int")

			return
		}
		limit = parsedLimit
	}

	if offsetStr, ok := c.GetQuery("offset"); ok {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			r.l.Error(err, "http - v1 - order - getAll - strconv.Atoi")
			errorResponse(c, http.StatusBadRequest, "failed conversation offset to int")

			return
		}
		offset = parsedOffset
	}

	if limit < 0 || offset < 0 {
		r.l.Error(errors.New("offset or limit is less than zero"), "http - v1 - order - getAll")
		errorResponse(c, http.StatusBadRequest, "wrong limit or offset format")

		return
	}

	orders, err := r.uc.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		r.l.Error(err, "http - v1 - order - getAll - GetAll")
		errorResponse(c, http.StatusInternalServerError, "order service problems")

		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary     Get Order By ID
// @Description Get Order By ID from Postgres
// @ID          get-orders-by-id
// @Tags  	    orders
// @Produce     json
// @Param       order_id path string true "Courier ID"
// @Success     200 {object} entity.OrderResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/{order_id} [get]
func (r *orderRoutes) get(c *gin.Context) {
	idStr := c.Param("order_id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		r.l.Error(err, "http - v1 - order - get - uuid.Parse")
		errorResponse(c, http.StatusBadRequest, "failed conversation id (string) to uuid")

		return
	}

	order, err := r.uc.Get(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - order - get - Get")
		errorResponse(c, http.StatusNotFound, "order service problems")

		return
	}

	c.JSON(http.StatusOK, order)
}

func ValidateOrderRequest(order CreateOrderRequest) error {
	if order.Weight <= 0 {
		return errors.New("the order weight can't be less than zero")
	}

	if order.Regions < 0 {
		return errors.New("the order region can't be less than zero")
	}

	for _, deliveryHours := range order.DeliveryHours {
		if err := parseTimeRange(deliveryHours); err != nil {
			return fmt.Errorf("invalid delivery hours: %w", err)
		}
	}

	if order.Cost < 0 {
		return errors.New("the order cost can't be less than zero")
	}

	return nil
}

type CreateOrderRequest struct {
	Weight        float32  `json:"weight" binding:"required"`
	Regions       int      `json:"regions" binding:"required"`
	DeliveryHours []string `json:"delivery_hours" binding:"required"`
	Cost          int      `json:"cost" binding:"required"`
}

// @Summary     Create Order
// @Description Create Order in Postgres
// @ID          create-order
// @Tags  	    orders
// @Accept      json
// @Produce     json
// @Param       request body object true "Order object"
// @Success     200 {object} entity.OrderResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/ [post]
func (r *orderRoutes) create(c *gin.Context) {
	ordersReq := make(map[string][]CreateOrderRequest)

	if err := c.ShouldBindJSON(&ordersReq); err != nil {
		r.l.Error(err, "http - v1 - order - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if _, ok := ordersReq["orders"]; !ok {
		r.l.Error(errors.New("invalid key in request body"), "http - v1 - order - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	for _, orderReq := range ordersReq["orders"] {
		if err := ValidateOrderRequest(orderReq); err != nil {
			r.l.Error(err, "http - v1 - order - create")
			errorResponse(c, http.StatusBadRequest, "invalid request body")

			return
		}
	}

	response := make(map[string][]*entity.OrderResponse)
	response["orders"] = make([]*entity.OrderResponse, 0)

	for _, orderReq := range ordersReq["orders"] {
		order, err := r.uc.Create(
			c.Request.Context(),
			&entity.Order{
				OrderResponse: entity.OrderResponse{
					OrderID:       uuid.New(),
					Weight:        orderReq.Weight,
					Regions:       orderReq.Regions,
					DeliveryHours: orderReq.DeliveryHours,
					Cost:          orderReq.Cost,
					CompletedTime: time.Time{},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)

		if err != nil {
			r.l.Error(err, "http - v1 - order - create")
			errorResponse(c, http.StatusInternalServerError, "order service problems")

			return
		}

		orders := response["orders"]
		orders = append(orders, order)
		response["orders"] = orders
	}

	c.JSON(http.StatusOK, response)
}

// @Summary     Complete Order
// @Description Complete Order
// @ID          complete-order
// @Tags  	    orders
// @Accept      json
// @Produce     json
// @Param       request body object true "Complete Info Order object"
// @Success     200 {array} entity.OrderResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/complete [post]
func (r *orderRoutes) complete(c *gin.Context) {
	completeInfoReq := make(map[string][]entity.CompleteInfo)

	if err := c.ShouldBindJSON(&completeInfoReq); err != nil {
		r.l.Error(err, "http - v1 - order - complete")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if _, ok := completeInfoReq["complete_info"]; !ok {
		r.l.Error(errors.New("invalid key in request body"), "http - v1 - order - complete")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	r.l.Debug(completeInfoReq["complete_info"])

	orders, err := r.uc.Complete(c.Request.Context(), completeInfoReq["complete_info"])
	if err != nil {
		r.l.Error(err, "http - v1 - order - complete")
		errorResponse(c, http.StatusBadRequest, "order service problem")

		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary     Set Courier ID to order
// @Description Set Courier ID to order
// @ID          set-order-courier-id
// @Tags  	    orders
// @Produce     json
// @Param       order_id query string true "Order ID"
// @Param       courier_id query string true "Courier ID"
// @Success     200 {object} entity.OrderResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/set_courier [put]
func (r *orderRoutes) setCourierID(c *gin.Context) {
	orderIDStr, courierIDStr := c.Query("order_id"), c.Query("courier_id")

	courierID, err := uuid.Parse(courierIDStr)
	if err != nil {
		r.l.Error(err, "http - v1 - order - setCourierID - uuid.Parse(courier_id)")
		errorResponse(c, http.StatusBadRequest, "failed conversation courier_id to uuid")

		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		r.l.Error(err, "http - v1 - order - setCourierID - uuid.Parse(order_id)")
		errorResponse(c, http.StatusBadRequest, "failed conversation order_id to uuid")

		return
	}

	order, err := r.uc.SetCourierID(c.Request.Context(), orderID, courierID)
	if err != nil {
		r.l.Error(err, "http - v1 - order - setCourierID")
		errorResponse(c, http.StatusInternalServerError, "order service problem")

		return
	}

	c.JSON(http.StatusOK, order)
}

type couriersAssignResponse struct {
	Date     time.Time                   `json:"date" binding:"require"`
	Couriers []*entity.CourierAssignment `json:"couriers" binding:"require"`
}

// @Summary     Assign Order to Courier
// @Description Assign Order to Courier
// @ID          assign-order
// @Tags  	    orders
// @Produce     json
// @Param       date query string false "Date"
// @Success     200 {object} couriersAssignResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/assign [post]
func (r *orderRoutes) assign(c *gin.Context) {
	date := time.Now().Truncate(24 * time.Hour)
	if dateStr, ok := c.GetQuery("date"); ok {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			r.l.Error(err, "http - v1 - order - assign - time.Parse(date)")
			errorResponse(c, http.StatusBadRequest, "failed conversation date to time")

			return
		}
		date = parsedDate
	}

	courierAssignments, err := r.uc.Assign(c.Request.Context(), date)
	if err != nil {
		r.l.Error(err, "http - v1 - order - assign")
		errorResponse(c, http.StatusInternalServerError, "order service problem")

		return
	}

	response := couriersAssignResponse{
		Date:     date,
		Couriers: courierAssignments,
	}

	c.JSON(http.StatusOK, response)
}
