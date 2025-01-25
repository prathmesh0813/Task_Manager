package middlewares

import "github.com/gin-gonic/gin"

// Set the response in gin context
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		response, _ := c.Get("response")
		message, _ := c.Get("message")

		errorStatus, _ := c.Get("error")

		formattedResponse := gin.H{
			"message": message,
			"error":   errorStatus,
			"data":    response,
		}

		c.JSON(c.Writer.Status(), formattedResponse)

	}
}
