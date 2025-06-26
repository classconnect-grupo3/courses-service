package controller

import (
	"log/slog"
	"net/http"

	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type TeacherActivityController struct {
	activityService service.TeacherActivityServiceInterface
	courseService   service.CourseServiceInterface
}

func NewTeacherActivityController(
	activityService service.TeacherActivityServiceInterface,
	courseService service.CourseServiceInterface,
) *TeacherActivityController {
	return &TeacherActivityController{
		activityService: activityService,
		courseService:   courseService,
	}
}

// @Summary Get course activity logs
// @Description Get activity logs for auxiliary teachers in a course (only for titular teacher)
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Param teacherId query string true "Teacher ID"
// @Success 200 {array} model.TeacherActivityLog
// @Router /courses/{courseId}/activity-logs [get]
func (c *TeacherActivityController) GetCourseActivityLogs(ctx *gin.Context) {
	slog.Debug("Getting course activity logs")
	courseID := ctx.Param("courseId")
	teacherID := ctx.Query("teacherId")

	if courseID == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	if teacherID == "" {
		slog.Error("Teacher ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Teacher ID is required"})
		return
	}

	// Verify that the requesting teacher is the titular teacher of the course
	course, err := c.courseService.GetCourseById(courseID)
	if err != nil {
		slog.Error("Error getting course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if course.TeacherUUID != teacherID {
		slog.Error("Only titular teacher can access activity logs")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only titular teacher can access activity logs"})
		return
	}

	logs, err := c.activityService.GetCourseActivityLogs(courseID)
	if err != nil {
		slog.Error("Error getting activity logs", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Activity logs retrieved", "courseID", courseID, "count", len(logs))
	ctx.JSON(http.StatusOK, logs)
}
