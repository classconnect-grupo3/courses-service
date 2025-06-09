package controller

import (
	"courses-service/src/ai"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type EnrollmentController struct {
	enrollmentService service.EnrollmentServiceInterface
	aiClient          *ai.AiClient
}

func NewEnrollmentController(enrollmentService service.EnrollmentServiceInterface, aiClient *ai.AiClient) *EnrollmentController {
	return &EnrollmentController{enrollmentService: enrollmentService, aiClient: aiClient}
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

	if courseID == "" {
		slog.Error("Invalid student ID or course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID or course ID"})
		return
	}

	var unenrollmentRequest schemas.UnenrollStudentRequest
	if err := ctx.ShouldBindJSON(&unenrollmentRequest); err != nil {
		slog.Error("Error binding unenrollment request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.enrollmentService.UnenrollStudent(unenrollmentRequest.StudentID, courseID)
	if err != nil {
		slog.Error("Error unenrolling student", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Student unenrolled from course", "studentId", unenrollmentRequest.StudentID, "courseId", courseID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Student successfully unenrolled from course"})
}

// @Summary Get enrollments by course ID
// @Description Get enrollments by course ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {array} model.Enrollment
// @Router /courses/{id}/enrollments [get]
func (c *EnrollmentController) GetEnrollmentsByCourseId(ctx *gin.Context) {
	slog.Debug("Getting enrollments by course ID", "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")

	enrollments, err := c.enrollmentService.GetEnrollmentsByCourseId(courseID)
	if err != nil {
		slog.Error("Error getting enrollments by course ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, enrollments)
}

// @Summary Set a course as favourite
// @Description Set a course as favourite
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param favouriteCourseRequest body schemas.SetFavouriteCourseRequest true "Favourite course request"
// @Success 200 {object} schemas.SetFavouriteCourseResponse
// @Router /courses/{id}/favourite [post]
func (c *EnrollmentController) SetFavouriteCourse(ctx *gin.Context) {
	slog.Debug("Setting favourite course", "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")

	if courseID == "" {
		slog.Error("Invalid course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var favouriteCourseRequest schemas.SetFavouriteCourseRequest
	if err := ctx.ShouldBindJSON(&favouriteCourseRequest); err != nil {
		slog.Error("Error binding favourite course request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.enrollmentService.SetFavouriteCourse(favouriteCourseRequest.StudentID, courseID)
	if err != nil {
		slog.Error("Error setting favourite course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Favourite course set", "studentId", favouriteCourseRequest.StudentID, "courseId", courseID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Favourite course set"})
}

// @Summary Unset a course as favourite
// @Description Unset a course as favourite
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param unsetFavouriteCourseRequest body schemas.UnsetFavouriteCourseRequest true "Unset favourite course request"
// @Success 200 {object} schemas.UnsetFavouriteCourseResponse
// @Router /courses/{id}/favourite [delete]
func (c *EnrollmentController) UnsetFavouriteCourse(ctx *gin.Context) {
	slog.Debug("Unsetting favourite course", "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")

	if courseID == "" {
		slog.Error("Invalid course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var unsetFavouriteCourseRequest schemas.UnsetFavouriteCourseRequest
	if err := ctx.ShouldBindJSON(&unsetFavouriteCourseRequest); err != nil {
		slog.Error("Error binding unset favourite course request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.enrollmentService.UnsetFavouriteCourse(unsetFavouriteCourseRequest.StudentID, courseID)
	if err != nil {
		slog.Error("Error unsetting favourite course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Favourite course unset", "studentId", unsetFavouriteCourseRequest.StudentID, "courseId", courseID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Favourite course unset"})
}

// @Summary Create a feedback for a course
// @Description Create a feedback for a course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param feedbackRequest body schemas.CreateStudentFeedbackRequest true "Feedback request"
// @Success 200 {object} model.StudentFeedback
// @Router /courses/{id}/student-feedback [post]
func (c *EnrollmentController) CreateFeedback(ctx *gin.Context) {
	slog.Debug("Creating feedback", "courseId", ctx.Param("id"))
	courseID := ctx.Param("id")

	if courseID == "" {
		slog.Error("Invalid course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var feedbackRequest schemas.CreateStudentFeedbackRequest
	if err := ctx.ShouldBindJSON(&feedbackRequest); err != nil {
		slog.Error("Error binding feedback request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(model.FeedbackTypes, feedbackRequest.FeedbackType) {
		slog.Error("Invalid feedback type")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback type"})
		return
	}

	err := c.enrollmentService.CreateStudentFeedback(feedbackRequest)
	if err != nil {
		slog.Error("Error creating feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Feedback created", "studentId", feedbackRequest.StudentUUID, "teacherId", feedbackRequest.TeacherUUID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Feedback created"})
}

// @Summary Get feedback by student ID
// @Description Get feedback by student ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param getFeedbackByStudentIdRequest body schemas.GetFeedbackByStudentIdRequest true "Get feedback by student ID request"
// @Success 200 {array} model.StudentFeedback
// @Router /feedback/student/{id} [get]
func (c *EnrollmentController) GetFeedbackByStudentId(ctx *gin.Context) {
	slog.Debug("Getting feedback by student ID", "studentId", ctx.Param("id"))
	studentID := ctx.Param("id")

	if studentID == "" {
		slog.Error("Invalid student ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest
	if err := ctx.ShouldBindJSON(&getFeedbackByStudentIdRequest); err != nil {
		slog.Error("Error binding get feedback by student ID request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := c.enrollmentService.GetFeedbackByStudentId(studentID, getFeedbackByStudentIdRequest)
	if err != nil {
		slog.Error("Error getting feedback by student ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Feedback retrieved", "studentId", studentID)
	ctx.JSON(http.StatusOK, feedback)
}

func (c *EnrollmentController) GetStudentFeedbackSummary(ctx *gin.Context) {
	slog.Debug("Getting student feedback summary", "studentId", ctx.Param("id"))
	studentID := ctx.Param("id")

	if studentID == "" {
		slog.Error("Invalid student ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	feedbacks, err := c.enrollmentService.GetFeedbackByStudentId(studentID, schemas.GetFeedbackByStudentIdRequest{})
	if err != nil {
		slog.Error("Error getting student feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(feedbacks) == 0 {
		slog.Error("No feedbacks found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No feedbacks found"})
		return
	}

	summary, err := c.aiClient.SummarizeStudentFeedbacks(feedbacks)
	if err != nil {
		slog.Error("Error summarizing student feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Student feedback summary retrieved", "summary", summary)
	ctx.JSON(http.StatusOK, summary)
}
