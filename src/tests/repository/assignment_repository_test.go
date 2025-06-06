package repository_test

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/tests/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var assignmentDBSetup *testutil.DBSetup

func init() {
	// Initialize database connection for assignment repository tests
	assignmentDBSetup = testutil.SetupTestDB()
}

func createTestAssignment() model.Assignment {
	dueDate := time.Now().Add(24 * time.Hour)
	return model.Assignment{
		Title:        "Test Assignment",
		Description:  "Test Description",
		Instructions: "Test Instructions",
		Type:         "exam",
		CourseID:     "course123",
		DueDate:      dueDate,
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
		TotalPoints:     10.0,
		PassingScore:    6.0,
		SubmissionRules: []string{"No cheating", "Submit on time"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func createTestAssignmentWithID(objectID primitive.ObjectID) model.Assignment {
	assignment := createTestAssignment()
	assignment.ID = objectID
	return assignment
}

// Tests for CreateAssignment
func TestCreateAssignment(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := createTestAssignment()

	// Test creating an assignment
	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)
	assert.NotNil(t, createdAssignment)

	// Verify the assignment was created
	assert.False(t, createdAssignment.ID.IsZero())
	assert.Equal(t, assignment.Title, createdAssignment.Title)
	assert.Equal(t, assignment.Description, createdAssignment.Description)
	assert.Equal(t, assignment.Instructions, createdAssignment.Instructions)
	assert.Equal(t, assignment.Type, createdAssignment.Type)
	assert.Equal(t, assignment.CourseID, createdAssignment.CourseID)
	assert.Equal(t, assignment.GracePeriod, createdAssignment.GracePeriod)
	assert.Equal(t, assignment.Status, createdAssignment.Status)
	assert.Equal(t, len(assignment.Questions), len(createdAssignment.Questions))
	assert.Equal(t, assignment.TotalPoints, createdAssignment.TotalPoints)
	assert.Equal(t, assignment.PassingScore, createdAssignment.PassingScore)
	assert.Equal(t, len(assignment.SubmissionRules), len(createdAssignment.SubmissionRules))
}

func TestCreateAssignmentWithMinimalData(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := model.Assignment{
		Title:       "Minimal Assignment",
		Description: "Minimal Description",
		Type:        "quiz",
		CourseID:    "course456",
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)
	assert.NotNil(t, createdAssignment)
	assert.False(t, createdAssignment.ID.IsZero())
	assert.Equal(t, assignment.Title, createdAssignment.Title)
	assert.Equal(t, assignment.Type, createdAssignment.Type)
	assert.Equal(t, assignment.CourseID, createdAssignment.CourseID)
}

// Tests for GetAssignments
func TestGetAssignments(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment1 := createTestAssignment()
	assignment1.Title = "Test Assignment 1"
	assignment1.Type = "exam"

	assignment2 := createTestAssignment()
	assignment2.Title = "Test Assignment 2"
	assignment2.Type = "homework"

	// Create test assignments
	_, err := assignmentRepository.CreateAssignment(assignment1)
	assert.NoError(t, err)
	_, err = assignmentRepository.CreateAssignment(assignment2)
	assert.NoError(t, err)

	// Get all assignments
	assignments, err := assignmentRepository.GetAssignments()
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 2, len(assignments))

	// Verify assignments are returned correctly
	titles := []string{assignments[0].Title, assignments[1].Title}
	assert.Contains(t, titles, "Test Assignment 1")
	assert.Contains(t, titles, "Test Assignment 2")
}

func TestGetAssignmentsEmpty(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignments, err := assignmentRepository.GetAssignments()
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 0, len(assignments))
}

// Tests for GetByID
func TestGetByID(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := createTestAssignment()

	// Create assignment
	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)

	// Get assignment by ID
	gotAssignment, err := assignmentRepository.GetByID(context.TODO(), createdAssignment.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, gotAssignment)

	// Verify assignment details
	assert.Equal(t, createdAssignment.ID, gotAssignment.ID)
	assert.Equal(t, createdAssignment.Title, gotAssignment.Title)
	assert.Equal(t, createdAssignment.Description, gotAssignment.Description)
	assert.Equal(t, createdAssignment.Instructions, gotAssignment.Instructions)
	assert.Equal(t, createdAssignment.Type, gotAssignment.Type)
	assert.Equal(t, createdAssignment.CourseID, gotAssignment.CourseID)
	assert.Equal(t, createdAssignment.GracePeriod, gotAssignment.GracePeriod)
	assert.Equal(t, createdAssignment.Status, gotAssignment.Status)
	assert.Equal(t, len(createdAssignment.Questions), len(gotAssignment.Questions))
	assert.Equal(t, createdAssignment.TotalPoints, gotAssignment.TotalPoints)
	assert.Equal(t, createdAssignment.PassingScore, gotAssignment.PassingScore)
}

func TestGetByIDNotFound(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	// Try to get non-existent assignment
	assignment, err := assignmentRepository.GetByID(context.TODO(), "663463666666666666666666")
	assert.NoError(t, err)
	assert.Nil(t, assignment)
}

func TestGetByIDWithInvalidID(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	// Try to get assignment with invalid ID
	assignment, err := assignmentRepository.GetByID(context.TODO(), "invalid-id")
	assert.Error(t, err)
	assert.Nil(t, assignment)
	assert.Contains(t, err.Error(), "failed to get assignment by id")
}

// Tests for GetAssignmentsByCourseId
func TestGetAssignmentsByCourseId(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	courseID := "course123"

	assignment1 := createTestAssignment()
	assignment1.Title = "Course Assignment 1"
	assignment1.CourseID = courseID

	assignment2 := createTestAssignment()
	assignment2.Title = "Course Assignment 2"
	assignment2.CourseID = courseID

	// Assignment for different course
	assignment3 := createTestAssignment()
	assignment3.Title = "Other Course Assignment"
	assignment3.CourseID = "othercourse456"

	// Create assignments
	_, err := assignmentRepository.CreateAssignment(assignment1)
	assert.NoError(t, err)
	_, err = assignmentRepository.CreateAssignment(assignment2)
	assert.NoError(t, err)
	_, err = assignmentRepository.CreateAssignment(assignment3)
	assert.NoError(t, err)

	// Get assignments for specific course
	assignments, err := assignmentRepository.GetAssignmentsByCourseId(courseID)
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 2, len(assignments))

	// Verify only assignments for the specified course are returned
	for _, assignment := range assignments {
		assert.Equal(t, courseID, assignment.CourseID)
	}

	titles := []string{assignments[0].Title, assignments[1].Title}
	assert.Contains(t, titles, "Course Assignment 1")
	assert.Contains(t, titles, "Course Assignment 2")
}

func TestGetAssignmentsByCourseIdEmpty(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignments, err := assignmentRepository.GetAssignmentsByCourseId("nonexistent-course")
	assert.NoError(t, err)
	assert.NotNil(t, assignments)
	assert.Equal(t, 0, len(assignments))
}

// Tests for UpdateAssignment
func TestUpdateAssignment(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := createTestAssignment()

	// Create assignment
	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)

	// Update assignment
	updateAssignment := model.Assignment{
		Title:        "Updated Assignment",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
		Type:         "homework",
		Status:       "draft",
		GracePeriod:  45,
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
		TotalPoints:     15.0,
		PassingScore:    9.0,
		SubmissionRules: []string{"Updated rule"},
	}

	updatedAssignment, err := assignmentRepository.UpdateAssignment(createdAssignment.ID.Hex(), updateAssignment)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAssignment)

	// Verify assignment was updated
	assert.Equal(t, createdAssignment.ID, updatedAssignment.ID)
	assert.Equal(t, updateAssignment.Title, updatedAssignment.Title)
	assert.Equal(t, updateAssignment.Description, updatedAssignment.Description)
	assert.Equal(t, updateAssignment.Instructions, updatedAssignment.Instructions)
	assert.Equal(t, updateAssignment.Type, updatedAssignment.Type)
	assert.Equal(t, updateAssignment.Status, updatedAssignment.Status)
	assert.Equal(t, updateAssignment.GracePeriod, updatedAssignment.GracePeriod)
	assert.Equal(t, updateAssignment.TotalPoints, updatedAssignment.TotalPoints)
	assert.Equal(t, updateAssignment.PassingScore, updatedAssignment.PassingScore)
	assert.Equal(t, len(updateAssignment.Questions), len(updatedAssignment.Questions))
	assert.Equal(t, len(updateAssignment.SubmissionRules), len(updatedAssignment.SubmissionRules))

	// Verify the course ID is preserved (not updated)
	assert.Equal(t, assignment.CourseID, updatedAssignment.CourseID)
}

func TestUpdateAssignmentPartial(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := createTestAssignment()

	// Create assignment
	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)

	// Update only title and description
	updateAssignment := model.Assignment{
		Title:       "Partially Updated Assignment",
		Description: "Partially Updated Description",
		// Other fields left empty to test partial update
	}

	updatedAssignment, err := assignmentRepository.UpdateAssignment(createdAssignment.ID.Hex(), updateAssignment)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAssignment)

	// Verify only specified fields were updated
	assert.Equal(t, updateAssignment.Title, updatedAssignment.Title)
	assert.Equal(t, updateAssignment.Description, updatedAssignment.Description)

	// Verify other fields remained unchanged
	assert.Equal(t, assignment.Type, updatedAssignment.Type)
	assert.Equal(t, assignment.CourseID, updatedAssignment.CourseID)
	assert.Equal(t, assignment.Status, updatedAssignment.Status)
	assert.Equal(t, assignment.GracePeriod, updatedAssignment.GracePeriod)
}

func TestUpdateAssignmentWithInvalidID(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	updateAssignment := model.Assignment{
		Title: "Updated Assignment",
	}

	// Try to update assignment with invalid ID
	updatedAssignment, err := assignmentRepository.UpdateAssignment("invalid-id", updateAssignment)
	assert.Error(t, err)
	assert.Nil(t, updatedAssignment)
	assert.Contains(t, err.Error(), "failed to update assignment")
}

func TestUpdateAssignmentNotFound(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	updateAssignment := model.Assignment{
		Title: "Updated Assignment",
	}

	// Try to update non-existent assignment
	updatedAssignment, err := assignmentRepository.UpdateAssignment("663463666666666666666666", updateAssignment)
	assert.NoError(t, err)           // Update operation succeeds even if no document is matched
	assert.Nil(t, updatedAssignment) // But returns nil because GetByID returns nil
}

// Tests for DeleteAssignment
func TestDeleteAssignment(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	assignment := createTestAssignment()

	// Create assignment
	createdAssignment, err := assignmentRepository.CreateAssignment(assignment)
	assert.NoError(t, err)

	// Delete assignment
	err = assignmentRepository.DeleteAssignment(createdAssignment.ID.Hex())
	assert.NoError(t, err)

	// Verify assignment was deleted
	deletedAssignment, err := assignmentRepository.GetByID(context.TODO(), createdAssignment.ID.Hex())
	assert.NoError(t, err)
	assert.Nil(t, deletedAssignment)
}

func TestDeleteAssignmentWithInvalidID(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	// Try to delete assignment with invalid ID
	err := assignmentRepository.DeleteAssignment("invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete assignment")
}

func TestDeleteAssignmentNotFound(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	// Try to delete non-existent assignment
	err := assignmentRepository.DeleteAssignment("663463666666666666666666")
	assert.NoError(t, err) // Delete operation succeeds even if no document is matched
}

// Integration tests
func TestCompleteAssignmentWorkflow(t *testing.T) {
	t.Cleanup(func() {
		assignmentDBSetup.CleanupCollection("assignments")
	})

	assignmentRepository := repository.NewAssignmentRepository(assignmentDBSetup.Client, assignmentDBSetup.DBName)

	// Create multiple assignments for the same course
	courseID := "integration-course"

	assignment1 := createTestAssignment()
	assignment1.Title = "Integration Test Assignment 1"
	assignment1.CourseID = courseID
	assignment1.Type = "exam"

	assignment2 := createTestAssignment()
	assignment2.Title = "Integration Test Assignment 2"
	assignment2.CourseID = courseID
	assignment2.Type = "homework"

	// Create assignments
	createdAssignment1, err := assignmentRepository.CreateAssignment(assignment1)
	assert.NoError(t, err)
	createdAssignment2, err := assignmentRepository.CreateAssignment(assignment2)
	assert.NoError(t, err)

	// Get all assignments
	allAssignments, err := assignmentRepository.GetAssignments()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allAssignments))

	// Get assignments by course ID
	courseAssignments, err := assignmentRepository.GetAssignmentsByCourseId(courseID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(courseAssignments))

	// Update one assignment
	updateAssignment := model.Assignment{
		Title:       "Updated Integration Assignment",
		Description: "Updated for integration test",
		Status:      "draft",
	}

	updatedAssignment, err := assignmentRepository.UpdateAssignment(createdAssignment1.ID.Hex(), updateAssignment)
	assert.NoError(t, err)
	assert.Equal(t, updateAssignment.Title, updatedAssignment.Title)
	assert.Equal(t, updateAssignment.Status, updatedAssignment.Status)

	// Delete one assignment
	err = assignmentRepository.DeleteAssignment(createdAssignment2.ID.Hex())
	assert.NoError(t, err)

	// Verify only one assignment remains
	remainingAssignments, err := assignmentRepository.GetAssignmentsByCourseId(courseID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(remainingAssignments))
	assert.Equal(t, updatedAssignment.Title, remainingAssignments[0].Title)
}
