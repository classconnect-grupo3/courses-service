package repository_test

import (
	"context"
	"testing"
	"time"

	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/tests/testutil"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var submissionDBSetup *testutil.DBSetup

func init() {
	// Initialize database connection for submission repository tests
	submissionDBSetup = testutil.SetupTestDB()
}

func createTestSubmission() model.Submission {
	return model.Submission{
		AssignmentID: "assignment123",
		StudentUUID:  "student123",
		StudentName:  "Test Student",
		Status:       model.SubmissionStatusDraft,
		Answers: []model.Answer{
			{
				QuestionID: "q1",
				Content:    "The answer is 4",
				Type:       "text",
			},
			{
				QuestionID: "q2",
				Content:    []string{"option1", "option3"},
				Type:       "multiple_choice",
			},
		},
		Score:     nil,
		Feedback:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestSubmissionWithID(objectID primitive.ObjectID) model.Submission {
	submission := createTestSubmission()
	submission.ID = objectID
	return submission
}

func createTestSubmissionWithDetails(assignmentID, studentUUID, studentName string, status model.SubmissionStatus) model.Submission {
	submission := createTestSubmission()
	submission.AssignmentID = assignmentID
	submission.StudentUUID = studentUUID
	submission.StudentName = studentName
	submission.Status = status
	return submission
}

// Tests for Create
func TestCreateSubmission(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := createTestSubmission()

	// Test creating a submission
	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)

	// Verify the submission was created
	assert.False(t, submission.ID.IsZero())
	assert.Equal(t, "assignment123", submission.AssignmentID)
	assert.Equal(t, "student123", submission.StudentUUID)
	assert.Equal(t, "Test Student", submission.StudentName)
	assert.Equal(t, model.SubmissionStatusDraft, submission.Status)
	assert.Equal(t, 2, len(submission.Answers))
	assert.Equal(t, "q1", submission.Answers[0].QuestionID)
	assert.Equal(t, "The answer is 4", submission.Answers[0].Content)
	assert.Equal(t, "text", submission.Answers[0].Type)
}

func TestCreateSubmissionWithMinimalData(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := model.Submission{
		AssignmentID: "assignment456",
		StudentUUID:  "student456",
		StudentName:  "Minimal Student",
		Status:       model.SubmissionStatusDraft,
		Answers:      []model.Answer{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)
	assert.False(t, submission.ID.IsZero())
	assert.Equal(t, "assignment456", submission.AssignmentID)
	assert.Equal(t, "student456", submission.StudentUUID)
	assert.Equal(t, "Minimal Student", submission.StudentName)
	assert.Equal(t, model.SubmissionStatusDraft, submission.Status)
	assert.Equal(t, 0, len(submission.Answers))
}

// Tests for Update
func TestUpdateSubmission(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := createTestSubmission()

	// Create submission
	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)

	// Update submission
	submission.Status = model.SubmissionStatusSubmitted
	submission.Answers = append(submission.Answers, model.Answer{
		QuestionID: "q3",
		Content:    "New answer",
		Type:       "text",
	})
	score := 85.5
	submission.Score = &score
	submission.Feedback = "Great work!"
	submittedAt := time.Now()
	submission.SubmittedAt = &submittedAt
	submission.UpdatedAt = time.Now()

	err = submissionRepository.Update(context.TODO(), &submission)
	assert.NoError(t, err)

	// Verify update by retrieving submission
	updatedSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedSubmission)
	assert.Equal(t, model.SubmissionStatusSubmitted, updatedSubmission.Status)
	assert.Equal(t, 3, len(updatedSubmission.Answers))
	assert.Equal(t, &score, updatedSubmission.Score)
	assert.Equal(t, "Great work!", updatedSubmission.Feedback)
	assert.NotNil(t, updatedSubmission.SubmittedAt)
}

func TestUpdateSubmissionPartial(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := createTestSubmission()

	// Create submission
	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)

	// Update only feedback
	submission.Feedback = "Needs improvement"
	submission.UpdatedAt = time.Now()

	err = submissionRepository.Update(context.TODO(), &submission)
	assert.NoError(t, err)

	// Verify partial update
	updatedSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedSubmission)
	assert.Equal(t, "Needs improvement", updatedSubmission.Feedback)
	assert.Equal(t, model.SubmissionStatusDraft, updatedSubmission.Status) // Should remain unchanged
	assert.Nil(t, updatedSubmission.Score)                                 // Should remain unchanged
}

// Tests for GetByID
func TestSubmissionGetByID(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := createTestSubmission()

	// Create submission
	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)

	// Get submission by ID
	gotSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, gotSubmission)

	// Verify submission details
	assert.Equal(t, submission.ID, gotSubmission.ID)
	assert.Equal(t, submission.AssignmentID, gotSubmission.AssignmentID)
	assert.Equal(t, submission.StudentUUID, gotSubmission.StudentUUID)
	assert.Equal(t, submission.StudentName, gotSubmission.StudentName)
	assert.Equal(t, submission.Status, gotSubmission.Status)
	assert.Equal(t, len(submission.Answers), len(gotSubmission.Answers))
	assert.Equal(t, submission.Score, gotSubmission.Score)
	assert.Equal(t, submission.Feedback, gotSubmission.Feedback)
}

func TestSubmissionGetByIDNotFound(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Try to get non-existent submission
	nonExistentID := primitive.NewObjectID()
	gotSubmission, err := submissionRepository.GetByID(context.TODO(), nonExistentID.Hex())
	assert.NoError(t, err)
	assert.Nil(t, gotSubmission)
}

func TestSubmissionGetByIDWithInvalidID(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Try to get submission with invalid ID
	gotSubmission, err := submissionRepository.GetByID(context.TODO(), "invalid-id")
	assert.Error(t, err)
	assert.Nil(t, gotSubmission)
}

// Tests for GetByAssignmentAndStudent
func TestGetByAssignmentAndStudent(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	submission := createTestSubmissionWithDetails("assignment123", "student123", "Test Student", model.SubmissionStatusDraft)

	// Create submission
	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)

	// Get submission by assignment and student
	gotSubmission, err := submissionRepository.GetByAssignmentAndStudent(context.TODO(), "assignment123", "student123")
	assert.NoError(t, err)
	assert.NotNil(t, gotSubmission)

	// Verify submission details
	assert.Equal(t, submission.ID, gotSubmission.ID)
	assert.Equal(t, "assignment123", gotSubmission.AssignmentID)
	assert.Equal(t, "student123", gotSubmission.StudentUUID)
	assert.Equal(t, "Test Student", gotSubmission.StudentName)
}

func TestGetByAssignmentAndStudentNotFound(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Try to get non-existent submission
	gotSubmission, err := submissionRepository.GetByAssignmentAndStudent(context.TODO(), "assignment999", "student999")
	assert.NoError(t, err)
	assert.Nil(t, gotSubmission)
}

func TestGetByAssignmentAndStudentWithOtherSubmissions(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Create multiple submissions
	submission1 := createTestSubmissionWithDetails("assignment123", "student123", "Student 1", model.SubmissionStatusDraft)
	submission2 := createTestSubmissionWithDetails("assignment123", "student456", "Student 2", model.SubmissionStatusSubmitted)
	submission3 := createTestSubmissionWithDetails("assignment456", "student123", "Student 1", model.SubmissionStatusLate)

	err := submissionRepository.Create(context.TODO(), &submission1)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission2)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission3)
	assert.NoError(t, err)

	// Get specific submission
	gotSubmission, err := submissionRepository.GetByAssignmentAndStudent(context.TODO(), "assignment123", "student123")
	assert.NoError(t, err)
	assert.NotNil(t, gotSubmission)
	assert.Equal(t, submission1.ID, gotSubmission.ID)
	assert.Equal(t, "assignment123", gotSubmission.AssignmentID)
	assert.Equal(t, "student123", gotSubmission.StudentUUID)
	assert.Equal(t, "Student 1", gotSubmission.StudentName)
}

// Tests for GetByAssignment
func TestGetByAssignment(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Create submissions for the same assignment
	submission1 := createTestSubmissionWithDetails("assignment123", "student1", "Student 1", model.SubmissionStatusSubmitted)
	submission2 := createTestSubmissionWithDetails("assignment123", "student2", "Student 2", model.SubmissionStatusDraft)
	submission3 := createTestSubmissionWithDetails("assignment456", "student3", "Student 3", model.SubmissionStatusLate)

	err := submissionRepository.Create(context.TODO(), &submission1)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission2)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission3)
	assert.NoError(t, err)

	// Get submissions by assignment
	submissions, err := submissionRepository.GetByAssignment(context.TODO(), "assignment123")
	assert.NoError(t, err)
	assert.NotNil(t, submissions)
	assert.Equal(t, 2, len(submissions))

	// Verify submissions belong to the correct assignment
	studentUUIDs := []string{submissions[0].StudentUUID, submissions[1].StudentUUID}
	assert.Contains(t, studentUUIDs, "student1")
	assert.Contains(t, studentUUIDs, "student2")

	for _, submission := range submissions {
		assert.Equal(t, "assignment123", submission.AssignmentID)
	}
}

func TestGetByAssignmentEmpty(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Get submissions for non-existent assignment
	submissions, err := submissionRepository.GetByAssignment(context.TODO(), "assignment999")
	assert.NoError(t, err)
	assert.NotNil(t, submissions)
	assert.Equal(t, 0, len(submissions))
}

// Tests for GetByStudent
func TestGetByStudent(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Create submissions for the same student
	submission1 := createTestSubmissionWithDetails("assignment123", "student123", "Test Student", model.SubmissionStatusSubmitted)
	submission2 := createTestSubmissionWithDetails("assignment456", "student123", "Test Student", model.SubmissionStatusDraft)
	submission3 := createTestSubmissionWithDetails("assignment789", "student456", "Other Student", model.SubmissionStatusLate)

	err := submissionRepository.Create(context.TODO(), &submission1)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission2)
	assert.NoError(t, err)
	err = submissionRepository.Create(context.TODO(), &submission3)
	assert.NoError(t, err)

	// Get submissions by student
	submissions, err := submissionRepository.GetByStudent(context.TODO(), "student123")
	assert.NoError(t, err)
	assert.NotNil(t, submissions)
	assert.Equal(t, 2, len(submissions))

	// Verify submissions belong to the correct student
	assignmentIDs := []string{submissions[0].AssignmentID, submissions[1].AssignmentID}
	assert.Contains(t, assignmentIDs, "assignment123")
	assert.Contains(t, assignmentIDs, "assignment456")

	for _, submission := range submissions {
		assert.Equal(t, "student123", submission.StudentUUID)
		assert.Equal(t, "Test Student", submission.StudentName)
	}
}

func TestGetByStudentEmpty(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// Get submissions for non-existent student
	submissions, err := submissionRepository.GetByStudent(context.TODO(), "student999")
	assert.NoError(t, err)
	assert.NotNil(t, submissions)
	assert.Equal(t, 0, len(submissions))
}

// Integration test: Complete submission workflow
func TestCompleteSubmissionWorkflow(t *testing.T) {
	t.Cleanup(func() {
		submissionDBSetup.CleanupCollection("submissions")
	})

	submissionRepository := repository.NewMongoSubmissionRepository(submissionDBSetup.Client.Database(submissionDBSetup.DBName))

	// 1. Create initial submission (draft)
	submission := createTestSubmissionWithDetails("assignment123", "student123", "Test Student", model.SubmissionStatusDraft)

	err := submissionRepository.Create(context.TODO(), &submission)
	assert.NoError(t, err)
	assert.False(t, submission.ID.IsZero())

	// 2. Verify creation
	gotSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, gotSubmission)
	assert.Equal(t, model.SubmissionStatusDraft, gotSubmission.Status)
	assert.Nil(t, gotSubmission.Score)
	assert.Equal(t, "", gotSubmission.Feedback)

	// 3. Update submission (submit)
	submission.Status = model.SubmissionStatusSubmitted
	submittedAt := time.Now()
	submission.SubmittedAt = &submittedAt
	submission.UpdatedAt = time.Now()

	err = submissionRepository.Update(context.TODO(), &submission)
	assert.NoError(t, err)

	// 4. Verify submission status
	submittedSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, submittedSubmission)
	assert.Equal(t, model.SubmissionStatusSubmitted, submittedSubmission.Status)
	assert.NotNil(t, submittedSubmission.SubmittedAt)

	// 5. Grade submission
	score := 92.5
	submission.Score = &score
	submission.Feedback = "Excellent work! Great understanding of the concepts."
	submission.UpdatedAt = time.Now()

	err = submissionRepository.Update(context.TODO(), &submission)
	assert.NoError(t, err)

	// 6. Verify grading
	gradedSubmission, err := submissionRepository.GetByID(context.TODO(), submission.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, gradedSubmission)
	assert.Equal(t, &score, gradedSubmission.Score)
	assert.Equal(t, "Excellent work! Great understanding of the concepts.", gradedSubmission.Feedback)

	// 7. Verify submission can be found by assignment and student
	foundSubmission, err := submissionRepository.GetByAssignmentAndStudent(context.TODO(), "assignment123", "student123")
	assert.NoError(t, err)
	assert.NotNil(t, foundSubmission)
	assert.Equal(t, submission.ID, foundSubmission.ID)

	// 8. Verify submission appears in assignment submissions
	assignmentSubmissions, err := submissionRepository.GetByAssignment(context.TODO(), "assignment123")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(assignmentSubmissions))
	assert.Equal(t, submission.ID, assignmentSubmissions[0].ID)

	// 9. Verify submission appears in student submissions
	studentSubmissions, err := submissionRepository.GetByStudent(context.TODO(), "student123")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(studentSubmissions))
	assert.Equal(t, submission.ID, studentSubmissions[0].ID)
}
