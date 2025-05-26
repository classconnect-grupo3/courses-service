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

func (c *EnrollmentController) EnrollStudent(ctx *gin.Context) {
	slog.Debug("Enrolling student", "studentId", ctx.Param("studentId"), "courseId", ctx.Param("courseId"))
	courseID := ctx.Param("courseId")

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

func (c *EnrollmentController) UnenrollStudent(ctx *gin.Context) {
	slog.Debug("Unenrolling student", "studentId", ctx.Param("studentId"), "courseId", ctx.Param("courseId"))
	courseID := ctx.Param("courseId")
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
