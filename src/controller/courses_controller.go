package controller

import (
	"log/slog"
	"net/http"
	"slices"

	"courses-service/src/ai"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	service  service.CourseServiceInterface
	aiClient *ai.AiClient
}

func NewCourseController(service service.CourseServiceInterface, aiClient *ai.AiClient) *CourseController {
	return &CourseController{service: service, aiClient: aiClient}
}

// @Summary Get all courses
// @Description Get all courses available in the database
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {array} model.Course
// @Router /courses [get]
func (c *CourseController) GetCourses(ctx *gin.Context) {
	slog.Debug("Getting courses")

	courses, err := c.service.GetCourses()
	if err != nil {
		slog.Error("Error getting courses", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Courses retrieved", "courses", courses)
	ctx.JSON(http.StatusOK, courses)
}

// @Summary Course creation
// @Description Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body schemas.CreateCourseRequest true "Course to create"
// @Success 201 {object} model.Course
// @Router /courses [post]
func (c *CourseController) CreateCourse(ctx *gin.Context) {
	slog.Debug("Creating course")

	var course schemas.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&course); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCourse, err := c.service.CreateCourse(course)
	if err != nil {
		slog.Error("Error creating course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course created", "course", createdCourse)
	ctx.JSON(http.StatusCreated, createdCourse)
}

// @Summary Get a course by ID
// @Description Get a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} model.Course
// @Router /courses/{id} [get]
func (c *CourseController) GetCourseById(ctx *gin.Context) {
	slog.Debug("Getting course by ID")

	id := ctx.Param("id")
	course, err := c.service.GetCourseById(id)
	if err != nil {
		slog.Error("Error getting course by ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

// @Summary Delete a course
// @Description Delete a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} schemas.DeleteCourseResponse
// @Router /courses/{id} [delete]
func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	slog.Debug("Deleting course")
	id := ctx.Param("id")

	err := c.service.DeleteCourse(id)
	if err != nil {
		slog.Error("Error deleting course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course deleted", "id", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// @Summary Get a course by teacher ID
// @Description Get a course by teacher ID
// @Tags courses
// @Accept json
// @Produce json
// @Param teacherId path string true "Teacher ID"
// @Success 200 {array} model.Course
// @Router /courses/teacher/{teacherId} [get]
func (c *CourseController) GetCourseByTeacherId(ctx *gin.Context) {
	slog.Debug("Getting course by teacher ID")
	teacherId := ctx.Param("teacherId")
	course, err := c.service.GetCourseByTeacherId(teacherId)
	if err != nil {
		slog.Error("Error getting course by teacher ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

// @Summary Get a course by title
// @Description Get a course by title
// @Tags courses
// @Accept json
// @Produce json
// @Param title path string true "Course title"
// @Success 200 {array} model.Course
// @Router /courses/title/{title} [get]
func (c *CourseController) GetCourseByTitle(ctx *gin.Context) {
	slog.Debug("Getting course by title")
	title := ctx.Param("title")
	course, err := c.service.GetCourseByTitle(title)
	if err != nil {
		slog.Error("Error getting course by title", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course retrieved", "course", course)
	ctx.JSON(http.StatusOK, course)
}

// @Summary Update a course
// @Description Update a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param course body schemas.UpdateCourseRequest true "Course to update"
// @Success 200 {object} model.Course
// @Router /courses/{id} [put]
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	slog.Debug("Updating course")
	id := ctx.Param("id")

	var updateCourseRequest schemas.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&updateCourseRequest); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCourse, err := c.service.UpdateCourse(id, updateCourseRequest)
	if err != nil {
		slog.Error("Error updating course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Course updated", "course", updatedCourse)
	ctx.JSON(http.StatusOK, updatedCourse)
}

// @Summary Get courses by student ID
// @Description Get courses by student ID
// @Tags courses
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {array} model.Course
// @Router /courses/student/{studentId} [get]
func (c *CourseController) GetCoursesByStudentId(ctx *gin.Context) {
	slog.Debug("Getting courses by student ID")
	studentId := ctx.Param("studentId")
	courses, err := c.service.GetCoursesByStudentId(studentId)
	if err != nil {
		slog.Error("Error getting courses by student ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Courses retrieved", "courses", courses)
	ctx.JSON(http.StatusOK, courses)
}

// @Summary Get courses by user ID
// @Description Get courses by user ID
// @Tags courses
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} model.Course
// @Router /courses/user/{userId} [get]
func (c *CourseController) GetCoursesByUserId(ctx *gin.Context) {
	slog.Debug("Getting courses by user ID")
	userId := ctx.Param("userId")
	courses, err := c.service.GetCoursesByUserId(userId)
	if err != nil {
		slog.Error("Error getting courses by user ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Courses retrieved", "courses", courses)
	ctx.JSON(http.StatusOK, courses)
}

// @Summary Add an aux teacher to a course
// @Description Add an aux teacher to a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
func (c *CourseController) AddAuxTeacherToCourse(ctx *gin.Context) {
	slog.Debug("Adding aux teacher to course")
	id := ctx.Param("id")
	if id == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	var auxTeacherRequest schemas.AddAuxTeacherToCourseRequest
	if err := ctx.ShouldBindJSON(&auxTeacherRequest); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherId := auxTeacherRequest.TeacherID
	auxTeacherId := auxTeacherRequest.AuxTeacherID
	course, err := c.service.AddAuxTeacherToCourse(id, teacherId, auxTeacherId)
	if err != nil {
		slog.Error("Error adding aux teacher to course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Aux teacher added to course", "course", course)
	ctx.JSON(http.StatusOK, course)
}

// @Summary Remove an aux teacher from a course
// @Description Remove an aux teacher from a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param removeAuxTeacherRequest body schemas.RemoveAuxTeacherFromCourseRequest true "Remove aux teacher from course request"
// @Success 200 {object} model.Course
// @Router /courses/{id}/remove-aux-teacher [delete]
func (c *CourseController) RemoveAuxTeacherFromCourse(ctx *gin.Context) {
	slog.Debug("Removing aux teacher from course")
	id := ctx.Param("id")
	if id == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	var removeAuxTeacherRequest schemas.RemoveAuxTeacherFromCourseRequest
	if err := ctx.ShouldBindJSON(&removeAuxTeacherRequest); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherId := removeAuxTeacherRequest.TeacherID
	auxTeacherId := removeAuxTeacherRequest.AuxTeacherID
	course, err := c.service.RemoveAuxTeacherFromCourse(id, teacherId, auxTeacherId)
	if err != nil {
		slog.Error("Error removing aux teacher from course", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Aux teacher removed from course", "course", course)
	ctx.JSON(http.StatusOK, course)
}

// @Summary Get favourite courses
// @Description Get favourite courses by student ID
// @Tags courses
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {array} model.Course
// @Router /courses/favourite/{studentId} [get]
func (c *CourseController) GetFavouriteCourses(ctx *gin.Context) {
	slog.Debug("Getting favourite courses")
	studentId := ctx.Param("studentId")
	if studentId == "" {
		slog.Error("Student ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	courses, err := c.service.GetFavouriteCourses(studentId)
	if err != nil {
		slog.Error("Error getting favourite courses", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("Favourite courses retrieved", "courses", courses)
	ctx.JSON(http.StatusOK, courses)
}

// @Summary Create course feedback
// @Description Create course feedback by course ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param feedback body schemas.CreateCourseFeedbackRequest true "Course feedback"
// @Success 200 {object} model.CourseFeedback
// @Router /courses/{id}/feedback [post]
func (c *CourseController) CreateCourseFeedback(ctx *gin.Context) {
	slog.Debug("Creating course feedback")
	courseId := ctx.Param("id")
	if courseId == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	var feedback schemas.CreateCourseFeedbackRequest
	if err := ctx.ShouldBindJSON(&feedback); err != nil {
		slog.Error("Error binding create course feedback request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(model.FeedbackTypes, feedback.FeedbackType) {
		slog.Error("Invalid feedback type")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback type"})
		return
	}

	feedbackModel, err := c.service.CreateCourseFeedback(courseId, feedback)
	if err != nil {
		slog.Error("Error creating course feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Course feedback created", "feedback", feedbackModel)
	ctx.JSON(http.StatusOK, feedbackModel)
}

// @Summary Get course feedback
// @Description Get course feedback by course ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param getCourseFeedbackRequest body schemas.GetCourseFeedbackRequest true "Get course feedback request"
// @Success 200 {array} model.CourseFeedback
// @Router /courses/{id}/feedback [get]
func (c *CourseController) GetCourseFeedback(ctx *gin.Context) {
	slog.Debug("Getting course feedback")
	courseId := ctx.Param("id")
	if courseId == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	var getCourseFeedbackRequest schemas.GetCourseFeedbackRequest
	if err := ctx.ShouldBindJSON(&getCourseFeedbackRequest); err != nil {
		slog.Error("Error binding get course feedback request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := c.service.GetCourseFeedback(courseId, getCourseFeedbackRequest)
	if err != nil {
		slog.Error("Error getting course feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Course feedback retrieved", "feedback", feedback)
	ctx.JSON(http.StatusOK, feedback)
}

// @Summary Get course feedback summary
// @Description Get course feedback summary by course ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {string} string "Course feedback summary"
// @Router /courses/{id}/feedback/summary [get]
func (c *CourseController) GetCourseFeedbackSummary(ctx *gin.Context) {
	slog.Debug("Getting course feedback summary")
	courseId := ctx.Param("id")
	if courseId == "" {
		slog.Error("Course ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	feedbacks, err := c.service.GetCourseFeedback(courseId, schemas.GetCourseFeedbackRequest{})
	if err != nil {
		slog.Error("Error getting course feedback", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(feedbacks) == 0 {
		slog.Error("No feedbacks found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No feedbacks found"})
		return
	}

	summary, err := c.aiClient.SummarizeCourseFeedbacks(feedbacks)
	if err != nil {
		slog.Error("Error getting course feedback summary", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Course feedback summary retrieved", "summary", summary)
	ctx.JSON(http.StatusOK, summary)
}
