package controller_test

import (
	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"courses-service/src/schemas"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	mockEnrollmentService      = &MockEnrollmentService{}
	mockErrorEnrollmentService = &MockEnrollmentServiceWithError{}
	normalEnrollmentController = controller.NewEnrollmentController(mockEnrollmentService)
	errorEnrollmentController  = controller.NewEnrollmentController(mockErrorEnrollmentService)
	normalEnrollmentRouter     = gin.Default()
	errorEnrollmentRouter      = gin.Default()
)

func init() {
	router.InitializeEnrollmentsRoutes(normalEnrollmentRouter, normalEnrollmentController)
	router.InitializeEnrollmentsRoutes(errorEnrollmentRouter, errorEnrollmentController)
}

type MockEnrollmentService struct{}

// CreateStudentFeedback implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentService) CreateStudentFeedback(feedbackRequest schemas.CreateStudentFeedbackRequest) error {
	return nil
}

// GetEnrollmentsByCourseId implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentService) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	return nil, nil
}

func (m *MockEnrollmentService) EnrollStudent(studentID, courseID string) error {
	return nil
}

func (m *MockEnrollmentService) UnenrollStudent(studentID, courseID string) error {
	return nil
}

func (m *MockEnrollmentService) SetFavouriteCourse(studentID, courseID string) error {
	return nil
}

func (m *MockEnrollmentService) UnsetFavouriteCourse(studentID, courseID string) error {
	return nil
}

type MockEnrollmentServiceWithError struct{}

// CreateStudentFeedback implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentServiceWithError) CreateStudentFeedback(feedbackRequest schemas.CreateStudentFeedbackRequest) error {
	return errors.New("Error creating student feedback")
}

// GetEnrollmentsByCourseId implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentServiceWithError) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	return nil, errors.New("Error getting enrollments by course ID")
}

func (m *MockEnrollmentServiceWithError) EnrollStudent(studentID, courseID string) error {
	return errors.New("Error enrolling student")
}

func (m *MockEnrollmentServiceWithError) UnenrollStudent(studentID, courseID string) error {
	return errors.New("Error unenrolling student")
}

func (m *MockEnrollmentServiceWithError) SetFavouriteCourse(studentID, courseID string) error {
	return errors.New("Error setting favourite course")
}

func (m *MockEnrollmentServiceWithError) UnsetFavouriteCourse(studentID, courseID string) error {
	return errors.New("Error unsetting favourite course")
}

func TestEnrollStudent(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/enroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Student successfully enrolled in course")
}

func TestEnrollStudentWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/enroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestEnrollStudentWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses//enroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEnrollStudentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/enroll", strings.NewReader(body))
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error enrolling student")
}

func TestUnenrollStudent(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/unenroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Student successfully unenrolled from course")
}

func TestUnenrollStudentWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/unenroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestUnenrollStudentWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses//unenroll", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnenrollStudentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/unenroll", strings.NewReader(body))
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error unenrolling student")
}

func TestGetEnrollmentsByCourseId(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/courses/course-123/enrollments", nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetEnrollmentsByCourseIdWithError(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/courses/course-123/enrollments", nil)
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting enrollments by course ID")
}

func TestSetFavouriteCourse(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Favourite course set")
}

func TestSetFavouriteCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestSetFavouriteCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses//favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid course ID")
}

func TestSetFavouriteCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/favourite", strings.NewReader(body))
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error setting favourite course")
}

func TestUnsetFavouriteCourse(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Favourite course unset")
}

func TestUnsetFavouriteCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestUnsetFavouriteCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses//favourite", strings.NewReader(body))
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid course ID")
}

func TestUnsetFavouriteCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_id": "123e4567-e89b-12d3-a456-426614174000"}`

	req, _ := http.NewRequest("DELETE", "/courses/course-123/favourite", strings.NewReader(body))
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error unsetting favourite course")
}

func TestCreateFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"teacher_uuid": "teacher-456",
		"course_id": "course-123",
		"feedback_type": "POSITIVO",
		"score": 4,
		"feedback": "Excellent work on the assignment!"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Feedback created")
}

func TestCreateFeedbackWithEmptyCourseID(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"teacher_uuid": "teacher-456",
		"course_id": "course-123",
		"feedback_type": "POSITIVO",
		"score": 4,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("POST", "/courses//feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateFeedbackWithInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"teacher_uuid": "teacher-456",
		"invalid_field": "invalid"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestCreateFeedbackWithMissingRequiredFields(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "",
		"teacher_uuid": "teacher-456",
		"course_id": "course-123",
		"feedback_type": "POSITIVO",
		"score": 4,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestCreateFeedbackWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"teacher_uuid": "teacher-456",
		"course_id": "course-123",
		"feedback_type": "POSITIVO",
		"score": 4,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error creating student feedback")
}

func TestCreateFeedbackWithDifferentFeedbackTypes(t *testing.T) {
	testCases := []struct {
		feedbackType string
		score        int
		expected     string
	}{
		{"POSITIVO", 5, "Feedback created"},
		{"NEUTRO", 3, "Feedback created"},
		{"NEGATIVO", 1, "Feedback created"},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		body := fmt.Sprintf(`{
			"student_uuid": "student-123",
			"teacher_uuid": "teacher-456",
			"course_id": "course-123",
			"feedback_type": "%s",
			"score": %d,
			"feedback": "Test feedback"
		}`, tc.feedbackType, tc.score)

		req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		normalEnrollmentRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), tc.expected)
	}
}

func TestCreateFeedbackWithInvalidFeedbackType(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"teacher_uuid": "teacher-456",
		"course_id": "course-123",
		"feedback_type": "INVALID_TYPE",
		"score": 3,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid feedback type")
}

func TestCreateFeedbackWithValidScoreBoundaries(t *testing.T) {
	testCases := []struct {
		score    int
		expected int
	}{
		{1, http.StatusOK}, // Lower boundary
		{5, http.StatusOK}, // Upper boundary
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		body := fmt.Sprintf(`{
			"student_uuid": "student-123",
			"teacher_uuid": "teacher-456",
			"course_id": "course-123",
			"feedback_type": "POSITIVO",
			"score": %d,
			"feedback": "Score boundary test"
		}`, tc.score)

		req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		normalEnrollmentRouter.ServeHTTP(w, req)

		assert.Equal(t, tc.expected, w.Code)
		if tc.expected == http.StatusOK {
			assert.Contains(t, w.Body.String(), "Feedback created")
		}
	}
}
