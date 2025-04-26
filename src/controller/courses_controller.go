package controller

import (
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
	courses, err := c.service.GetCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, courses)
}
