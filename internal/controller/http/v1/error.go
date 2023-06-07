package v1

import (
	"github.com/gin-gonic/gin"
)

// The response struct represents the format of error responses that will be
// sent by the API.
type response struct {
	Error string `json:"error" example:"message"`
}

// The errorResponse function is a utility function that generates an error
// response with the given msg and HTTP code, and aborts the current
// request with that response.
func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
