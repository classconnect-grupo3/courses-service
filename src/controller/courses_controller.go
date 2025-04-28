package controller

import (
	"log"
	"log/slog"
	"net/http"

	"courses-service/src/schemas"
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

func (c *CoursesController) CreateCourse(ctx *gin.Context) {
	slog.Debug("Creating course")

	var course schemas.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&course); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCourse, err := c.service.CreateCourse(course)
	if err != nil {
		slog.Error("Error creating course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course created", "course", createdCourse)
	ctx.JSON(http.StatusCreated, createdCourse)
}

func (c *CoursesController) GetCourseById(ctx *gin.Context) {
	slog.Debug("Getting course by ID")

	id := ctx.Param("id")
	course, err := c.service.GetCourseById(id)
	if err != nil {
		slog.Error("Error getting course by ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

func (c *CoursesController) DeleteCourse(ctx *gin.Context) {
	slog.Debug("Deleting course")
	id := ctx.Param("id")

	err := c.service.DeleteCourse(id)
	if err != nil {
		slog.Error("Error deleting course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course deleted", "id", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func (c *CoursesController) GetCourseByTeacherId(ctx *gin.Context) {
	slog.Debug("Getting course by teacher ID")
	teacherId := ctx.Param("teacherId")
	log.Printf("The teacher ID is %v", teacherId)
	course, err := c.service.GetCourseByTeacherId(teacherId)
	if err != nil {
		slog.Error("Error getting course by teacher ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

func (c *CoursesController) GetCourseByTitle(ctx *gin.Context) {
	slog.Debug("Getting course by title")
	title := ctx.Param("title")
	course, err := c.service.GetCourseByTitle(title)
	if err != nil {
		slog.Error("Error getting course by title", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}
