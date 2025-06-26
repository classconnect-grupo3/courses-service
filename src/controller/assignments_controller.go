package controller

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"courses-service/src/queues"
	"courses-service/src/schemas"
	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type AssignmentsController struct {
	service            service.AssignmentServiceInterface
	notificationsQueue queues.NotificationsQueueInterface
	activityService    service.TeacherActivityServiceInterface
}

func NewAssignmentsController(
	service service.AssignmentServiceInterface,
	notificationsQueue queues.NotificationsQueueInterface,
	activityService service.TeacherActivityServiceInterface,
) *AssignmentsController {
	return &AssignmentsController{
		service:            service,
		notificationsQueue: notificationsQueue,
		activityService:    activityService,
	}
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
	log.Println("Creating assignment")

	var assignment schemas.CreateAssignmentRequest
	if err := ctx.ShouldBindJSON(&assignment); err != nil {
		log.Println("Error binding JSON:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAssignment, err := c.service.CreateAssignment(assignment)
	if err != nil {
		log.Println("Error creating assignment:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity if teacher is auxiliary
	teacherUUID := ctx.GetHeader("X-Teacher-UUID")
	if teacherUUID != "" {
		c.activityService.LogActivityIfAuxTeacher(
			createdAssignment.CourseID,
			teacherUUID,
			"CREATE_ASSIGNMENT",
			fmt.Sprintf("Created assignment: %s", createdAssignment.Title),
		)
	}

	queueMessage := queues.NewAssignmentCreatedMessage(
		createdAssignment.CourseID,
		createdAssignment.ID.Hex(),
		createdAssignment.Title,
		createdAssignment.DueDate,
	)
	err = c.notificationsQueue.Publish(queueMessage)
	if err != nil {
		log.Println("Error publishing message:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Assignment created:", createdAssignment.ID)
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

	// Log activity if teacher is auxiliary
	teacherUUID := ctx.GetHeader("X-Teacher-UUID")
	if teacherUUID != "" {
		c.activityService.LogActivityIfAuxTeacher(
			updatedAssignment.CourseID,
			teacherUUID,
			"UPDATE_ASSIGNMENT",
			fmt.Sprintf("Updated assignment: %s", updatedAssignment.Title),
		)
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

	// Get assignment info before deleting for logging
	assignment, err := c.service.GetAssignmentById(id)
	if err != nil {
		slog.Error("Error getting assignment for deletion", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.DeleteAssignment(id); err != nil {
		slog.Error("Error deleting assignment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity if teacher is auxiliary
	teacherUUID := ctx.GetHeader("X-Teacher-UUID")
	if teacherUUID != "" && assignment != nil {
		c.activityService.LogActivityIfAuxTeacher(
			assignment.CourseID,
			teacherUUID,
			"DELETE_ASSIGNMENT",
			fmt.Sprintf("Deleted assignment: %s", assignment.Title),
		)
	}

	slog.Debug("Assignment deleted")
	ctx.JSON(http.StatusNoContent, nil)
}
