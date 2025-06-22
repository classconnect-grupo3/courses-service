package controller

import (
	"fmt"
	"net/http"
	"time"

	"courses-service/src/model"
	"courses-service/src/queues"
	"courses-service/src/schemas"
	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct {
	submissionService  service.SubmissionServiceInterface
	notificationsQueue queues.NotificationsQueueInterface
}

func NewSubmissionController(submissionService service.SubmissionServiceInterface, notificationsQueue queues.NotificationsQueueInterface) *SubmissionController {
	return &SubmissionController{
		submissionService:  submissionService,
		notificationsQueue: notificationsQueue,
	}
}

type CreateSubmissionRequest struct {
	Answers []model.Answer `json:"answers"`
}

// @Summary Create a submission
// @Description Create a submission
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param submission body CreateSubmissionRequest true "Submission to create"
// @Success 201 {object} model.Submission
// @Router /assignments/{assignmentId}/submissions [post]
func (c *SubmissionController) CreateSubmission(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	var req CreateSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get student info from context (assuming middleware sets this)
	studentUUID := ctx.GetString("student_uuid")
	studentName := ctx.GetString("student_name")

	submission, err := c.submissionService.GetOrCreateSubmission(ctx, assignmentID, studentUUID, studentName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	submission.Answers = req.Answers
	submission.UpdatedAt = time.Now()

	if err := c.submissionService.UpdateSubmission(ctx, submission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submission)
}

// @Summary Get a submission by ID
// @Description Get a submission by ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param id path string true "Submission ID"
// @Success 200 {object} model.Submission
// @Router /assignments/{assignmentId}/submissions/{id} [get]
func (c *SubmissionController) GetSubmission(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	id := ctx.Param("id")

	submission, err := c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if submission == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	// Validate submission belongs to the assignment
	if submission.AssignmentID != assignmentID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	ctx.JSON(http.StatusOK, submission)
}

// @Summary Update a submission
// @Description Update a submission by ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param id path string true "Submission ID"
// @Param submission body model.Submission true "Submission to update"
// @Success 200 {object} model.Submission
// @Router /assignments/{assignmentId}/submissions/{id} [put]
func (c *SubmissionController) UpdateSubmission(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	id := ctx.Param("id")

	var submission model.Submission
	if err := ctx.ShouldBindJSON(&submission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate submission ID matches URL
	if submission.ID.Hex() != id {
		fmt.Printf("submission ID mismatch: %s != %s\n", submission.ID.Hex(), id)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "submission ID mismatch"})
		return
	}

	// Validate submission belongs to the assignment
	if submission.AssignmentID != assignmentID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "assignment ID mismatch"})
		return
	}

	// Validate student ownership
	studentUUID := ctx.GetString("student_uuid")
	if submission.StudentUUID != studentUUID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.submissionService.UpdateSubmission(ctx, &submission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submission)
}

// @Summary Submit a submission
// @Description Submit a submission by ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param id path string true "Submission ID"
// @Success 200 {object} model.Submission
// @Router /assignments/{assignmentId}/submissions/{id}/submit [post]
func (c *SubmissionController) SubmitSubmission(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	id := ctx.Param("id")

	// Validate submission belongs to the assignment
	submission, err := c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if submission == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}
	if submission.AssignmentID != assignmentID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	if err := c.submissionService.SubmitSubmission(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated submission after auto-correction
	updatedSubmission, err := c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send notification about the corrected submission
	c.sendCorrectionNotification(updatedSubmission, assignmentID)

	ctx.JSON(http.StatusOK, updatedSubmission)
}

// sendCorrectionNotification sends a notification about the corrected submission
func (c *SubmissionController) sendCorrectionNotification(submission *model.Submission, assignmentID string) {
	if c.notificationsQueue == nil {
		return // Skip if notifications are not configured
	}

	correctionType := "automatic"
	needsManualReview := false

	if submission.NeedsManualReview != nil && *submission.NeedsManualReview {
		correctionType = "needs_manual_review"
		needsManualReview = true
	}

	queueMessage := queues.NewSubmissionCorrectedMessage(
		assignmentID,
		submission.ID.Hex(),
		submission.StudentUUID,
		submission.Score,
		submission.Feedback,
		submission.AIScore,
		submission.AIFeedback,
		correctionType,
		needsManualReview,
	)

	// if err := c.notificationsQueue.Publish(queueMessage); err != nil {
	// 	// Log the error but don't fail the response
	// 	fmt.Printf("Error publishing correction notification: %v\n", err)
	// } TODO: Uncomment this when the notifications queue is implemented
	fmt.Println("queueMessage: ", queueMessage)
}

// @Summary Get submissions by assignment ID
// @Description Get submissions by assignment ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Success 200 {array} model.Submission
// @Router /assignments/{assignmentId}/submissions [get]
func (c *SubmissionController) GetSubmissionsByAssignment(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")

	submissions, err := c.submissionService.GetSubmissionsByAssignment(ctx, assignmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submissions)
}

// @Summary Get submissions by student ID
// @Description Get submissions by student ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param studentUUID path string true "Student ID"
// @Success 200 {array} model.Submission
// @Router /students/{studentUUID}/submissions [get]
func (c *SubmissionController) GetSubmissionsByStudent(ctx *gin.Context) {
	studentUUID := ctx.Param("studentUUID")

	// Validate student access
	requestingStudentUUID := ctx.GetString("student_uuid")
	if studentUUID != requestingStudentUUID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	submissions, err := c.submissionService.GetSubmissionsByStudent(ctx, studentUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submissions)
}

// @Summary Grade a submission
// @Description Grade a submission by ID (for teachers)
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param id path string true "Submission ID"
// @Param gradeRequest body schemas.GradeSubmissionRequest true "Grade request"
// @Success 200 {object} model.Submission
// @Router /assignments/{assignmentId}/submissions/{id}/grade [put]
func (c *SubmissionController) GradeSubmission(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	id := ctx.Param("id")

	var gradeRequest schemas.GradeSubmissionRequest
	if err := ctx.ShouldBindJSON(&gradeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get teacher info from context
	teacherUUID := ctx.GetString("teacher_uuid")

	// Validate teacher permissions for this assignment
	if err := c.submissionService.ValidateTeacherPermissions(ctx, assignmentID, teacherUUID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Validate submission belongs to the assignment
	submission, err := c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if submission == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}
	if submission.AssignmentID != assignmentID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	// Grade the submission
	gradedSubmission, err := c.submissionService.GradeSubmission(ctx, id, gradeRequest.Score, gradeRequest.Feedback)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gradedSubmission)
}

// @Summary Generate feedback summary
// @Description Generate an AI summary of the feedback for a submission
// @Tags submissions
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param id path string true "Submission ID"
// @Success 200 {object} schemas.AiSummaryResponse
// @Router /assignments/{assignmentId}/submissions/{id}/feedback-summary [get]
func (c *SubmissionController) GenerateFeedbackSummary(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")
	id := ctx.Param("id")

	// Get teacher info from context
	teacherUUID := ctx.GetString("teacher_uuid")

	// Validate teacher permissions for this assignment
	if err := c.submissionService.ValidateTeacherPermissions(ctx, assignmentID, teacherUUID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Validate submission belongs to the assignment
	submission, err := c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if submission == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}
	if submission.AssignmentID != assignmentID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	// Generate feedback summary
	summary, err := c.submissionService.GenerateFeedbackSummary(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
