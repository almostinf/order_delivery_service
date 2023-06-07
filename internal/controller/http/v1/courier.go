package v1

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/almostinf/order_delivery_service/internal/entity"
	"github.com/almostinf/order_delivery_service/internal/usecase"
	"github.com/almostinf/order_delivery_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type courierRoutes struct {
	uc usecase.CourierUseCase
	l  logger.Interface
}

func newCourierRoutes(handler *gin.RouterGroup, uc usecase.CourierUseCase, l logger.Interface) {
	r := &courierRoutes{uc, l}

	h := handler.Group("/couriers")
	{
		h.GET("/", r.getAll)
		h.GET("/:courier_id", r.get)
		h.POST("/", r.create)
		h.GET("/meta-info/:courier_id", r.getMetaInfo)
		h.GET("/assignments", r.getAssignments)
	}
}

type getAllCouriersResponse struct {
	Couriers []*entity.CourierResponse `json:"couriers" binging:"require"`
	Limit    int                       `json:"limit" binding:"require"`
	Offset   int                       `json:"offset" binding:"require"`
}

// @Summary     Get All Couriers
// @Description Get All Couriers from Postgres
// @ID          get-all-couriers
// @Tags  	    couriers
// @Produce     json
// @Param       limit query int false "Limit the number of results (default: 1)"
// @Param       offset query int false "Offset the list of results (default: 0)"
// @Success     200 {object} getAllCouriersResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /couriers/ [get]
func (r *courierRoutes) getAll(c *gin.Context) {
	limit, offset := 1, 0
	if limitStr, ok := c.GetQuery("limit"); ok {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			r.l.Error(err, "http - v1 - courier - getAll - strconv.Atoi")
			errorResponse(c, http.StatusBadRequest, "failed conversation offset to int")

			return
		}
		limit = parsedLimit
	}

	if offsetStr, ok := c.GetQuery("offset"); ok {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			r.l.Error(err, "http - v1 - courier - getAll - strconv.Atoi")
			errorResponse(c, http.StatusBadRequest, "failed conversation offset to int")

			return
		}
		offset = parsedOffset
	}

	if limit < 0 || offset < 0 {
		r.l.Error(errors.New("offset or limit is less than zero"), "http - v1 - courier - getAll")
		errorResponse(c, http.StatusBadRequest, "wrong limit or offset format")

		return
	}

	couriers, err := r.uc.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getAll - GetAll")
		errorResponse(c, http.StatusInternalServerError, "courier service problems")

		return
	}

	response := getAllCouriersResponse{
		Couriers: couriers,
		Limit:    limit,
		Offset:   offset,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary     Get Courier by ID in path
// @Description Get Courier by ID from Postgres
// @ID          get-courier-by-id
// @Tags  	    couriers
// @Produce     json
// @Param       courier_id path string true "Courier ID"
// @Success     200 {object} entity.CourierResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /couriers/{courier_id} [get]
func (r *courierRoutes) get(c *gin.Context) {
	idStr := c.Param("courier_id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - get - uuid.Parse")
		errorResponse(c, http.StatusBadRequest, "failed conversation id (string) to uuid")

		return
	}

	courier, err := r.uc.Get(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - get - Get")
		errorResponse(c, http.StatusNotFound, "courier service problems")

		return
	}

	c.JSON(http.StatusOK, courier)
}

func parseTimeRange(s string) error {
	re := regexp.MustCompile(`^\d{2}:\d{2}-\d{2}:\d{2}$`)
	if !re.MatchString(s) {
		return errors.New("invalid time range format")
	}

	layout := "15:04"
	parts := strings.Split(s, "-")
	startTime, err := time.Parse(layout, parts[0])
	if err != nil {
		return errors.New("invalid start time format")
	}

	endTime, err := time.Parse(layout, parts[1])
	if err != nil {
		return errors.New("invalid end time format")
	}

	if endTime.Sub(startTime) <= 0 {
		return errors.New("the end time must not be less or equal than the start time")
	}

	return nil
}

func validateRegions(regions []int) error {
	if len(regions) == 0 {
		return errors.New("invalid courier regions format (regions are empty)")
	}

	for _, region := range regions {
		if region < 0 {
			return errors.New("invalid courier region format")
		}
	}

	return nil
}

func ValidateCourierRequest(courier CreateCourierRequest) error {
	if courier.CourierType != "FOOT" && courier.CourierType != "BIKE" && courier.CourierType != "AUTO" {
		return errors.New("invalid courier type format")
	}

	if err := validateRegions(courier.Regions); err != nil {
		return err
	}

	for _, workingTime := range courier.WorkingHours {
		if err := parseTimeRange(workingTime); err != nil {
			return fmt.Errorf("invalid working time: %w", err)
		}
	}

	return nil
}

type CreateCourierRequest struct {
	CourierType  string   `json:"courier_type" binding:"required"`
	Regions      []int    `json:"regions" binding:"required"`
	WorkingHours []string `json:"working_hours" binding:"required"`
}

// @Summary     Create Courier
// @Description Create Courier in Postgres
// @ID          create-courier
// @Tags  	    couriers
// @Accept      json
// @Produce     json
// @Param       request body object true "Courier object"
// @Success     200 {object} entity.CourierResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /couriers/ [post]
func (r *courierRoutes) create(c *gin.Context) {
	couriersReq := make(map[string][]CreateCourierRequest)

	if err := c.ShouldBindJSON(&couriersReq); err != nil {
		r.l.Error(err, "http - v1 - courier - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if _, ok := couriersReq["couriers"]; !ok {
		r.l.Error(errors.New("invalid key in request body"), "http - v1 - courier - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	for _, courierReq := range couriersReq["couriers"] {
		if err := ValidateCourierRequest(courierReq); err != nil {
			r.l.Error(err, "http - v1 - courier - create")
			errorResponse(c, http.StatusBadRequest, "invalid request body")

			return
		}
	}

	response := make(map[string][]*entity.CourierResponse)
	response["couriers"] = make([]*entity.CourierResponse, 0)

	for _, courierReq := range couriersReq["couriers"] {
		courier, err := r.uc.Create(
			c.Request.Context(),
			&entity.Courier{
				CourierResponse: entity.CourierResponse{
					CourierID:    uuid.New(),
					CourierType:  courierReq.CourierType,
					Regions:      courierReq.Regions,
					WorkingHours: courierReq.WorkingHours,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)

		if err != nil {
			r.l.Error(err, "http - v1 - courier - create")
			errorResponse(c, http.StatusInternalServerError, "courier service problems")

			return
		}

		couriers := response["couriers"]
		couriers = append(couriers, courier)
		response["couriers"] = couriers
	}

	c.JSON(http.StatusOK, response)
}

// @Summary     Get MetaInfo about Courier
// @Description Get MetaInfo about Courier from Postgres
// @ID          get-courier-metainfo
// @Tags  	    couriers
// @Produce     json
// @Param       courier_id path string true "Courier ID"
// @Param       start_date query string true "Start Date"
// @Param       end_date query string true "End Date"
// @Success     200 {object} entity.CourierMetaInfo
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /couriers/meta-info/{courier_id} [get]
func (r *courierRoutes) getMetaInfo(c *gin.Context) {
	idStr := c.Param("courier_id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getMetaInfo - uuid.Parse")
		errorResponse(c, http.StatusBadRequest, "failed conversation id (string) to uuid")

		return
	}

	startDateStr, endDateStr := c.Query("start_date"), c.Query("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getMetaInfo - time.Parse(start_date)")
		errorResponse(c, http.StatusBadRequest, "failed conversation start_date to time")

		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getMetaInfo - time.Parse(end_date)")
		errorResponse(c, http.StatusBadRequest, "failed conversation end_date to time")

		return
	}

	if endDate.Sub(startDate) <= 0 {
		r.l.Error(err, "http - v1 - courier - getMetaInfo")
		errorResponse(c, http.StatusBadRequest, "the end_date must not be less or equal than the start_date")

		return
	}

	courierMetaInfo, err := r.uc.GetMetaInfo(c.Request.Context(), id, startDate, endDate)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getMetaInfo")
		errorResponse(c, http.StatusInternalServerError, "courier service problems")

		return
	}

	c.JSON(http.StatusOK, courierMetaInfo)
}

// @Summary     Get Assignments of Courier
// @Description Get Assignments of Courier from Postgres
// @ID          get-courier-assignments
// @Tags  	    couriers
// @Produce     json
// @Param       date query string false "Date"
// @Param       courier_id query string false "Courier ID"
// @Success     200 {object} couriersAssignResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /couriers/assignments [get]
func (r *courierRoutes) getAssignments(c *gin.Context) {
	date := time.Now().Truncate(24 * time.Hour)
	isAllCouriers := true
	var courierID uuid.UUID

	if dateStr, ok := c.GetQuery("date"); ok {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, dateStr)
		if err != nil {
			r.l.Error(err, "http - v1 - courier - getAssignments - time.Parse(date)")
			errorResponse(c, http.StatusBadRequest, "failed conversation date to time")

			return
		}
		date = parsedDate
	}

	if courierIDStr, ok := c.GetQuery("courier_id"); ok {
		parsedCourierID, err := uuid.Parse(courierIDStr)
		if err != nil {
			r.l.Error(err, "http - v1 - courier - getAssignments - uuid.Parse(courier_id)")
			errorResponse(c, http.StatusBadRequest, "failed conversation courier_id to uuid")

			return
		}
		isAllCouriers = false
		courierID = parsedCourierID
	}

	courierAssignments, err := r.uc.GetAssignments(c.Request.Context(), date, courierID, isAllCouriers)
	if err != nil {
		r.l.Error(err, "http - v1 - courier - getAssignments")
		errorResponse(c, http.StatusInternalServerError, "courier service problem")

		return
	}

	response := couriersAssignResponse{
		Date:     date,
		Couriers: courierAssignments,
	}

	c.JSON(http.StatusOK, response)
}
