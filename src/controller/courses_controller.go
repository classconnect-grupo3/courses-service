package controller

import (
	"log/slog"
	"net/http"

	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type CoursesController struct {
	service *service.CourseService
}

func NewCoursesController(service *service.CourseService) *CoursesController {
	return &CoursesController{service: service}
}

func (c *CoursesController) GetCourses(ctx *gin.Context) {
	slog.Debug("Getting courses")
	courses, err := c.service.GetCourses()
	if err != nil {
		slog.Error("Error getting courses", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Courses retrieved", "courses", courses)
	ctx.JSON(http.StatusOK, courses)
}

