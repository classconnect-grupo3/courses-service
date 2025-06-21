package controller_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/schemas"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mockSubmissionService      = &MockSubmissionService{}
	mockSubmissionErrorService = &MockSubmissionServiceWithError{}
	normalSubmissionController = controller.NewSubmissionController(mockSubmissionService)
	errorSubmissionController  = controller.NewSubmissionController(mockSubmissionErrorService)
	normalSubmissionRouter     = gin.Default()
	errorSubmissionRouter      = gin.Default()
)

// InitializeSubmissionRoutesForTest initializes submission routes without authentication middleware for testing
func InitializeSubmissionRoutesForTest(r *gin.Engine, controller *controller.SubmissionController) {
	// Add routes without authentication middleware
	r.POST("/assignments/:assignmentId/submissions", mockStudentAuthMiddleware(), controller.CreateSubmission)
	r.GET("/assignments/:assignmentId/submissions/:id", controller.GetSubmission)
	r.PUT("/assignments/:assignmentId/submissions/:id", mockStudentAuthMiddleware(), controller.UpdateSubmission)
	r.POST("/assignments/:assignmentId/submissions/:id/submit", controller.SubmitSubmission)
	r.GET("/students/:studentUUID/submissions", mockStudentAuthMiddleware(), controller.GetSubmissionsByStudent)
	r.PUT("/assignments/:assignmentId/submissions/:id/grade", mockTeacherAuthMiddleware(), controller.GradeSubmission)
	r.GET("/assignments/:assignmentId/submissions", controller.GetSubmissionsByAssignment)
}

// mockStudentAuthMiddleware simulates student authentication middleware for testing
func mockStudentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract values from request headers or use defaults for testing
		studentUUID := c.GetHeader("Student-UUID")
		if studentUUID == "" {
			studentUUID = "student123" // default for tests
		}
		studentName := c.GetHeader("Student-Name")
		if studentName == "" {
			studentName = "Test Student" // default for tests
		}

		c.Set("student_uuid", studentUUID)
		c.Set("student_name", studentName)
		c.Next()
	}
}

// mockTeacherAuthMiddleware simulates teacher authentication middleware for testing
func mockTeacherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract values from request headers or use defaults for testing
		teacherUUID := c.GetHeader("Teacher-UUID")
		if teacherUUID == "" {
			teacherUUID = "teacher123" // default for tests
		}

		c.Set("teacher_uuid", teacherUUID)
		c.Next()
	}
}

func init() {
	gin.SetMode(gin.TestMode)
	InitializeSubmissionRoutesForTest(normalSubmissionRouter, normalSubmissionController)
	InitializeSubmissionRoutesForTest(errorSubmissionRouter, errorSubmissionController)
}

type MockSubmissionService struct{}

// AutoCorrectSubmission implements service.SubmissionServiceInterface.
func (m *MockSubmissionService) AutoCorrectSubmission(ctx context.Context, submissionID string) error {
	panic("unimplemented")
}

// GenerateFeedbackSummary implements service.SubmissionServiceInterface.
func (m *MockSubmissionService) GenerateFeedbackSummary(ctx context.Context, submissionID string) (*schemas.AiSummaryResponse, error) {
	panic("unimplemented")
}

func (m *MockSubmissionService) CreateSubmission(ctx context.Context, submission *model.Submission) error {
	submission.ID = primitive.NewObjectID()
	submission.CreatedAt = time.Now()
	submission.UpdatedAt = time.Now()
	return nil
}

func (m *MockSubmissionService) UpdateSubmission(ctx context.Context, submission *model.Submission) error {
	submission.UpdatedAt = time.Now()
	return nil
}

func (m *MockSubmissionService) SubmitSubmission(ctx context.Context, submissionID string) error {
	return nil
}

func (m *MockSubmissionService) GetSubmission(ctx context.Context, id string) (*model.Submission, error) {
	if id == "nonexistent" {
		return nil, nil
	}
	if id == "different-assignment" {
		return &model.Submission{
			ID:           primitive.NewObjectID(),
			AssignmentID: "different-assignment-id",
			StudentUUID:  "student123",
			StudentName:  "Test Student",
			Status:       model.SubmissionStatusDraft,
			Answers: []model.Answer{
				{
					QuestionID: "q1",
					Content:    "Test answer",
					Type:       "text",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return &model.Submission{
		ID:           mustParseSubmissionObjectID(id),
		AssignmentID: "assignment123",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
		Status:       model.SubmissionStatusDraft,
		Answers: []model.Answer{
			{
				QuestionID: "q1",
				Content:    "Test answer",
				Type:       "text",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockSubmissionService) GetSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return []model.Submission{
		{
			ID:           primitive.NewObjectID(),
			AssignmentID: assignmentID,
			StudentUUID:  "student123",
			StudentName:  "Test Student 1",
			Status:       model.SubmissionStatusSubmitted,
			Answers: []model.Answer{
				{
					QuestionID: "q1",
					Content:    "Answer 1",
					Type:       "text",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:           primitive.NewObjectID(),
			AssignmentID: assignmentID,
			StudentUUID:  "student456",
			StudentName:  "Test Student 2",
			Status:       model.SubmissionStatusDraft,
			Answers: []model.Answer{
				{
					QuestionID: "q1",
					Content:    "Answer 2",
					Type:       "text",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *MockSubmissionService) GetSubmissionsByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	return []model.Submission{
		{
			ID:           primitive.NewObjectID(),
			AssignmentID: "assignment123",
			StudentUUID:  studentUUID,
			StudentName:  "Test Student",
			Status:       model.SubmissionStatusSubmitted,
			Answers: []model.Answer{
				{
					QuestionID: "q1",
					Content:    "Student answer 1",
					Type:       "text",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:           primitive.NewObjectID(),
			AssignmentID: "assignment456",
			StudentUUID:  studentUUID,
			StudentName:  "Test Student",
			Status:       model.SubmissionStatusDraft,
			Answers: []model.Answer{
				{
					QuestionID: "q2",
					Content:    "Student answer 2",
					Type:       "text",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *MockSubmissionService) GetOrCreateSubmission(ctx context.Context, assignmentID, studentUUID, studentName string) (*model.Submission, error) {
	return &model.Submission{
		ID:           primitive.NewObjectID(),
		AssignmentID: assignmentID,
		StudentUUID:  studentUUID,
		StudentName:  studentName,
		Status:       model.SubmissionStatusDraft,
		Answers:      []model.Answer{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *MockSubmissionService) GradeSubmission(ctx context.Context, submissionID string, score *float64, feedback string) (*model.Submission, error) {
	return &model.Submission{
		ID:           mustParseSubmissionObjectID(submissionID),
		AssignmentID: "assignment123",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
		Status:       model.SubmissionStatusSubmitted,
		Answers: []model.Answer{
			{
				QuestionID: "q1",
				Content:    "Graded answer",
				Type:       "text",
			},
		},
		Score:     score,
		Feedback:  feedback,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockSubmissionService) ValidateTeacherPermissions(ctx context.Context, assignmentID, teacherUUID string) error {
	if teacherUUID == "unauthorized-teacher" {
		return errors.New("teacher not authorized to grade this assignment")
	}
	return nil
}

type MockSubmissionServiceWithError struct{}

// AutoCorrectSubmission implements service.SubmissionServiceInterface.
func (m *MockSubmissionServiceWithError) AutoCorrectSubmission(ctx context.Context, submissionID string) error {
	panic("unimplemented")
}

// GenerateFeedbackSummary implements service.SubmissionServiceInterface.
func (m *MockSubmissionServiceWithError) GenerateFeedbackSummary(ctx context.Context, submissionID string) (*schemas.AiSummaryResponse, error) {
	panic("unimplemented")
}

func (m *MockSubmissionServiceWithError) CreateSubmission(ctx context.Context, submission *model.Submission) error {
	return errors.New("error creating submission")
}

func (m *MockSubmissionServiceWithError) UpdateSubmission(ctx context.Context, submission *model.Submission) error {
	return errors.New("error updating submission")
}

func (m *MockSubmissionServiceWithError) SubmitSubmission(ctx context.Context, submissionID string) error {
	return errors.New("error submitting submission")
}

func (m *MockSubmissionServiceWithError) GetSubmission(ctx context.Context, id string) (*model.Submission, error) {
	return nil, errors.New("error getting submission")
}

func (m *MockSubmissionServiceWithError) GetSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return nil, errors.New("error getting submissions by assignment")
}

func (m *MockSubmissionServiceWithError) GetSubmissionsByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	return nil, errors.New("error getting submissions by student")
}

func (m *MockSubmissionServiceWithError) GetOrCreateSubmission(ctx context.Context, assignmentID, studentUUID, studentName string) (*model.Submission, error) {
	return nil, errors.New("error getting or creating submission")
}

func (m *MockSubmissionServiceWithError) GradeSubmission(ctx context.Context, submissionID string, score *float64, feedback string) (*model.Submission, error) {
	return nil, errors.New("error grading submission")
}

func (m *MockSubmissionServiceWithError) ValidateTeacherPermissions(ctx context.Context, assignmentID, teacherUUID string) error {
	if teacherUUID == "unauthorized-teacher" {
		return errors.New("teacher not authorized to grade this assignment")
	}
	return nil
}

// Helper function to create consistent ObjectIDs for testing
func mustParseSubmissionObjectID(id string) primitive.ObjectID {
	switch id {
	case "valid-submission-id":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
		return objectID
	default:
		return primitive.NewObjectID()
	}
}

// Tests for CreateSubmission
func TestCreateSubmission(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"answers": [
			{
				"question_id": "q1",
				"content": "Test answer",
				"type": "text"
			}
		]
	}`

	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")
	req.Header.Set("Student-Name", "Test Student")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "student123")
}

func TestCreateSubmissionWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid json`
	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")
	req.Header.Set("Student-Name", "Test Student")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSubmissionWithError(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"answers": [
			{
				"question_id": "q1",
				"content": "Test answer",
				"type": "text"
			}
		]
	}`

	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")
	req.Header.Set("Student-Name", "Test Student")

	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting or creating submission")
}

// Tests for GetSubmission
func TestGetSubmission(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions/valid-submission-id", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "assignment123")
}

func TestGetSubmissionNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions/nonexistent", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestGetSubmissionDifferentAssignment(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions/different-assignment", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestGetSubmissionWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions/valid-submission-id", nil)
	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting submission")
}

// Tests for UpdateSubmission
func TestUpdateSubmission(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"id": "123456789012345678901234",
		"assignment_id": "assignment123",
		"student_uuid": "student123",
		"student_name": "Test Student",
		"status": "draft",
		"answers": [
			{
				"question_id": "q1",
				"content": "Updated answer",
				"type": "text"
			}
		],
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/123456789012345678901234", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSubmissionWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid json`
	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/123456789012345678901234", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSubmissionWithIDMismatch(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"id": "999999999999999999999999",
		"assignment_id": "assignment123",
		"student_uuid": "student123",
		"student_name": "Test Student",
		"status": "draft",
		"answers": [],
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/123456789012345678901234", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "submission ID mismatch")
}

func TestUpdateSubmissionWithAssignmentMismatch(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"id": "123456789012345678901234",
		"assignment_id": "different-assignment",
		"student_uuid": "student123",
		"student_name": "Test Student",
		"status": "draft",
		"answers": [],
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/123456789012345678901234", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "assignment ID mismatch")
}

func TestUpdateSubmissionUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"id": "123456789012345678901234",
		"assignment_id": "assignment123",
		"student_uuid": "different-student",
		"student_name": "Test Student",
		"status": "draft",
		"answers": [],
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/123456789012345678901234", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

// Tests for SubmitSubmission
func TestSubmitSubmission(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions/valid-submission-id/submit", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubmitSubmissionNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions/nonexistent/submit", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestSubmitSubmissionDifferentAssignment(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions/different-assignment/submit", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestSubmitSubmissionWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/assignments/assignment123/submissions/valid-submission-id/submit", nil)
	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting submission")
}

// Tests for GetSubmissionsByAssignment
func TestGetSubmissionsByAssignment(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions", nil)
	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Student 1")
	assert.Contains(t, w.Body.String(), "Test Student 2")
}

func TestGetSubmissionsByAssignmentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assignments/assignment123/submissions", nil)
	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting submissions by assignment")
}

// Tests for GetSubmissionsByStudent
func TestGetSubmissionsByStudent(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/students/student123/submissions", nil)
	req.Header.Set("Student-UUID", "student123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "assignment123")
	assert.Contains(t, w.Body.String(), "assignment456")
}

func TestGetSubmissionsByStudentUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/students/student123/submissions", nil)
	req.Header.Set("Student-UUID", "different-student")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestGetSubmissionsByStudentWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/students/student123/submissions", nil)
	req.Header.Set("Student-UUID", "student123")

	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting submissions by student")
}

// Tests for GradeSubmission
func TestGradeSubmission(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/valid-submission-id/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "teacher123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "85.5")
	assert.Contains(t, w.Body.String(), "Great work!")
}

func TestGradeSubmissionWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid json`
	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/valid-submission-id/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "teacher123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGradeSubmissionUnauthorizedTeacher(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/valid-submission-id/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "unauthorized-teacher")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "teacher not authorized to grade this assignment")
}

func TestGradeSubmissionNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/nonexistent/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "teacher123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestGradeSubmissionDifferentAssignment(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/different-assignment/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "teacher123")

	normalSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "submission not found")
}

func TestGradeSubmissionWithInvalidTeacher(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/valid-submission-id/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "unauthorized-teacher")

	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "teacher not authorized to grade this assignment")
}

func TestGradeSubmissionWithError(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{
		"score": 85.5,
		"feedback": "Great work!"
	}`

	req, _ := http.NewRequest("PUT", "/assignments/assignment123/submissions/valid-submission-id/grade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teacher-UUID", "teacher123")

	errorSubmissionRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error getting submission")
}
