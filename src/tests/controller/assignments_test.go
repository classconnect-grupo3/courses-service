package controller_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/queues"
	"courses-service/src/router"
	"courses-service/src/schemas"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mockAssignmentService      = &MockAssignmentService{}
	mockAssignmentErrorService = &MockAssignmentServiceWithError{}
	normalAssignmentController = controller.NewAssignmentsController(mockAssignmentService, SharedMockNotificationsQueue, SharedMockActivityService)
	errorAssignmentController  = controller.NewAssignmentsController(mockAssignmentErrorService, SharedMockNotificationsQueue, SharedMockActivityService)
	normalAssignmentRouter     = gin.Default()
	errorAssignmentRouter      = gin.Default()
)

func init() {
	gin.SetMode(gin.TestMode)
	router.InitializeAssignmentsRoutes(normalAssignmentRouter, normalAssignmentController)
	router.InitializeAssignmentsRoutes(errorAssignmentRouter, errorAssignmentController)
}

type MockNotificationsQueue struct{}

func (m *MockNotificationsQueue) Publish(message queues.QueueMessage) error {
	return nil
}

type MockTeacherActivityService struct{}

func (m *MockTeacherActivityService) LogActivityIfAuxTeacher(courseID, teacherUUID, activityType, description string) {
	// Mock implementation - do nothing
}

func (m *MockTeacherActivityService) GetCourseActivityLogs(courseID string) ([]*model.TeacherActivityLog, error) {
	return []*model.TeacherActivityLog{}, nil
}

type MockAssignmentService struct{}

func (m *MockAssignmentService) CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error) {
	return &model.Assignment{
		ID:           primitive.NewObjectID(),
		Title:        c.Title,
		Description:  c.Description,
		Instructions: c.Instructions,
		Type:         c.Type,
		CourseID:     c.CourseID,
		DueDate:      c.DueDate,
		GracePeriod:  c.GracePeriod,
		Status:       c.Status,
		Questions:    c.Questions,
		TotalPoints:  c.TotalPoints,
		PassingScore: c.PassingScore,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *MockAssignmentService) GetAssignments() ([]*model.Assignment, error) {
	return []*model.Assignment{
		{
			ID:           primitive.NewObjectID(),
			Title:        "Test Assignment 1",
			Description:  "Test Description 1",
			Instructions: "Test Instructions 1",
			Type:         "exam",
			CourseID:     "course123",
			DueDate:      time.Now().Add(24 * time.Hour),
			GracePeriod:  30,
			Status:       "published",
			Questions: []model.Question{
				{
					ID:             "q1",
					Text:           "What is 2+2?",
					Type:           model.QuestionTypeMultipleChoice,
					Options:        []string{"3", "4", "5"},
					CorrectAnswers: []string{"4"},
					Points:         10.0,
					Order:          1,
				},
			},
			TotalPoints:  10.0,
			PassingScore: 6.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           primitive.NewObjectID(),
			Title:        "Test Assignment 2",
			Description:  "Test Description 2",
			Instructions: "Test Instructions 2",
			Type:         "homework",
			CourseID:     "course456",
			DueDate:      time.Now().Add(48 * time.Hour),
			GracePeriod:  15,
			Status:       "draft",
			Questions: []model.Question{
				{
					ID:             "q2",
					Text:           "Explain the concept of recursion",
					Type:           model.QuestionTypeText,
					Options:        []string{},
					CorrectAnswers: []string{},
					Points:         20.0,
					Order:          1,
				},
			},
			TotalPoints:  20.0,
			PassingScore: 12.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}, nil
}

func (m *MockAssignmentService) GetAssignmentById(id string) (*model.Assignment, error) {
	if id == "nonexistent" {
		return nil, nil
	}
	return &model.Assignment{
		ID:           primitive.NewObjectID(),
		Title:        "Test Assignment",
		Description:  "Test Description",
		Instructions: "Test Instructions",
		Type:         "quiz",
		CourseID:     "course123",
		DueDate:      time.Now().Add(24 * time.Hour),
		GracePeriod:  30,
		Status:       "published",
		Questions: []model.Question{
			{
				ID:             "q1",
				Text:           "What is the capital of France?",
				Type:           model.QuestionTypeMultipleChoice,
				Options:        []string{"London", "Paris", "Berlin"},
				CorrectAnswers: []string{"Paris"},
				Points:         10.0,
				Order:          1,
			},
		},
		TotalPoints:  10.0,
		PassingScore: 6.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *MockAssignmentService) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	return []*model.Assignment{
		{
			ID:           primitive.NewObjectID(),
			Title:        "Course Assignment 1",
			Description:  "Course Assignment Description 1",
			Instructions: "Course Assignment Instructions 1",
			Type:         "exam",
			CourseID:     courseId,
			DueDate:      time.Now().Add(24 * time.Hour),
			GracePeriod:  30,
			Status:       "published",
			Questions: []model.Question{
				{
					ID:             "q1",
					Text:           "Question 1",
					Type:           model.QuestionTypeText,
					Options:        []string{},
					CorrectAnswers: []string{},
					Points:         15.0,
					Order:          1,
				},
			},
			TotalPoints:  15.0,
			PassingScore: 9.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}, nil
}

func (m *MockAssignmentService) UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error) {
	return &model.Assignment{
		ID:           primitive.NewObjectID(),
		Title:        updateAssignmentRequest.Title,
		Description:  updateAssignmentRequest.Description,
		Instructions: updateAssignmentRequest.Instructions,
		Type:         updateAssignmentRequest.Type,
		CourseID:     "course123",
		DueDate:      updateAssignmentRequest.DueDate,
		GracePeriod:  updateAssignmentRequest.GracePeriod,
		Status:       updateAssignmentRequest.Status,
		Questions:    updateAssignmentRequest.Questions,
		TotalPoints:  updateAssignmentRequest.TotalPoints,
		PassingScore: updateAssignmentRequest.PassingScore,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *MockAssignmentService) DeleteAssignment(id string) error {
	return nil
}

type MockAssignmentServiceWithError struct{}

func (m *MockAssignmentServiceWithError) CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error) {
	return nil, errors.New("error creating assignment")
}

func (m *MockAssignmentServiceWithError) GetAssignments() ([]*model.Assignment, error) {
	return nil, errors.New("error getting assignments")
}

func (m *MockAssignmentServiceWithError) GetAssignmentById(id string) (*model.Assignment, error) {
	return nil, errors.New("error getting assignment by id")
}

func (m *MockAssignmentServiceWithError) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	return nil, errors.New("error getting assignments by course id")
}

func (m *MockAssignmentServiceWithError) UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error) {
	return nil, errors.New("error updating assignment")
}

func (m *MockAssignmentServiceWithError) DeleteAssignment(id string) error {
	return errors.New("error deleting assignment")
}

// Tests for GetAssignments
func TestGetAssignments(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments", nil)
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Assignment 1")
	assert.Contains(t, w.Body.String(), "Test Assignment 2")
}

func TestGetAssignmentsWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments", nil)
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting assignments")
}

// Tests for CreateAssignment
func TestCreateAssignment(t *testing.T) {
	w := httptest.NewRecorder()

	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	body := `{
		"title": "New Assignment",
		"description": "New Assignment Description", 
		"instructions": "New Assignment Instructions",
		"type": "exam",
		"course_id": "course123",
		"due_date": "` + dueDate + `",
		"grace_period": 30,
		"status": "published",
		"questions": [
			{
				"id": "q1",
				"text": "What is 2+2?",
				"type": "multiple_choice",
				"options": ["3", "4", "5"],
				"correct_answers": ["4"],
				"points": 10.0,
				"order": 1
			}
		],
		"total_points": 10.0,
		"passing_score": 6.0
	}`
	req, _ := http.NewRequest("POST", "/assignments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New Assignment")
}

func TestCreateAssignmentWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`
	req, _ := http.NewRequest("POST", "/assignments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAssignmentWithMalformedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid json`
	req, _ := http.NewRequest("POST", "/assignments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAssignmentWithError(t *testing.T) {
	w := httptest.NewRecorder()

	dueDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	body := `{
		"title": "New Assignment",
		"description": "New Assignment Description", 
		"instructions": "New Assignment Instructions",
		"type": "exam",
		"course_id": "course123",
		"due_date": "` + dueDate + `",
		"grace_period": 30,
		"status": "published",
		"questions": [
			{
				"id": "q1",
				"text": "What is 2+2?",
				"type": "multiple_choice",
				"options": ["3", "4", "5"],
				"correct_answers": ["4"],
				"points": 10.0,
				"order": 1
			}
		],
		"total_points": 10.0,
		"passing_score": 6.0
	}`
	req, _ := http.NewRequest("POST", "/assignments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error creating assignment")
}

// Tests for GetAssignmentById
func TestGetAssignmentById(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/123", nil)
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Assignment")
}

func TestGetAssignmentByIdNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/nonexistent", nil)
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "assignment not found")
}

func TestGetAssignmentByIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/123", nil)
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting assignment by id")
}

// Tests for GetAssignmentsByCourseId
func TestGetAssignmentsByCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/course/course123", nil)
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Course Assignment 1")
}

func TestGetAssignmentsByCourseIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/course/course123", nil)
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting assignments by course id")
}

// Tests for UpdateAssignment
func TestUpdateAssignment(t *testing.T) {
	w := httptest.NewRecorder()

	dueDate := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	body := `{
		"title": "Updated Assignment",
		"description": "Updated Description",
		"instructions": "Updated Instructions",
		"type": "homework",
		"due_date": "` + dueDate + `",
		"grace_period": 45,
		"status": "draft",
		"questions": [
			{
				"id": "q1",
				"text": "Updated question?",
				"type": "text",
				"options": [],
				"correct_answers": [],
				"points": 15.0,
				"order": 1
			}
		],
		"total_points": 15.0,
		"passing_score": 9.0
	}`
	req, _ := http.NewRequest("PUT", "/assignments/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Assignment")
}

func TestUpdateAssignmentWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid json`
	req, _ := http.NewRequest("PUT", "/assignments/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAssignmentWithError(t *testing.T) {
	w := httptest.NewRecorder()

	dueDate := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	body := `{
		"title": "Updated Assignment",
		"description": "Updated Description",
		"instructions": "Updated Instructions",
		"type": "homework",
		"due_date": "` + dueDate + `",
		"grace_period": 45,
		"status": "draft",
		"questions": [
			{
				"id": "q1",
				"text": "Updated question?",
				"type": "text",
				"options": [],
				"correct_answers": [],
				"points": 15.0,
				"order": 1
			}
		],
		"total_points": 15.0,
		"passing_score": 9.0
	}`
	req, _ := http.NewRequest("PUT", "/assignments/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error updating assignment")
}

// Tests for DeleteAssignment
func TestDeleteAssignment(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/assignments/123", nil)
	normalAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteAssignmentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/assignments/123", nil)
	errorAssignmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error deleting assignment")
}
