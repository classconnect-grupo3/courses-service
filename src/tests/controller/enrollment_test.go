package controller_test

import (
	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"errors"
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

type MockEnrollmentServiceWithError struct{}

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
