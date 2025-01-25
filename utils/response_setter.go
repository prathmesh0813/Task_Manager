package utils

import "github.com/gin-gonic/gin"

// SetErrorResponse sets the response, message, error, and status for a given context.
func SetResponse(c *gin.Context, response any, message string, err bool, statusCode int) {
	c.Set("response", response)
	c.Set("message", message)
	c.Set("error", err)
	c.Status(statusCode)
}
