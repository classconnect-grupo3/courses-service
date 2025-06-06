package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmissionMockRepository struct{}

func (m *SubmissionMockRepository) Create(ctx context.Context, submission *model.Submission) error {
	submission.ID = primitive.NewObjectID()
	return nil
}

func (m *SubmissionMockRepository) Update(ctx context.Context, submission *model.Submission) error {
	return nil
}

func (m *SubmissionMockRepository) GetByID(ctx context.Context, id string) (*model.Submission, error) {
	if id == "valid-submission-id" {
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
	if id == "nonexistent" {
		return nil, nil
	}
	return nil, errors.New("repository error")
}

func (m *SubmissionMockRepository) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error) {
	if assignmentID == "existing-assignment" && studentUUID == "existing-student" {
		return &model.Submission{
			ID:           primitive.NewObjectID(),
			AssignmentID: assignmentID,
			StudentUUID:  studentUUID,
			StudentName:  "Existing Student",
			Status:       model.SubmissionStatusDraft,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}, nil
	}
	if assignmentID == "new-assignment" && studentUUID == "new-student" {
		return nil, nil // No existing submission
	}
	return nil, errors.New("repository error")
}

func (m *SubmissionMockRepository) GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	if assignmentID == "assignment123" {
		return []model.Submission{
			{
				ID:           primitive.NewObjectID(),
				AssignmentID: assignmentID,
				StudentUUID:  "student1",
				StudentName:  "Student 1",
				Status:       model.SubmissionStatusSubmitted,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}, nil
	}
	return nil, errors.New("repository error")
}

func (m *SubmissionMockRepository) GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	if studentUUID == "student123" {
		return []model.Submission{
			{
				ID:           primitive.NewObjectID(),
				AssignmentID: "assignment1",
				StudentUUID:  studentUUID,
				StudentName:  "Test Student",
				Status:       model.SubmissionStatusSubmitted,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}, nil
	}
	return nil, errors.New("repository error")
}

type SubmissionMockRepositoryWithError struct{}

func (m *SubmissionMockRepositoryWithError) Create(ctx context.Context, submission *model.Submission) error {
	return errors.New("repository create error")
}

func (m *SubmissionMockRepositoryWithError) Update(ctx context.Context, submission *model.Submission) error {
	return errors.New("repository update error")
}

func (m *SubmissionMockRepositoryWithError) GetByID(ctx context.Context, id string) (*model.Submission, error) {
	return nil, errors.New("repository get error")
}

func (m *SubmissionMockRepositoryWithError) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error) {
	return nil, errors.New("repository get error")
}

func (m *SubmissionMockRepositoryWithError) GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return nil, errors.New("repository get error")
}

func (m *SubmissionMockRepositoryWithError) GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	return nil, errors.New("repository get error")
}

type AssignmentMockRepository struct{}

func (m *AssignmentMockRepository) GetByID(ctx context.Context, id string) (*model.Assignment, error) {
	if id == "assignment123" {
		dueDate := time.Now().Add(24 * time.Hour) // Due in 24 hours
		return &model.Assignment{
			ID:          mustParseSubmissionObjectID("assignment123"),
			Title:       "Test Assignment",
			CourseID:    "course123",
			DueDate:     dueDate,
			GracePeriod: 30, // 30 minutes grace period
			Status:      "published",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if id == "nonexistent-assignment" {
		return nil, nil
	}
	return nil, errors.New("repository error")
}

func (m *AssignmentMockRepository) CreateAssignment(assignment model.Assignment) (*model.Assignment, error) {
	return nil, nil
}

func (m *AssignmentMockRepository) GetAssignments() ([]*model.Assignment, error) {
	return nil, nil
}

func (m *AssignmentMockRepository) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	return nil, nil
}

func (m *AssignmentMockRepository) UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error) {
	return nil, nil
}

func (m *AssignmentMockRepository) DeleteAssignment(id string) error {
	return nil
}

type CourseMockService struct{}

func (m *CourseMockService) GetCourseById(id string) (*model.Course, error) {
	if id == "course123" {
		return &model.Course{
			ID:          mustParseSubmissionObjectID("course123"),
			Title:       "Test Course",
			TeacherUUID: "teacher123",
			AuxTeachers: []string{"aux-teacher1", "aux-teacher2"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if id == "nonexistent-course" {
		return nil, nil
	}
	return nil, errors.New("course service error")
}

func (m *CourseMockService) GetCourses() ([]*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) DeleteCourse(id string) error {
	return nil
}

func (m *CourseMockService) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error) {
	return nil, nil
}

func (m *CourseMockService) GetCourseByTitle(title string) ([]*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) AddAuxTeacherToCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}

func (m *CourseMockService) GetFavouriteCourses(studentId string) ([]*model.Course, error) {
	return nil, nil
}

// Helper function to create consistent ObjectIDs for testing
func mustParseSubmissionObjectID(id string) primitive.ObjectID {
	switch id {
	case "assignment123":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
		return objectID
	case "course123":
		objectID, _ := primitive.ObjectIDFromHex("345678901234567890123456")
		return objectID
	case "valid-submission-id":
		objectID, _ := primitive.ObjectIDFromHex("456789012345678901234567")
		return objectID
	default:
		return primitive.NewObjectID()
	}
}

// Tests for CreateSubmission
func TestCreateSubmission(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission := &model.Submission{
		AssignmentID: "assignment123",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
		Answers: []model.Answer{
			{
				QuestionID: "q1",
				Content:    "Test answer",
				Type:       "text",
			},
		},
	}

	err := submissionService.CreateSubmission(context.TODO(), submission)
	assert.NoError(t, err)
	assert.False(t, submission.ID.IsZero())
	assert.Equal(t, model.SubmissionStatusDraft, submission.Status)
	assert.False(t, submission.CreatedAt.IsZero())
	assert.False(t, submission.UpdatedAt.IsZero())
}

func TestCreateSubmissionWithNonexistentAssignment(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission := &model.Submission{
		AssignmentID: "nonexistent-assignment",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
	}

	err := submissionService.CreateSubmission(context.TODO(), submission)
	assert.Error(t, err)
	assert.Equal(t, service.ErrAssignmentNotFound, err)
}

func TestCreateSubmissionWithRepositoryError(t *testing.T) {
	submissionRepo := &SubmissionMockRepositoryWithError{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission := &model.Submission{
		AssignmentID: "assignment123",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
	}

	err := submissionService.CreateSubmission(context.TODO(), submission)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository create error")
}

// Tests for GetSubmission
func TestGetSubmission(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission, err := submissionService.GetSubmission(context.TODO(), "valid-submission-id")
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.Equal(t, "assignment123", submission.AssignmentID)
	assert.Equal(t, "student123", submission.StudentUUID)
}

func TestGetSubmissionWithNonexistentID(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission, err := submissionService.GetSubmission(context.TODO(), "nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, submission)
}

// Tests for GetOrCreateSubmission
func TestGetOrCreateSubmissionExisting(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission, err := submissionService.GetOrCreateSubmission(context.TODO(), "existing-assignment", "existing-student", "Existing Student")
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.Equal(t, "existing-assignment", submission.AssignmentID)
	assert.Equal(t, "existing-student", submission.StudentUUID)
	assert.Equal(t, "Existing Student", submission.StudentName)
}

func TestGetOrCreateSubmissionNew(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	submission, err := submissionService.GetOrCreateSubmission(context.TODO(), "new-assignment", "new-student", "New Student")
	assert.NoError(t, err)
	assert.NotNil(t, submission)
	assert.Equal(t, "new-assignment", submission.AssignmentID)
	assert.Equal(t, "new-student", submission.StudentUUID)
	assert.Equal(t, "New Student", submission.StudentName)
	assert.Equal(t, model.SubmissionStatusDraft, submission.Status)
	assert.False(t, submission.ID.IsZero())
}

// Tests for GradeSubmission
func TestGradeSubmission(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	score := 85.5
	feedback := "Great work!"

	gradedSubmission, err := submissionService.GradeSubmission(context.TODO(), "valid-submission-id", &score, feedback)
	assert.NoError(t, err)
	assert.NotNil(t, gradedSubmission)
	assert.Equal(t, &score, gradedSubmission.Score)
	assert.Equal(t, feedback, gradedSubmission.Feedback)
	assert.False(t, gradedSubmission.UpdatedAt.IsZero())
}

func TestGradeSubmissionWithNonexistentSubmission(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	score := 85.5
	feedback := "Great work!"

	gradedSubmission, err := submissionService.GradeSubmission(context.TODO(), "nonexistent", &score, feedback)
	assert.Error(t, err)
	assert.Nil(t, gradedSubmission)
	assert.Equal(t, service.ErrSubmissionNotFound, err)
}

// Tests for ValidateTeacherPermissions
func TestValidateTeacherPermissionsMainTeacher(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	err := submissionService.ValidateTeacherPermissions(context.TODO(), "assignment123", "teacher123")
	assert.NoError(t, err)
}

func TestValidateTeacherPermissionsAuxTeacher(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	err := submissionService.ValidateTeacherPermissions(context.TODO(), "assignment123", "aux-teacher1")
	assert.NoError(t, err)
}

func TestValidateTeacherPermissionsUnauthorized(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	err := submissionService.ValidateTeacherPermissions(context.TODO(), "assignment123", "unauthorized-teacher")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher not authorized to grade this assignment")
}

func TestValidateTeacherPermissionsWithNonexistentAssignment(t *testing.T) {
	submissionRepo := &SubmissionMockRepository{}
	assignmentRepo := &AssignmentMockRepository{}
	courseService := &CourseMockService{}
	submissionService := service.NewSubmissionService(submissionRepo, assignmentRepo, courseService)

	err := submissionService.ValidateTeacherPermissions(context.TODO(), "nonexistent-assignment", "teacher123")
	assert.Error(t, err)
	assert.Equal(t, service.ErrAssignmentNotFound, err)
}
