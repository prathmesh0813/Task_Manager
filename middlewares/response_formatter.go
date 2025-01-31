package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Set the response in gin context
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		elapsed := time.Since(startTime)
		response, _ := c.Get("response")
		message, _ := c.Get("message")
		errorStatus, _ := c.Get("error")
		requestid, _ := c.Get("request_id")
		formattedResponse := gin.H{
			"message":       message,
			"error":         errorStatus,
			"data":          response,
			"executionTime": elapsed,
			"request_id":    requestid,
		}

		c.JSON(c.Writer.Status(), formattedResponse)

	}
}
