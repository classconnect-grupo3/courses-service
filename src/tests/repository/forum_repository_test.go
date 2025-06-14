package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to create a test course
func createTestCourse(t *testing.T, courseRepo *repository.CourseRepository) *model.Course {
	course := model.Course{
		Title:          "Test Course for Forum",
		Description:    "Test Description",
		TeacherUUID:    "teacher-123",
		TeacherName:    "Test Teacher",
		Capacity:       30,
		StudentsAmount: 0,
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(24 * time.Hour * 30),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Modules:        []model.Module{},
	}

	createdCourse, err := courseRepo.CreateCourse(course)
	if err != nil {
		t.Fatalf("Failed to create test course: %v", err)
	}
	return createdCourse
}

// Helper function to create a test question
func createTestQuestion(t *testing.T, forumRepo *repository.ForumRepository, courseID string) *model.ForumQuestion {
	question := model.ForumQuestion{
		CourseID:    courseID,
		AuthorID:    "author-123",
		Title:       "Test Question",
		Description: "This is a test question",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral, model.QuestionTagTeoria},
	}

	createdQuestion, err := forumRepo.CreateQuestion(question)
	if err != nil {
		t.Fatalf("Failed to create test question: %v", err)
	}
	return createdQuestion
}

func TestCreateQuestion(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course
	course := createTestCourse(t, courseRepo)

	// Test creating a question
	question := model.ForumQuestion{
		CourseID:    course.ID.Hex(),
		AuthorID:    "author-123",
		Title:       "How to implement clean architecture?",
		Description: "I need help understanding clean architecture principles",
		Tags:        []model.QuestionTag{model.QuestionTagTeoria, model.QuestionTagNecesitoAyuda},
	}

	createdQuestion, err := forumRepo.CreateQuestion(question)
	if err != nil {
		t.Fatalf("Failed to create question: %v", err)
	}

	// Verify question properties
	if createdQuestion.ID.IsZero() {
		t.Error("Expected question to have an ID after creation")
	}

	if createdQuestion.CourseID != question.CourseID {
		t.Errorf("Expected course ID %s, got %s", question.CourseID, createdQuestion.CourseID)
	}

	if createdQuestion.AuthorID != question.AuthorID {
		t.Errorf("Expected author ID %s, got %s", question.AuthorID, createdQuestion.AuthorID)
	}

	if createdQuestion.Title != question.Title {
		t.Errorf("Expected title %s, got %s", question.Title, createdQuestion.Title)
	}

	if createdQuestion.Description != question.Description {
		t.Errorf("Expected description %s, got %s", question.Description, createdQuestion.Description)
	}

	if len(createdQuestion.Tags) != len(question.Tags) {
		t.Errorf("Expected %d tags, got %d", len(question.Tags), len(createdQuestion.Tags))
	}

	if createdQuestion.Status != model.QuestionStatusOpen {
		t.Errorf("Expected status %s, got %s", model.QuestionStatusOpen, createdQuestion.Status)
	}

	if len(createdQuestion.Votes) != 0 {
		t.Errorf("Expected 0 votes, got %d", len(createdQuestion.Votes))
	}

	if len(createdQuestion.Answers) != 0 {
		t.Errorf("Expected 0 answers, got %d", len(createdQuestion.Answers))
	}

	if createdQuestion.AcceptedAnswerID != nil {
		t.Error("Expected no accepted answer ID")
	}
}

func TestGetQuestionById(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Test getting question by ID
	retrievedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get question by ID: %v", err)
	}

	if retrievedQuestion.ID != question.ID {
		t.Errorf("Expected question ID %s, got %s", question.ID.Hex(), retrievedQuestion.ID.Hex())
	}

	if retrievedQuestion.Title != question.Title {
		t.Errorf("Expected title %s, got %s", question.Title, retrievedQuestion.Title)
	}

	// Test with invalid ID
	_, err = forumRepo.GetQuestionById("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}

	// Test with non-existent ID
	nonExistentID := primitive.NewObjectID().Hex()
	_, err = forumRepo.GetQuestionById(nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent question ID, got nil")
	}
}

func TestGetQuestionsByCourseId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create test courses
	course1 := createTestCourse(t, courseRepo)
	course2 := createTestCourse(t, courseRepo)

	// Create questions for course1
	createTestQuestion(t, forumRepo, course1.ID.Hex())
	createTestQuestion(t, forumRepo, course1.ID.Hex())

	// Create question for course2
	createTestQuestion(t, forumRepo, course2.ID.Hex())

	// Test getting questions for course1
	questions, err := forumRepo.GetQuestionsByCourseId(course1.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get questions by course ID: %v", err)
	}

	if len(questions) != 2 {
		t.Errorf("Expected 2 questions for course1, got %d", len(questions))
	}

	// Verify questions are sorted by created_at descending (newest first)
	if len(questions) >= 2 {
		if questions[0].CreatedAt.Before(questions[1].CreatedAt) {
			t.Error("Expected questions to be sorted by created_at descending")
		}
	}

	// Test getting questions for course2
	questions, err = forumRepo.GetQuestionsByCourseId(course2.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get questions by course ID: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 question for course2, got %d", len(questions))
	}

	// Test with non-existent course
	questions, err = forumRepo.GetQuestionsByCourseId("non-existent-course")
	if err != nil {
		t.Fatalf("Failed to get questions for non-existent course: %v", err)
	}

	if len(questions) != 0 {
		t.Errorf("Expected 0 questions for non-existent course, got %d", len(questions))
	}
}

func TestUpdateQuestion(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Add a small delay to ensure updated_at will be different
	time.Sleep(10 * time.Millisecond)

	// Test updating question
	updateData := model.ForumQuestion{
		Title:       "Updated Question Title",
		Description: "Updated description",
		Tags:        []model.QuestionTag{model.QuestionTagPractica},
		Status:      model.QuestionStatusClosed,
	}

	updatedQuestion, err := forumRepo.UpdateQuestion(question.ID.Hex(), updateData)
	if err != nil {
		t.Fatalf("Failed to update question: %v", err)
	}

	if updatedQuestion.Title != updateData.Title {
		t.Errorf("Expected title %s, got %s", updateData.Title, updatedQuestion.Title)
	}

	if updatedQuestion.Description != updateData.Description {
		t.Errorf("Expected description %s, got %s", updateData.Description, updatedQuestion.Description)
	}

	if len(updatedQuestion.Tags) != len(updateData.Tags) {
		t.Errorf("Expected %d tags, got %d", len(updateData.Tags), len(updatedQuestion.Tags))
	}

	if updatedQuestion.Status != updateData.Status {
		t.Errorf("Expected status %s, got %s", updateData.Status, updatedQuestion.Status)
	}

	// Test with invalid ID
	_, err = forumRepo.UpdateQuestion("invalid-id", updateData)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestDeleteQuestion(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Test deleting question
	err := forumRepo.DeleteQuestion(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete question: %v", err)
	}

	// Verify question is deleted
	_, err = forumRepo.GetQuestionById(question.ID.Hex())
	if err == nil {
		t.Error("Expected error when getting deleted question, got nil")
	}

	// Test deleting non-existent question
	err = forumRepo.DeleteQuestion(question.ID.Hex())
	if err == nil {
		t.Error("Expected error when deleting non-existent question, got nil")
	}

	// Test with invalid ID
	err = forumRepo.DeleteQuestion("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestAddAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Test adding answer
	answer := model.ForumAnswer{
		AuthorID: "answer-author-123",
		Content:  "This is a test answer",
	}

	addedAnswer, err := forumRepo.AddAnswer(question.ID.Hex(), answer)
	if err != nil {
		t.Fatalf("Failed to add answer: %v", err)
	}

	if addedAnswer.ID == "" {
		t.Error("Expected answer to have an ID after creation")
	}

	if addedAnswer.AuthorID != answer.AuthorID {
		t.Errorf("Expected author ID %s, got %s", answer.AuthorID, addedAnswer.AuthorID)
	}

	if addedAnswer.Content != answer.Content {
		t.Errorf("Expected content %s, got %s", answer.Content, addedAnswer.Content)
	}

	if len(addedAnswer.Votes) != 0 {
		t.Errorf("Expected 0 votes, got %d", len(addedAnswer.Votes))
	}

	if addedAnswer.IsAccepted {
		t.Error("Expected answer to not be accepted initially")
	}

	// Verify answer is added to question
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Answers) != 1 {
		t.Errorf("Expected 1 answer in question, got %d", len(updatedQuestion.Answers))
	}

	// Test with invalid question ID
	_, err = forumRepo.AddAnswer("invalid-id", answer)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestUpdateAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course, question, and answer
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())
	answer := model.ForumAnswer{
		AuthorID: "answer-author-123",
		Content:  "Original answer content",
	}
	addedAnswer, _ := forumRepo.AddAnswer(question.ID.Hex(), answer)

	// Test updating answer
	newContent := "Updated answer content"
	updatedAnswer, err := forumRepo.UpdateAnswer(question.ID.Hex(), addedAnswer.ID, newContent)
	if err != nil {
		t.Fatalf("Failed to update answer: %v", err)
	}

	if updatedAnswer.Content != newContent {
		t.Errorf("Expected content %s, got %s", newContent, updatedAnswer.Content)
	}

	// Test with invalid question ID
	_, err = forumRepo.UpdateAnswer("invalid-id", addedAnswer.ID, newContent)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}

	// Test with invalid answer ID
	_, err = forumRepo.UpdateAnswer(question.ID.Hex(), "invalid-answer-id", newContent)
	if err == nil {
		t.Error("Expected error for invalid answer ID, got nil")
	}
}

func TestDeleteAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course, question, and answer
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())
	answer := model.ForumAnswer{
		AuthorID: "answer-author-123",
		Content:  "Answer to be deleted",
	}
	addedAnswer, _ := forumRepo.AddAnswer(question.ID.Hex(), answer)

	// Test deleting answer
	err := forumRepo.DeleteAnswer(question.ID.Hex(), addedAnswer.ID)
	if err != nil {
		t.Fatalf("Failed to delete answer: %v", err)
	}

	// Verify answer is deleted
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Answers) != 0 {
		t.Errorf("Expected 0 answers after deletion, got %d", len(updatedQuestion.Answers))
	}

	// Test with invalid question ID
	err = forumRepo.DeleteAnswer("invalid-id", addedAnswer.ID)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestAcceptAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course, question, and answers
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	answer1 := model.ForumAnswer{
		AuthorID: "answer-author-1",
		Content:  "First answer",
	}
	addedAnswer1, _ := forumRepo.AddAnswer(question.ID.Hex(), answer1)

	answer2 := model.ForumAnswer{
		AuthorID: "answer-author-2",
		Content:  "Second answer",
	}
	addedAnswer2, _ := forumRepo.AddAnswer(question.ID.Hex(), answer2)

	// Test accepting first answer
	err := forumRepo.AcceptAnswer(question.ID.Hex(), addedAnswer1.ID)
	if err != nil {
		t.Fatalf("Failed to accept answer: %v", err)
	}

	// Verify answer is accepted and question status is updated
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if updatedQuestion.Status != model.QuestionStatusResolved {
		t.Errorf("Expected question status %s, got %s", model.QuestionStatusResolved, updatedQuestion.Status)
	}

	if updatedQuestion.AcceptedAnswerID == nil || *updatedQuestion.AcceptedAnswerID != addedAnswer1.ID {
		t.Errorf("Expected accepted answer ID %s, got %v", addedAnswer1.ID, updatedQuestion.AcceptedAnswerID)
	}

	// Find and verify the accepted answer
	var acceptedAnswer *model.ForumAnswer
	for _, ans := range updatedQuestion.Answers {
		if ans.ID == addedAnswer1.ID {
			acceptedAnswer = &ans
			break
		}
	}

	if acceptedAnswer == nil {
		t.Fatal("Could not find accepted answer")
	}

	if !acceptedAnswer.IsAccepted {
		t.Error("Expected answer to be marked as accepted")
	}

	// Test accepting second answer (should unmark first)
	err = forumRepo.AcceptAnswer(question.ID.Hex(), addedAnswer2.ID)
	if err != nil {
		t.Fatalf("Failed to accept second answer: %v", err)
	}

	// Verify only second answer is accepted
	updatedQuestion, err = forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	acceptedCount := 0
	for _, ans := range updatedQuestion.Answers {
		if ans.IsAccepted {
			acceptedCount++
			if ans.ID != addedAnswer2.ID {
				t.Errorf("Expected answer %s to be accepted, but %s is accepted", addedAnswer2.ID, ans.ID)
			}
		}
	}

	if acceptedCount != 1 {
		t.Errorf("Expected exactly 1 accepted answer, got %d", acceptedCount)
	}

	// Test with invalid question ID
	err = forumRepo.AcceptAnswer("invalid-id", addedAnswer1.ID)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}

	// Test with invalid answer ID
	err = forumRepo.AcceptAnswer(question.ID.Hex(), "invalid-answer-id")
	if err == nil {
		t.Error("Expected error for invalid answer ID, got nil")
	}
}

func TestAddVoteToQuestion(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Test adding upvote
	err := forumRepo.AddVoteToQuestion(question.ID.Hex(), "voter-123", model.VoteTypeUp)
	if err != nil {
		t.Fatalf("Failed to add upvote to question: %v", err)
	}

	// Verify vote is added
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Votes) != 1 {
		t.Errorf("Expected 1 vote, got %d", len(updatedQuestion.Votes))
	}

	if updatedQuestion.Votes[0].UserID != "voter-123" {
		t.Errorf("Expected voter ID voter-123, got %s", updatedQuestion.Votes[0].UserID)
	}

	if updatedQuestion.Votes[0].VoteType != model.VoteTypeUp {
		t.Errorf("Expected vote type %d, got %d", model.VoteTypeUp, updatedQuestion.Votes[0].VoteType)
	}

	// Test changing vote (should replace existing vote)
	err = forumRepo.AddVoteToQuestion(question.ID.Hex(), "voter-123", model.VoteTypeDown)
	if err != nil {
		t.Fatalf("Failed to change vote on question: %v", err)
	}

	// Verify vote is replaced
	updatedQuestion, err = forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Votes) != 1 {
		t.Errorf("Expected 1 vote after change, got %d", len(updatedQuestion.Votes))
	}

	if updatedQuestion.Votes[0].VoteType != model.VoteTypeDown {
		t.Errorf("Expected vote type %d, got %d", model.VoteTypeDown, updatedQuestion.Votes[0].VoteType)
	}

	// Test adding vote from different user
	err = forumRepo.AddVoteToQuestion(question.ID.Hex(), "voter-456", model.VoteTypeUp)
	if err != nil {
		t.Fatalf("Failed to add vote from different user: %v", err)
	}

	// Verify both votes exist
	updatedQuestion, err = forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Votes) != 2 {
		t.Errorf("Expected 2 votes, got %d", len(updatedQuestion.Votes))
	}

	// Test with invalid question ID
	err = forumRepo.AddVoteToQuestion("invalid-id", "voter-123", model.VoteTypeUp)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestAddVoteToAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course, question, and answer
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())
	answer := model.ForumAnswer{
		AuthorID: "answer-author-123",
		Content:  "Test answer for voting",
	}
	addedAnswer, _ := forumRepo.AddAnswer(question.ID.Hex(), answer)

	// Test adding upvote to answer
	err := forumRepo.AddVoteToAnswer(question.ID.Hex(), addedAnswer.ID, "voter-123", model.VoteTypeUp)
	if err != nil {
		t.Fatalf("Failed to add upvote to answer: %v", err)
	}

	// Verify vote is added
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Answers) != 1 {
		t.Fatalf("Expected 1 answer, got %d", len(updatedQuestion.Answers))
	}

	answerVotes := updatedQuestion.Answers[0].Votes
	if len(answerVotes) != 1 {
		t.Errorf("Expected 1 vote on answer, got %d", len(answerVotes))
	}

	if answerVotes[0].UserID != "voter-123" {
		t.Errorf("Expected voter ID voter-123, got %s", answerVotes[0].UserID)
	}

	if answerVotes[0].VoteType != model.VoteTypeUp {
		t.Errorf("Expected vote type %d, got %d", model.VoteTypeUp, answerVotes[0].VoteType)
	}

	// Test with invalid question ID
	err = forumRepo.AddVoteToAnswer("invalid-id", addedAnswer.ID, "voter-123", model.VoteTypeUp)
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}

	// Test with invalid answer ID
	err = forumRepo.AddVoteToAnswer(question.ID.Hex(), "invalid-answer-id", "voter-123", model.VoteTypeUp)
	if err == nil {
		t.Error("Expected error for invalid answer ID, got nil")
	}
}

func TestRemoveVoteFromQuestion(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course and question
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())

	// Add a vote first
	forumRepo.AddVoteToQuestion(question.ID.Hex(), "voter-123", model.VoteTypeUp)

	// Test removing vote
	err := forumRepo.RemoveVoteFromQuestion(question.ID.Hex(), "voter-123")
	if err != nil {
		t.Fatalf("Failed to remove vote from question: %v", err)
	}

	// Verify vote is removed
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Votes) != 0 {
		t.Errorf("Expected 0 votes after removal, got %d", len(updatedQuestion.Votes))
	}

	// Test removing non-existent vote (should not error)
	err = forumRepo.RemoveVoteFromQuestion(question.ID.Hex(), "non-existent-voter")
	if err != nil {
		t.Fatalf("Unexpected error when removing non-existent vote: %v", err)
	}

	// Test with invalid question ID
	err = forumRepo.RemoveVoteFromQuestion("invalid-id", "voter-123")
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestRemoveVoteFromAnswer(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course, question, and answer
	course := createTestCourse(t, courseRepo)
	question := createTestQuestion(t, forumRepo, course.ID.Hex())
	answer := model.ForumAnswer{
		AuthorID: "answer-author-123",
		Content:  "Test answer for vote removal",
	}
	addedAnswer, _ := forumRepo.AddAnswer(question.ID.Hex(), answer)

	// Add a vote first
	forumRepo.AddVoteToAnswer(question.ID.Hex(), addedAnswer.ID, "voter-123", model.VoteTypeUp)

	// Test removing vote
	err := forumRepo.RemoveVoteFromAnswer(question.ID.Hex(), addedAnswer.ID, "voter-123")
	if err != nil {
		t.Fatalf("Failed to remove vote from answer: %v", err)
	}

	// Verify vote is removed
	updatedQuestion, err := forumRepo.GetQuestionById(question.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get updated question: %v", err)
	}

	if len(updatedQuestion.Answers[0].Votes) != 0 {
		t.Errorf("Expected 0 votes on answer after removal, got %d", len(updatedQuestion.Answers[0].Votes))
	}

	// Test with invalid question ID
	err = forumRepo.RemoveVoteFromAnswer("invalid-id", addedAnswer.ID, "voter-123")
	if err == nil {
		t.Error("Expected error for invalid question ID, got nil")
	}
}

func TestSearchQuestions(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("forum_questions")
	})

	forumRepo := repository.NewForumRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create test courses
	course1 := createTestCourse(t, courseRepo)
	course2 := createTestCourse(t, courseRepo)

	// Create questions with different properties
	question1 := model.ForumQuestion{
		CourseID:    course1.ID.Hex(),
		AuthorID:    "author-1",
		Title:       "How to implement clean architecture?",
		Description: "I need help with clean architecture patterns",
		Tags:        []model.QuestionTag{model.QuestionTagTeoria, model.QuestionTagNecesitoAyuda},
		Status:      model.QuestionStatusOpen,
	}
	forumRepo.CreateQuestion(question1)

	question2 := model.ForumQuestion{
		CourseID:    course1.ID.Hex(),
		AuthorID:    "author-2",
		Title:       "Database design best practices",
		Description: "What are the best practices for database design?",
		Tags:        []model.QuestionTag{model.QuestionTagPractica},
		Status:      model.QuestionStatusResolved,
	}
	forumRepo.CreateQuestion(question2)

	question3 := model.ForumQuestion{
		CourseID:    course2.ID.Hex(),
		AuthorID:    "author-3",
		Title:       "Testing strategies",
		Description: "How to write effective tests?",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral},
		Status:      model.QuestionStatusOpen,
	}
	forumRepo.CreateQuestion(question3)

	// Test search by course ID only
	questions, err := forumRepo.SearchQuestions(course1.ID.Hex(), "", nil, "")
	if err != nil {
		t.Fatalf("Failed to search questions: %v", err)
	}

	if len(questions) != 2 {
		t.Errorf("Expected 2 questions for course1, got %d", len(questions))
	}

	// Test search by query text
	questions, err = forumRepo.SearchQuestions(course1.ID.Hex(), "architecture", nil, "")
	if err != nil {
		t.Fatalf("Failed to search questions by query: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 question matching 'architecture', got %d", len(questions))
	}

	if questions[0].Title != question1.Title {
		t.Errorf("Expected question with title '%s', got '%s'", question1.Title, questions[0].Title)
	}

	// Test search by tags
	questions, err = forumRepo.SearchQuestions(course1.ID.Hex(), "", []model.QuestionTag{model.QuestionTagPractica}, "")
	if err != nil {
		t.Fatalf("Failed to search questions by tags: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 question with 'practica' tag, got %d", len(questions))
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 resolved question, got %d", len(questions))
	}

	// Test search with multiple filters
	questions, err = forumRepo.SearchQuestions(course1.ID.Hex(), "", []model.QuestionTag{model.QuestionTagTeoria}, model.QuestionStatusOpen)
	if err != nil {
		t.Fatalf("Failed to search questions with multiple filters: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 question matching multiple filters, got %d", len(questions))
	}

	// Test search with no results
	questions, err = forumRepo.SearchQuestions(course1.ID.Hex(), "nonexistent", nil, "")
	if err != nil {
		t.Fatalf("Failed to search questions with no results: %v", err)
	}

	if len(questions) != 0 {
		t.Errorf("Expected 0 questions for non-matching query, got %d", len(questions))
	}
}
