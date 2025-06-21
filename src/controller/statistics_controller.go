package controller

import (
	"courses-service/src/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StatisticsController handles requests for statistics data
type StatisticsController struct {
	statisticsService service.StatisticsServiceInterface
}

// NewStatisticsController creates a new statistics controller
func NewStatisticsController(statisticsService service.StatisticsServiceInterface) *StatisticsController {
	return &StatisticsController{
		statisticsService: statisticsService,
	}
}

// GetCourseStatistics returns statistics for a specific course
func (c *StatisticsController) GetCourseStatistics(ctx *gin.Context) {
	courseID := ctx.Param("courseId")
	if courseID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	// Parse time range from query parameters (optional)
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format, use YYYY-MM-DD"})
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format, use YYYY-MM-DD"})
			return
		}
	}

	data, _, err := c.statisticsService.ExportCourseStatsCSV(ctx, courseID, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"csv": string(data)})
}

// Returns statistics for a specific student
func (c *StatisticsController) GetStudentStatistics(ctx *gin.Context) {
	studentID := ctx.Param("studentId")
	if studentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
		return
	}

	courseID := ctx.Query("course_id")
	if courseID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	// Parse time range from query parameters (optional)
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format, use YYYY-MM-DD"})
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format, use YYYY-MM-DD"})
			return
		}
	}

	data, _, err := c.statisticsService.ExportStudentStatsCSV(ctx, studentID, courseID, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"csv": string(data)})
}

// GetTeacherCoursesStatistics returns aggregated statistics for all courses of a teacher
func (c *StatisticsController) GetTeacherCoursesStatistics(ctx *gin.Context) {
	teacherID := ctx.Param("teacherId")
	if teacherID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "teacher_id is required"})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format, use YYYY-MM-DD"})
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format, use YYYY-MM-DD"})
			return
		}
	}

	data, _, err := c.statisticsService.ExportTeacherCoursesStatsCSV(ctx, teacherID, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"csv": string(data)})
}
