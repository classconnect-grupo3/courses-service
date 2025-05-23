package controller

import (
	"log/slog"
	"net/http"

	"courses-service/src/model"
	"courses-service/src/schemas"

	"github.com/gin-gonic/gin"
)

type CourseService interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error)
}

type CourseController struct {
	service CourseService
}

func NewCourseController(service CourseService) *CourseController {
	return &CourseController{service: service}
}

func (c *CourseController) GetCourses(ctx *gin.Context) {
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

func (c *CourseController) CreateCourse(ctx *gin.Context) {
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

func (c *CourseController) GetCourseById(ctx *gin.Context) {
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

func (c *CourseController) DeleteCourse(ctx *gin.Context) {
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

func (c *CourseController) GetCourseByTeacherId(ctx *gin.Context) {
	slog.Debug("Getting course by teacher ID")
	teacherId := ctx.Param("teacherId")
	course, err := c.service.GetCourseByTeacherId(teacherId)
	if err != nil {
		slog.Error("Error getting course by teacher ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

func (c *CourseController) GetCourseByTitle(ctx *gin.Context) {
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

func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	slog.Debug("Updating course")
	id := ctx.Param("id")

	var updateCourseRequest schemas.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&updateCourseRequest); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCourse, err := c.service.UpdateCourse(id, updateCourseRequest)
	if err != nil {
		slog.Error("Error updating course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course updated", "course", updatedCourse)
	ctx.JSON(http.StatusOK, updatedCourse)
}
