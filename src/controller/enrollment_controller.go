package controller

import (
	"courses-service/src/ai"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type EnrollmentController struct {
	enrollmentService service.EnrollmentServiceInterface
	aiClient          *ai.AiClient
	activityService   service.TeacherActivityServiceInterface
}

func NewEnrollmentController(enrollmentService service.EnrollmentServiceInterface, aiClient *ai.AiClient, activityService service.TeacherActivityServiceInterface) *EnrollmentController {
	return &EnrollmentController{
		enrollmentService: enrollmentService,
		aiClient:          aiClient,
		activityService:   activityService,
	}
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
// @Param studentId query string true "Student ID"
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

	studentId := ctx.Query("studentId")
	if studentId == "" {
		slog.Error("Student ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	err := c.enrollmentService.UnenrollStudent(studentId, courseID)
	if err != nil {
		slog.Error("Error unenrolling student", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Student unenrolled from course", "studentId", studentId, "courseId", courseID)
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
// @Param studentId query string true "Student ID"
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

	studentId := ctx.Query("studentId")
	if studentId == "" {
		slog.Error("Student ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}
	err := c.enrollmentService.UnsetFavouriteCourse(studentId, courseID)
	if err != nil {
		slog.Error("Error unsetting favourite course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Favourite course unset", "studentId", studentId, "courseId", courseID)
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
// @Router /feedback/student/{id} [put]
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

// @Summary Get student feedback summary
// @Description Get student feedback summary by student ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} schemas.AiSummaryResponse
// @Router /feedback/student/{id}/summary [get]
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
	ctx.JSON(http.StatusOK, schemas.AiSummaryResponse{Summary: summary})
}

// @Summary Approve a student in a course
// @Description Approve a student by changing their enrollment status to completed
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param studentId path string true "Student ID"
// @Success 200 {object} schemas.ApproveStudentResponse
// @Router /courses/{id}/students/{studentId}/approve [put]
func (c *EnrollmentController) ApproveStudent(ctx *gin.Context) {
	slog.Debug("Approving student", "courseId", ctx.Param("courseId"), "studentId", ctx.Param("studentId"))

	courseID := ctx.Param("id")
	studentID := ctx.Param("studentId")

	if courseID == "" {
		slog.Error("Invalid course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	if studentID == "" {
		slog.Error("Invalid student ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	err := c.enrollmentService.ApproveStudent(studentID, courseID)
	if err != nil {
		slog.Error("Error approving student", "error", err, "studentId", studentID, "courseId", courseID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity if teacher is auxiliary
	teacherUUID := ctx.GetString("teacher_uuid")
	if teacherUUID != "" {
		c.activityService.LogActivityIfAuxTeacher(
			courseID,
			teacherUUID,
			"APPROVE_STUDENT",
			fmt.Sprintf("Approved student: %s", studentID),
		)
	}

	slog.Debug("Student approved successfully", "studentId", studentID, "courseId", courseID)
	ctx.JSON(http.StatusOK, schemas.ApproveStudentResponse{
		Message:   "Student approved successfully",
		StudentID: studentID,
		CourseID:  courseID,
	})
}

// @Summary Disapprove a student in a course
// @Description Disapprove a student by changing their enrollment status to dropped with a reason
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param studentId path string true "Student ID"
// @Param disapproveRequest body schemas.DisapproveStudentRequest true "Disapprove request"
// @Success 200 {object} schemas.DisapproveStudentResponse
// @Router /courses/{id}/students/{studentId}/disapprove [put]
func (c *EnrollmentController) DisapproveStudent(ctx *gin.Context) {
	slog.Debug("Disapproving student", "courseId", ctx.Param("id"), "studentId", ctx.Param("studentId"))

	courseID := ctx.Param("id")
	studentID := ctx.Param("studentId")

	if courseID == "" {
		slog.Error("Invalid course ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	if studentID == "" {
		slog.Error("Invalid student ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	var disapproveRequest schemas.DisapproveStudentRequest
	if err := ctx.ShouldBindJSON(&disapproveRequest); err != nil {
		slog.Error("Error binding disapprove request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.enrollmentService.DisapproveStudent(studentID, courseID, disapproveRequest.Reason)
	if err != nil {
		slog.Error("Error disapproving student", "error", err, "studentId", studentID, "courseId", courseID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity if teacher is auxiliary
	teacherUUID := ctx.GetString("teacher_uuid")
	if teacherUUID != "" {
		c.activityService.LogActivityIfAuxTeacher(
			courseID,
			teacherUUID,
			"DISAPPROVE_STUDENT",
			fmt.Sprintf("Disapproved student: %s (reason: %s)", studentID, disapproveRequest.Reason),
		)
	}

	slog.Debug("Student disapproved successfully", "studentId", studentID, "courseId", courseID, "reason", disapproveRequest.Reason)
	ctx.JSON(http.StatusOK, schemas.DisapproveStudentResponse{
		Message:   "Student disapproved successfully",
		StudentID: studentID,
		CourseID:  courseID,
		Reason:    disapproveRequest.Reason,
	})
}
