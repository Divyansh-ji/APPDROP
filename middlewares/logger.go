package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		c.Next()

		latency := time.Since(start)
		statusCode = c.Writer.Status()

		log.Printf("[%s] %d %s %s %s (%s)",
			method, statusCode, path, clientIP, latency, c.Errors.String(),
		)
	}
}
