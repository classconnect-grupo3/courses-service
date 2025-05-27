package controller

import (
	"log/slog"
	"net/http"

	"courses-service/src/model"
	"courses-service/src/schemas"

	"github.com/gin-gonic/gin"
)

type AssignmentService interface {
	CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetAssignmentById(id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error)
	DeleteAssignment(id string) error
}

type AssignmentsController struct {
	service AssignmentService
}

func NewAssignmentsController(service AssignmentService) *AssignmentsController {
	return &AssignmentsController{service: service}
}

// @Summary Get all assignments
// @Description Get all assignments
// @Tags assignments
// @Accept json
// @Produce json
// @Router /assignments [get]
// @Success 200 {array} model.Assignment
func (c *AssignmentsController) GetAssignments(ctx *gin.Context) {
	slog.Debug("Getting assignments")

	assignments, err := c.service.GetAssignments()
	if err != nil {
		slog.Error("Error getting assignments", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Assignments retrieved", "assignments", assignments)
	ctx.JSON(http.StatusOK, assignments)
}

// @Summary Create an assignment
// @Description Create an assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignment body schemas.CreateAssignmentRequest true "Assignment to create"
// @Success 201 {object} model.Assignment
// @Router /assignments [post]
func (c *AssignmentsController) CreateAssignment(ctx *gin.Context) {
	slog.Debug("Creating assignment")

	var assignment schemas.CreateAssignmentRequest
	if err := ctx.ShouldBindJSON(&assignment); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAssignment, err := c.service.CreateAssignment(assignment)
	if err != nil {
		slog.Error("Error creating assignment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Assignment created", "assignment", createdAssignment)
	ctx.JSON(http.StatusCreated, createdAssignment)
}

// @Summary Get an assignment by ID
// @Description Get an assignment by ID
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Success 200 {object} model.Assignment
// @Router /assignments/{assignmentId} [get]
func (c *AssignmentsController) GetAssignmentById(ctx *gin.Context) {
	slog.Debug("Getting assignment by ID")
	id := ctx.Param("assignmentId")

	assignment, err := c.service.GetAssignmentById(id)
	if err != nil {
		slog.Error("Error getting assignment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if assignment == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "assignment not found"})
		return
	}

	slog.Debug("Assignment retrieved", "assignment", assignment)
	ctx.JSON(http.StatusOK, assignment)
}

// @Summary Get assignments by course ID
// @Description Get assignments by course ID
// @Tags assignments
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Success 200 {array} model.Assignment
// @Router /assignments/course/{courseId} [get]
func (c *AssignmentsController) GetAssignmentsByCourseId(ctx *gin.Context) {
	slog.Debug("Getting assignments by course ID")
	courseId := ctx.Param("courseId")

	assignments, err := c.service.GetAssignmentsByCourseId(courseId)
	if err != nil {
		slog.Error("Error getting assignments", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Assignments retrieved", "assignments", assignments)
	ctx.JSON(http.StatusOK, assignments)
}

// @Summary Update an assignment
// @Description Update an assignment by ID
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Param assignment body schemas.UpdateAssignmentRequest true "Assignment to update"
// @Success 200 {object} model.Assignment
// @Router /assignments/{assignmentId} [put]
func (c *AssignmentsController) UpdateAssignment(ctx *gin.Context) {
	slog.Debug("Updating assignment")
	id := ctx.Param("assignmentId")

	var updateAssignmentRequest schemas.UpdateAssignmentRequest
	if err := ctx.ShouldBindJSON(&updateAssignmentRequest); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAssignment, err := c.service.UpdateAssignment(id, updateAssignmentRequest)
	if err != nil {
		slog.Error("Error updating assignment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Assignment updated", "assignment", updatedAssignment)
	ctx.JSON(http.StatusOK, updatedAssignment)
}

// @Summary Delete an assignment
// @Description Delete an assignment by ID
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignmentId path string true "Assignment ID"
// @Success 200 {string} string "Assignment deleted successfully"
// @Router /assignments/{assignmentId} [delete]
func (c *AssignmentsController) DeleteAssignment(ctx *gin.Context) {
	slog.Debug("Deleting assignment")
	id := ctx.Param("assignmentId")

	if err := c.service.DeleteAssignment(id); err != nil {
		slog.Error("Error deleting assignment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Assignment deleted")
	ctx.JSON(http.StatusNoContent, nil)
}
