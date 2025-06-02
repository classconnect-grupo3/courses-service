package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TeacherAuth is a middleware that extracts teacher information from the X-Teacher-UUID and X-Teacher-Name
// headers and sets them in the context for handlers to use.
func TeacherAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherUUID := c.GetHeader("X-Teacher-UUID")
		teacherName := c.GetHeader("X-Teacher-Name")

		if teacherUUID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-Teacher-UUID header is required"})
			c.Abort()
			return
		}

		// Set values in context for downstream handlers
		c.Set("teacher_uuid", teacherUUID)
		c.Set("teacher_name", teacherName)

		c.Next()
	}
} 