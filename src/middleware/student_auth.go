package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// StudentAuth is a middleware that extracts student information from the X-Student-UUID and X-Student-Name
// headers and sets them in the context for handlers to use.
func StudentAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		studentUUID := c.GetHeader("X-Student-UUID")
		studentName := c.GetHeader("X-Student-Name")

		if studentUUID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-Student-UUID header is required"})
			c.Abort()
			return
		}

		// Set values in context for downstream handlers
		c.Set("student_uuid", studentUUID)
		c.Set("student_name", studentName)

		c.Next()
	}
}
