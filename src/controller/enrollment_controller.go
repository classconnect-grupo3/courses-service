package controller

import (
	"courses-service/src/schemas"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnrollmentService interface {
	EnrollStudent(studentID, courseID string) error
	UnenrollStudent(studentID, courseID string) error
}

type EnrollmentController struct {
	enrollmentService EnrollmentService
}

func NewEnrollmentController(enrollmentService EnrollmentService) *EnrollmentController {
	return &EnrollmentController{enrollmentService: enrollmentService}
}

// @Summary Enroll a student in a course
// @Description Enroll a student in a course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param enrollmentRequest body schemas.EnrollStudentRequest true "Enrollment request"
// @Router /courses/{id}/enroll [post]
func (c *EnrollmentController) EnrollStudent(ctx *gin.Context) {
	slog.Debug("Enrolling student", "studentId", ctx.Param("studentId"), "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")

	if courseID == "" {
		slog.Error("Invalid course ID", "courseId", courseID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var enrollmentRequest schemas.EnrollStudentRequest
	if err := ctx.ShouldBindJSON(&enrollmentRequest); err != nil {
		slog.Error("Error binding enrollment request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.enrollmentService.EnrollStudent(enrollmentRequest.StudentID, courseID)
	if err != nil {
		slog.Error("Error enrolling student", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Student enrolled in course", "studentId", enrollmentRequest.StudentID, "courseId", courseID)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Student successfully enrolled in course"})
}

// @Summary Unenroll a student from a course
// @Description Unenroll a student from a course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param unenrollmentRequest body schemas.UnenrollStudentRequest true "Unenrollment request"
// @Success 200 {object} schemas.UnenrollStudentResponse
// @Router /courses/{id}/unenroll [delete]
func (c *EnrollmentController) UnenrollStudent(ctx *gin.Context) {
	slog.Debug("Unenrolling student", "studentId", ctx.Param("studentId"), "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")
	studentID := ctx.Param("studentId")

	if studentID == "" || courseID == "" {
		slog.Error("Invalid student ID or course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID or course ID"})
		return
	}

	err := c.enrollmentService.UnenrollStudent(studentID, courseID)
	if err != nil {
		slog.Error("Error unenrolling student", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Student unenrolled from course", "studentId", studentID, "courseId", courseID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Student successfully unenrolled from course"})
}
