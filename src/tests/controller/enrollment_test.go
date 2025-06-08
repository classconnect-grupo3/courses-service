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

// GetFeedbackByStudentId implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentService) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	if studentID == "student-with-feedback" {
		return []*model.StudentFeedback{
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-123",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent work!",
			},
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-456",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Good effort",
			},
		}, nil
	}
	return []*model.StudentFeedback{}, nil
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

// GetFeedbackByStudentId implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentServiceWithError) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	return nil, errors.New("Error getting feedback by student ID")
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

	req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("POST", "/courses//student-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

		req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
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

		req, _ := http.NewRequest("POST", "/courses/course-123/student-feedback", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		normalEnrollmentRouter.ServeHTTP(w, req)

		assert.Equal(t, tc.expected, w.Code)
		if tc.expected == http.StatusOK {
			assert.Contains(t, w.Body.String(), "Feedback created")
		}
	}
}

func TestGetFeedbackByStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123", "feedback_type": "POSITIVO"}`

	req, _ := http.NewRequest("GET", "/feedback/student/student-with-feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Verify that response contains feedback data
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-with-feedback")
	assert.Contains(t, responseBody, "teacher-123")
	assert.Contains(t, responseBody, "Excellent work!")
}

func TestGetFeedbackByStudentIdWithEmptyStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123"}`

	req, _ := http.NewRequest("GET", "/feedback/student/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetFeedbackByStudentIdWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("GET", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	// Should still work as the request body is optional (filters)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetFeedbackByStudentIdWithMalformedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123", "invalid_json"`

	req, _ := http.NewRequest("GET", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestGetFeedbackByStudentIdWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123"}`

	req, _ := http.NewRequest("GET", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting feedback by student ID")
}

func TestGetFeedbackByStudentIdWithNoFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123"}`

	req, _ := http.NewRequest("GET", "/feedback/student/student-without-feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String())) // Empty array response
}

func TestGetFeedbackByStudentIdWithDifferentFilters(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		expectedCode int
		studentID    string
	}{
		{
			name:         "Filter by course ID",
			body:         `{"course_id": "specific-course"}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
		{
			name:         "Filter by feedback type",
			body:         `{"feedback_type": "POSITIVO"}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
		{
			name:         "Filter by score range",
			body:         `{"start_score": 4, "end_score": 5}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
		{
			name:         "Filter by date range",
			body:         `{"start_date": "2024-01-01T00:00:00Z", "end_date": "2024-12-31T23:59:59Z"}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
		{
			name:         "Combined filters",
			body:         `{"course_id": "course-123", "feedback_type": "POSITIVO", "start_score": 4}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
		{
			name:         "Empty filters",
			body:         `{}`,
			expectedCode: http.StatusOK,
			studentID:    "student-with-feedback",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/feedback/student/%s", tc.studentID), strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			normalEnrollmentRouter.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedCode == http.StatusOK {
				responseBody := w.Body.String()
				// For student-with-feedback, we should get feedback data
				if tc.studentID == "student-with-feedback" {
					assert.Contains(t, responseBody, tc.studentID)
				}
			}
		})
	}
}
