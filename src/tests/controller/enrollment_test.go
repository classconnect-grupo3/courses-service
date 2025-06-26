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
	mockActivityService        = &MockTeacherActivityService{}
	normalEnrollmentController = controller.NewEnrollmentController(mockEnrollmentService, nil, mockActivityService)
	errorEnrollmentController  = controller.NewEnrollmentController(mockErrorEnrollmentService, nil, mockActivityService)
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
				CourseID:     "course-123",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent work!",
			},
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-456",
				CourseID:     "course-456",
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

// ApproveStudent implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentService) ApproveStudent(studentID, courseID string) error {
	if studentID == "error-student" || courseID == "error-course" {
		return errors.New("error approving student")
	}
	return nil
}

// DisapproveStudent implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentService) DisapproveStudent(studentID, courseID, reason string) error {
	if studentID == "error-student" || courseID == "error-course" {
		return errors.New("error disapproving student")
	}
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

// ApproveStudent implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentServiceWithError) ApproveStudent(studentID, courseID string) error {
	return errors.New("Error approving student")
}

// DisapproveStudent implements service.EnrollmentServiceInterface.
func (m *MockEnrollmentServiceWithError) DisapproveStudent(studentID, courseID, reason string) error {
	return errors.New("Error disapproving student")
}

type MockTeacherActivityService struct{}

func (m *MockTeacherActivityService) LogActivityIfAuxTeacher(courseID, teacherUUID, activityType, description string) {
	// Mock implementation - do nothing
}

func (m *MockTeacherActivityService) GetCourseActivityLogs(courseID string) ([]*model.TeacherActivityLog, error) {
	return []*model.TeacherActivityLog{}, nil
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
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/course-123/unenroll?studentId=%s", studentId), nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Student successfully unenrolled from course")
}

func TestUnenrollStudentWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/courses/course-123/unenroll?studentId=", nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Student ID is required")
}

func TestUnenrollStudentWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses//unenroll?studentId=%s", studentId), nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnenrollStudentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/course-123/unenroll?studentId=%s", studentId), nil)
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
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/course-123/favourite?studentId=%s", studentId), nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Favourite course unset")
}

func TestUnsetFavouriteCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/courses/course-123/favourite?studentId=", nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Student ID is required")
}

func TestUnsetFavouriteCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses//favourite?studentId=%s", studentId), nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid course ID")
}

func TestUnsetFavouriteCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	studentId := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/course-123/favourite?studentId=%s", studentId), nil)
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

	req, _ := http.NewRequest("PUT", "/feedback/student/student-with-feedback", strings.NewReader(body))
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

	req, _ := http.NewRequest("PUT", "/feedback/student/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetFeedbackByStudentIdWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("PUT", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	// Should still work as the request body is optional (filters)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetFeedbackByStudentIdWithMalformedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123", "invalid_json"`

	req, _ := http.NewRequest("PUT", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestGetFeedbackByStudentIdWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123"}`

	req, _ := http.NewRequest("PUT", "/feedback/student/student-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting feedback by student ID")
}

func TestGetFeedbackByStudentIdWithNoFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"course_id": "course-123"}`

	req, _ := http.NewRequest("PUT", "/feedback/student/student-without-feedback", strings.NewReader(body))
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

			req, _ := http.NewRequest("PUT", fmt.Sprintf("/feedback/student/%s", tc.studentID), strings.NewReader(tc.body))
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

// Tests for ApproveStudent endpoint
func TestApproveStudent(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response structure
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Student approved successfully")
	assert.Contains(t, responseBody, "student-456")
	assert.Contains(t, responseBody, "course-123")
}

func TestApproveStudentWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses//students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestApproveStudentWithEmptyStudentId(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students//approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Student ID is required")
}

func TestApproveStudentWithBothEmptyIds(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses//students//approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestApproveStudentWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error approving student")
}

func TestApproveStudentWithErrorStudentId(t *testing.T) {
	w := httptest.NewRecorder()

	// Using "error-student" which triggers an error in our mock
	req, _ := http.NewRequest("PUT", "/courses/course-123/students/error-student/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error approving student")
}

func TestApproveStudentWithErrorCourseId(t *testing.T) {
	w := httptest.NewRecorder()

	// Using "error-course" which triggers an error in our mock
	req, _ := http.NewRequest("PUT", "/courses/error-course/students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error approving student")
}

func TestApproveStudentWithoutTeacherHeaders(t *testing.T) {
	w := httptest.NewRecorder()

	// No teacher headers - should fail authorization
	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-456/approve", nil)
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "X-Teacher-UUID header is required")
}

func TestApproveStudentWithEmptyTeacherUUID(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "") // Empty UUID
	req.Header.Set("X-Teacher-Name", "Test Teacher")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "X-Teacher-UUID header is required")
}

func TestApproveStudentWithOnlyTeacherUUID(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-456/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	// No X-Teacher-Name header - should still work as it's optional
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Student approved successfully")
	assert.Contains(t, responseBody, "student-456")
	assert.Contains(t, responseBody, "course-123")
}

func TestApproveStudentWithValidUUIDs(t *testing.T) {
	testCases := []struct {
		name      string
		courseId  string
		studentId string
		expected  int
	}{
		{
			name:      "Valid standard IDs",
			courseId:  "course-123",
			studentId: "student-456",
			expected:  http.StatusOK,
		},
		{
			name:      "Valid UUID format course",
			courseId:  "123e4567-e89b-12d3-a456-426614174000",
			studentId: "student-456",
			expected:  http.StatusOK,
		},
		{
			name:      "Valid UUID format student",
			courseId:  "course-123",
			studentId: "123e4567-e89b-12d3-a456-426614174000",
			expected:  http.StatusOK,
		},
		{
			name:      "Both UUID format",
			courseId:  "123e4567-e89b-12d3-a456-426614174000",
			studentId: "987fcdeb-51c2-43d4-b567-531028391849",
			expected:  http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			url := fmt.Sprintf("/courses/%s/students/%s/approve", tc.courseId, tc.studentId)
			req, _ := http.NewRequest("PUT", url, nil)
			req.Header.Set("X-Teacher-UUID", "teacher-123")
			req.Header.Set("X-Teacher-Name", "Test Teacher")
			normalEnrollmentRouter.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
			if tc.expected == http.StatusOK {
				responseBody := w.Body.String()
				assert.Contains(t, responseBody, "Student approved successfully")
				assert.Contains(t, responseBody, tc.studentId)
				assert.Contains(t, responseBody, tc.courseId)
			}
		})
	}
}

func TestApproveStudentResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/approve", nil)
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response structure
	expected := `{"message":"Student approved successfully","student_id":"student-123","course_id":"course-123"}`
	assert.JSONEq(t, expected, w.Body.String())
}

// Tests for DisapproveStudent endpoint

func TestDisapproveStudent(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Student did not meet the requirements"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Student disapproved successfully")
	assert.Contains(t, w.Body.String(), "Student did not meet the requirements")
}

func TestDisapproveStudentWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses//students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestDisapproveStudentWithEmptyStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students//disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Student ID is required")
}

func TestDisapproveStudentWithBothEmptyIds(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses//students//disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestDisapproveStudentWithInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": json}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestDisapproveStudentWithMissingReason(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestDisapproveStudentWithEmptyReason(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": ""}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestDisapproveStudentWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	errorEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error disapproving student")
}

func TestDisapproveStudentWithErrorStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/error-student/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error disapproving student")
}

func TestDisapproveStudentWithErrorCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses/error-course/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error disapproving student")
}

func TestDisapproveStudentWithLongReason(t *testing.T) {
	w := httptest.NewRecorder()
	longReason := strings.Repeat("This is a very long reason. ", 50)
	body := fmt.Sprintf(`{"reason": "%s"}`, longReason)

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Student disapproved successfully")
	assert.Contains(t, w.Body.String(), longReason)
}

func TestDisapproveStudentWithSpecialCharactersInReason(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Student failed due to: 1) Poor attendance, 2) Low grades, 3) Missed deadlines & assignments"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Student disapproved successfully")
	assert.Contains(t, w.Body.String(), "Poor attendance")
}

func TestDisapproveStudentResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Academic performance below standards"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	req.Header.Set("X-Teacher-UUID", "teacher-123")
	req.Header.Set("X-Teacher-Name", "Teacher Name")
	normalEnrollmentRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response structure
	expected := `{"message":"Student disapproved successfully","student_id":"student-123","course_id":"course-123","reason":"Academic performance below standards"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestDisapproveStudentWithoutTeacherHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"reason": "Some reason"}`

	req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
	// No teacher headers set
	normalEnrollmentRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDisapproveStudentWithDifferentReasonTypes(t *testing.T) {
	testCases := []struct {
		name     string
		reason   string
		expected string
	}{
		{
			name:     "Academic failure",
			reason:   "Failed to meet academic requirements",
			expected: "Failed to meet academic requirements",
		},
		{
			name:     "Behavioral issues",
			reason:   "Inappropriate behavior in class",
			expected: "Inappropriate behavior in class",
		},
		{
			name:     "Attendance problems",
			reason:   "Excessive absences",
			expected: "Excessive absences",
		},
		{
			name:     "Short reason",
			reason:   "Failed",
			expected: "Failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := fmt.Sprintf(`{"reason": "%s"}`, tc.reason)

			req, _ := http.NewRequest("PUT", "/courses/course-123/students/student-123/disapprove", strings.NewReader(body))
			req.Header.Set("X-Teacher-UUID", "teacher-123")
			req.Header.Set("X-Teacher-Name", "Teacher Name")
			normalEnrollmentRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), tc.expected)
		})
	}
}
