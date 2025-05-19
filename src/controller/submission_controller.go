package controller

import (
	"net/http"
	"time"

	"courses-service/src/model"
	"courses-service/src/service"
	"github.com/gin-gonic/gin"
)

type SubmissionController struct {
	submissionService *service.SubmissionService
}

func NewSubmissionController(submissionService *service.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}

type CreateSubmissionRequest struct {
	Answers []model.Answer `json:"answers"`
}

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

	submission, err = c.submissionService.GetSubmission(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submission)
}

func (c *SubmissionController) GetSubmissionsByAssignment(ctx *gin.Context) {
	assignmentID := ctx.Param("assignmentId")

	submissions, err := c.submissionService.GetSubmissionsByAssignment(ctx, assignmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, submissions)
}

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