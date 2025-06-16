package service_test

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAssignmentRepository struct{}

func (m *MockAssignmentRepository) CreateAssignment(assignment model.Assignment) (*model.Assignment, error) {
	if assignment.CourseID == "error-creating-assignment" {
		return nil, errors.New("Error creating assignment")
	}

	// Simulate successful creation
	assignment.ID = primitive.NewObjectID()
	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()

	return &assignment, nil
}

func (m *MockAssignmentRepository) GetAssignments() ([]*model.Assignment, error) {
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
					Text:           "Explain recursion",
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

func (m *MockAssignmentRepository) GetByID(ctx context.Context, id string) (*model.Assignment, error) {
	if id == "valid-assignment-id" {
		return &model.Assignment{
			ID:           mustParseAssignmentObjectID(id),
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
					Text:           "Test Question",
					Type:           model.QuestionTypeText,
					Options:        []string{},
					CorrectAnswers: []string{},
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
	if id == "error-assignment-id" {
		return nil, errors.New("Error getting assignment by ID")
	}
	return nil, nil // Assignment not found
}

func (m *MockAssignmentRepository) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	if courseId == "course-with-assignments" {
		return []*model.Assignment{
			{
				ID:           primitive.NewObjectID(),
				Title:        "Course Assignment 1",
				Description:  "Assignment 1 for course",
				Instructions: "Instructions for assignment 1",
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
			{
				ID:           primitive.NewObjectID(),
				Title:        "Course Assignment 2",
				Description:  "Assignment 2 for course",
				Instructions: "Instructions for assignment 2",
				Type:         "homework",
				CourseID:     courseId,
				DueDate:      time.Now().Add(48 * time.Hour),
				GracePeriod:  15,
				Status:       "draft",
				Questions: []model.Question{
					{
						ID:             "q2",
						Text:           "Question 2",
						Type:           model.QuestionTypeMultipleChoice,
						Options:        []string{"A", "B", "C"},
						CorrectAnswers: []string{"B"},
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
	if courseId == "empty-course" {
		return []*model.Assignment{}, nil
	}
	if courseId == "error-course" {
		return nil, errors.New("Error getting assignments by course ID")
	}
	return []*model.Assignment{}, nil
}

func (m *MockAssignmentRepository) UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error) {
	if id == "valid-assignment-id" {
		updateAssignment.ID = mustParseAssignmentObjectID(id)
		updateAssignment.UpdatedAt = time.Now()
		return &updateAssignment, nil
	}
	if id == "error-updating-assignment" {
		return nil, errors.New("Error updating assignment")
	}
	return nil, errors.New("assignment not found")
}

func (m *MockAssignmentRepository) DeleteAssignment(id string) error {
	if id == "valid-assignment-id" {
		return nil
	}
	if id == "error-deleting-assignment" {
		return errors.New("Error deleting assignment")
	}
	return errors.New("assignment not found")
}

type MockCourseService struct{}

// GetCourseMembers implements service.CourseServiceInterface.
func (m *MockCourseService) GetCourseMembers(courseId string) (*schemas.CourseMembersResponse, error) {
	panic("unimplemented")
}

// GetCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseService) GetCourseFeedback(courseId string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	panic("unimplemented")
}

// CreateCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseService) CreateCourseFeedback(courseId string, feedbackRequest schemas.CreateCourseFeedbackRequest) (*model.CourseFeedback, error) {
	if courseId == "error-creating-feedback" {
		return nil, errors.New("Error creating feedback")
	}
	if courseId == "invalid-score" {
		return nil, errors.New("Score must be between 1 and 5")
	}
	return nil, nil
}

func (m *MockCourseService) GetCourseById(id string) (*model.Course, error) {
	if id == "valid-course-id" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Test Course",
			Description: "Test Course Description",
			TeacherUUID: "teacher-123",
			Capacity:    30,
		}, nil
	}
	if id == "error-course-id" {
		return nil, errors.New("Error getting course")
	}
	return nil, nil // Course not found
}

// Mock implementations for other CourseService methods (not used in assignment service but required by interface)
func (m *MockCourseService) GetCourses() ([]*model.Course, error) { return nil, nil }
func (m *MockCourseService) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) DeleteCourse(id string, teacherId string) error { return nil }
func (m *MockCourseService) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error) {
	return nil, nil
}
func (m *MockCourseService) GetCourseByTitle(title string) ([]*model.Course, error) { return nil, nil }
func (m *MockCourseService) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) AddAuxTeacherToCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseService) GetFavouriteCourses(studentId string) ([]*model.Course, error) {
	return nil, nil
}

// Helper function to create consistent ObjectIDs for testing
func mustParseAssignmentObjectID(id string) primitive.ObjectID {
	switch id {
	case "valid-assignment-id":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
		return objectID
	default:
		return primitive.NewObjectID()
	}
}

// Tests for GetAssignments
func TestGetAssignments(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignments, err := assignmentService.GetAssignments()
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 2, len(assignments))
	assert.Equal(t, "Test Assignment 1", assignments[0].Title)
	assert.Equal(t, "Test Assignment 2", assignments[1].Title)
}

// Tests for GetAssignmentById
func TestGetAssignmentById(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignment, err := assignmentService.GetAssignmentById("valid-assignment-id")
	assert.NoError(t, err)
	assert.NotNil(t, assignment)
	assert.Equal(t, "Test Assignment", assignment.Title)
	assert.Equal(t, "Test Description", assignment.Description)
}

func TestGetAssignmentByIdWithEmptyId(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignment, err := assignmentService.GetAssignmentById("")
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "id is required")
}

func TestGetAssignmentByIdWithError(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignment, err := assignmentService.GetAssignmentById("error-assignment-id")
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "Error getting assignment by ID")
}

func TestGetAssignmentByIdNotFound(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignment, err := assignmentService.GetAssignmentById("nonexistent-assignment-id")
	assert.NoError(t, err)
	assert.Nil(t, assignment)
}

// Tests for CreateAssignment
func TestCreateAssignment(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	request := schemas.CreateAssignmentRequest{
		Title:        "New Assignment",
		Description:  "New Assignment Description",
		Instructions: "New Assignment Instructions",
		Type:         "exam",
		CourseID:     "valid-course-id",
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
	}

	assignment, err := assignmentService.CreateAssignment(request)
	assert.NoError(t, err)
	assert.NotNil(t, assignment)
	assert.Equal(t, request.Title, assignment.Title)
	assert.Equal(t, request.Description, assignment.Description)
	assert.Equal(t, request.Instructions, assignment.Instructions)
	assert.Equal(t, request.Type, assignment.Type)
	assert.Equal(t, request.CourseID, assignment.CourseID)
	assert.Equal(t, request.GracePeriod, assignment.GracePeriod)
	assert.Equal(t, request.Status, assignment.Status)
	assert.Equal(t, request.TotalPoints, assignment.TotalPoints)
	assert.Equal(t, request.PassingScore, assignment.PassingScore)
	assert.False(t, assignment.ID.IsZero())
}

func TestCreateAssignmentWithCourseNotFound(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	request := schemas.CreateAssignmentRequest{
		Title:        "New Assignment",
		Description:  "New Assignment Description",
		Instructions: "New Assignment Instructions",
		Type:         "exam",
		CourseID:     "nonexistent-course-id",
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
	}

	assignment, err := assignmentService.CreateAssignment(request)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "course not found")
}

func TestCreateAssignmentWithErrorGettingCourse(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	request := schemas.CreateAssignmentRequest{
		Title:        "New Assignment",
		Description:  "New Assignment Description",
		Instructions: "New Assignment Instructions",
		Type:         "exam",
		CourseID:     "error-course-id",
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
	}

	assignment, err := assignmentService.CreateAssignment(request)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "Error getting course")
}

func TestCreateAssignmentWithErrorCreating(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	request := schemas.CreateAssignmentRequest{
		Title:        "New Assignment",
		Description:  "New Assignment Description",
		Instructions: "New Assignment Instructions",
		Type:         "exam",
		CourseID:     "error-creating-assignment",
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
	}

	assignment, err := assignmentService.CreateAssignment(request)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "course not found")
}

// Tests for UpdateAssignment
func TestUpdateAssignment(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	updateRequest := schemas.UpdateAssignmentRequest{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		DueDate:      time.Now().Add(48 * time.Hour),
		GracePeriod:  45,
		Status:       "draft",
		Questions: []model.Question{
			{
				ID:             "q1",
				Text:           "Updated question?",
				Type:           model.QuestionTypeText,
				Options:        []string{},
				CorrectAnswers: []string{},
				Points:         15.0,
				Order:          1,
			},
		},
		TotalPoints:  15.0,
		PassingScore: 9.0,
	}

	assignment, err := assignmentService.UpdateAssignment("valid-assignment-id", updateRequest)
	assert.NoError(t, err)
	assert.NotNil(t, assignment)
	assert.Equal(t, updateRequest.Title, assignment.Title)
	assert.Equal(t, updateRequest.Description, assignment.Description)
	assert.Equal(t, updateRequest.Instructions, assignment.Instructions)
	assert.Equal(t, updateRequest.Type, assignment.Type)
	assert.Equal(t, updateRequest.GracePeriod, assignment.GracePeriod)
	assert.Equal(t, updateRequest.Status, assignment.Status)
	assert.Equal(t, updateRequest.TotalPoints, assignment.TotalPoints)
	assert.Equal(t, updateRequest.PassingScore, assignment.PassingScore)
}

func TestUpdateAssignmentWithEmptyId(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	updateRequest := schemas.UpdateAssignmentRequest{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		DueDate:      time.Now().Add(48 * time.Hour),
		GracePeriod:  45,
		Status:       "draft",
		Questions:    []model.Question{},
		TotalPoints:  15.0,
		PassingScore: 9.0,
	}

	assignment, err := assignmentService.UpdateAssignment("", updateRequest)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "id is required")
}

func TestUpdateAssignmentWithAssignmentNotFound(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	updateRequest := schemas.UpdateAssignmentRequest{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		DueDate:      time.Now().Add(48 * time.Hour),
		GracePeriod:  45,
		Status:       "draft",
		Questions:    []model.Question{},
		TotalPoints:  15.0,
		PassingScore: 9.0,
	}

	assignment, err := assignmentService.UpdateAssignment("nonexistent-assignment-id", updateRequest)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "assignment not found")
}

func TestUpdateAssignmentWithErrorGettingAssignment(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	updateRequest := schemas.UpdateAssignmentRequest{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		DueDate:      time.Now().Add(48 * time.Hour),
		GracePeriod:  45,
		Status:       "draft",
		Questions:    []model.Question{},
		TotalPoints:  15.0,
		PassingScore: 9.0,
	}

	assignment, err := assignmentService.UpdateAssignment("error-assignment-id", updateRequest)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "Error getting assignment by ID")
}

func TestUpdateAssignmentWithErrorUpdating(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	updateRequest := schemas.UpdateAssignmentRequest{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		DueDate:      time.Now().Add(48 * time.Hour),
		GracePeriod:  45,
		Status:       "draft",
		Questions:    []model.Question{},
		TotalPoints:  15.0,
		PassingScore: 9.0,
	}

	assignment, err := assignmentService.UpdateAssignment("error-updating-assignment", updateRequest)
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "assignment not found")
}

// Tests for DeleteAssignment
func TestDeleteAssignment(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	err := assignmentService.DeleteAssignment("valid-assignment-id")
	assert.NoError(t, err)
}

func TestDeleteAssignmentWithEmptyId(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	err := assignmentService.DeleteAssignment("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is required")
}

func TestDeleteAssignmentWithError(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	err := assignmentService.DeleteAssignment("error-deleting-assignment")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error deleting assignment")
}

func TestDeleteAssignmentNotFound(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	err := assignmentService.DeleteAssignment("nonexistent-assignment-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "assignment not found")
}

// Tests for GetAssignmentsByCourseId
func TestGetAssignmentsByCourseId(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignments, err := assignmentService.GetAssignmentsByCourseId("course-with-assignments")
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 2, len(assignments))
	assert.Equal(t, "Course Assignment 1", assignments[0].Title)
	assert.Equal(t, "Course Assignment 2", assignments[1].Title)
}

func TestGetAssignmentsByCourseIdWithEmptyCourse(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignments, err := assignmentService.GetAssignmentsByCourseId("empty-course")
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 0, len(assignments))
}

func TestGetAssignmentsByCourseIdWithEmptyId(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignments, err := assignmentService.GetAssignmentsByCourseId("")
	assert.Error(t, err)
	assert.Nil(t, assignments)
	assert.Contains(t, err.Error(), "course id is required")
}

func TestGetAssignmentsByCourseIdWithError(t *testing.T) {
	assignmentService := service.NewAssignmentService(&MockAssignmentRepository{}, &MockCourseService{})

	assignments, err := assignmentService.GetAssignmentsByCourseId("error-course")
	assert.Error(t, err)
	assert.Nil(t, assignments)
	assert.Contains(t, err.Error(), "Error getting assignments by course ID")
}
