package integration_tests

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	_ "github.com/lib/pq"
)

const (
	// Attempts connection
	host       = "test:8080"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()

	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP /couriers:

// HTTP GET: /couriers
func TestHTTPGetAllCouriers(t *testing.T) {
	Test(t,
		Description("get all couriers success"),
		Get(basePath+"/couriers?limit=10&offset=10"),
		Expect().Status().Equal(http.StatusOK),
		// Expect().Body().String().Contains(`[{`),
	)

	Test(t,
		Description("forgot to put an offset in the request"),
		Get(basePath+"/couriers?limit=10"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("limit less than zero"),
		Get(basePath+"/couriers?limit=-10&offset=10"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("wrong limit or offset format"),
	)
}

// HTTP GET: /couriers/:courier_id
func TestHTTPGetCourier(t *testing.T) {
	// 9789176b-966b-44b3-b52a-1dde8b2fdc3f
	Test(t,
		Description("invalid courier_id"),
		Get(basePath+"/couriers/afdsaf"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation id (string) to int64"),
	)
}

// HTTP POST: /couriers
func TestHTTPCreateCourier(t *testing.T) {
	body := `
	{
		"couriers": [
			{
				"courier_type": "AUTO",
				"regions": [1, 2, 3],
				"working_hours": ["10:00-23:00", "15:00-18:00"]
			},
			{
				"courier_type": "AUTO",
				"regions": [1, 3, 4],
				"working_hours": ["19:00-20:00", "14:00-19:00"]
			}
		]
	}
	`

	Test(t,
		Description("create courier success"),
		Post(basePath+"/couriers/"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
	)

	body = `
	{
		"couriers": [
			{
				"courier_type": "AuUTO",
				"regions": [1, 2, 3],
				"working_hours": ["10:00--25:00", "15:00-18:00"]
			},
			{
				"courier_type": "AUTO",
				"regions": [1, -3, 4],
				"working_hours": ["19:00-20:00", "14:00-19:00"]
			}
		]
	}
	`

	Test(t,
		Description("invalid courier"),
		Post(basePath+"/couriers/"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("invalid request body"),
	)
}

// HTTP GET: /couriers/meta-info/:courier_id
func TestHTTPGetCourierMetaInfo(t *testing.T) {
	Test(t,
		Description("forgot to put the courier_id"),
		Get(basePath+"/couriers/meta-info/"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation id (string) to int64"),
	)

	Test(t,
		Description("invalid courier_id"),
		Get(basePath+"/couriers/meta-info/afdsf"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation id (string) to int64"),
	)

	Test(t,
		Description("invalid start_date"),
		Get(basePath+`/couriers/meta-info/1?start_date=2006-13-02&end_date=2006-01-02`),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation start_date to time"),
	)

	Test(t,
		Description("invalid end_date"),
		Get(basePath+`/couriers/meta-info/1?start_date=2006-11-02&end_date=2006-11-34`),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation end_date to time"),
	)

	Test(t,
		Description("the end_date is less than the start_date"),
		Get(basePath+"/couriers/meta-info/1?start_date=2006-12-02&end_date=2006-01-02"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("the end_date must not be less or equal than the start_date"),
	)
}

// HTTP GET: /couriers/assignments
func TestHTTPGetCouriersAssignments(t *testing.T) {
	Test(t,
		Description("get couriers assignments successfully"),
		Get(basePath+"/couriers/assignments"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("get couriers assignments successfully with date"),
		Get(basePath+"/couriers/assignments?date=2019-08-24"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("wrong courier_id"),
		Get(basePath+"/couriers/assignments?courier_id=adffafd"),
		Expect().Status().Equal(http.StatusBadRequest),
	)
}

// HTTP /orders:

// HTTP GET: /orders
func TestHTTPGetAllOrders(t *testing.T) {
	Test(t,
		Description("get all orders success"),
		Get(basePath+"/orders?limit=10&offset=10"),
		Expect().Status().Equal(http.StatusOK),
		// Expect().Body().String().Contains(`[{`),
	)

	Test(t,
		Description("forgot to put an offset in the request"),
		Get(basePath+"/orders?limit=10"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("limit less than zero"),
		Get(basePath+"/orders?limit=-10&offset=10"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("wrong limit or offset format"),
	)
}

// HTTP GET: /orders/:order_id
func TestHTTPGetOrder(t *testing.T) {
	// 9789176b-966b-44b3-b52a-1dde8b2fdc3f
	Test(t,
		Description("invalid order_id"),
		Get(basePath+"/orders/afdsfa"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("failed conversation id (string) to int64"),
	)
}

// HTTP POST: /orders
func TetsHTTPCreateOrder(t *testing.T) {
	body := `
	{
		"orders": [
			{
				"weight": 35,
				"regions": 1,
				"delivery_hours": [
					"02:00-12:00",
					"15:00-19:00"
				],
				"cost": 1600
			},
			{
				"weight": 5,
				"regions": 1,
				"delivery_hours": [
					"15:00-16:00",
					"19:00-21:00"
				],
				"cost": 1000
			}
		]
	}
	`

	Test(t,
		Description("create order success"),
		Post(basePath+"/orders/"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
	)

	body = `
	{
		"orders": [
			{
				"weight": 35,
				"regions": -1,
				"delivery_hours": [
					"31:00-12:00",
					"15:00-19:00"
				],
				"cost": 1600
			},
			{
				"weight": 5,
				"regions": 1,
				"delivery_hours": [
					"15:00-16:00",
					"19:00--21:00"
				],
				"cost": 1000
			}
		]
	}
	`

	Test(t,
		Description("invalid order"),
		Post(basePath+"/orders/"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("invalid request body"),
	)
}

// HTTP POST: /orders/complete
func TestHTTPCompleteOrder(t *testing.T) {
	body := `
	{
		"complete_info": [
			{
				"courier_id": 1,
				"order_id": 2,
				"complete_time": "2019-08-25T14:15:22Z"
			},
			{
				"courier_id": 1,
				"order_id": 3,
				"complete_time": "2019-09-24T14:15:22Z"
			}
		]
	}
	`

	Test(t,
		Description("complete fails with internal error"),
		Post(basePath+"/orders/complete"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
	)

	body = `
	{
		"complete_info": [
			{
				"courier_id": afdsaf,
				"order_id": 1,
				"complete_time": "2019-08-25T14:15:22Z"
			}
		]
	}
	`

	Test(t,
		Description("invalid courier_id"),
		Post(basePath+"/orders/complete"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
	)
}

// HTTP POST: /orders/assign
func TestHTTPAssignOrders(t *testing.T) {
	Test(t,
		Description("assign orders successfully"),
		Post(basePath+"/orders/assign"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("assign orders successfully with date"),
		Post(basePath+"/orders/assign?date=2019-08-24"),
		Expect().Status().Equal(http.StatusOK),
	)

	Test(t,
		Description("wrong date"),
		Post(basePath+"/orders/assign?date=-2019-08--24"),
		Expect().Status().Equal(http.StatusBadRequest),
	)
}
