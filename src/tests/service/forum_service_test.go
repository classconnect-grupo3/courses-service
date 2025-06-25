package service_test

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock Forum Repository
type MockForumRepository struct{}

func (m *MockForumRepository) CreateQuestion(question model.ForumQuestion) (*model.ForumQuestion, error) {
	if question.CourseID == "error-course" {
		return nil, errors.New("failed to create question")
	}

	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()
	question.Status = model.QuestionStatusOpen
	question.Votes = []model.Vote{}
	question.Answers = []model.ForumAnswer{}

	return &question, nil
}

func (m *MockForumRepository) GetQuestionById(id string) (*model.ForumQuestion, error) {
	if id == "non-existent-question" {
		return nil, errors.New("question not found")
	}
	if id == "question-with-answers" {
		return &model.ForumQuestion{
			ID:          mustParseForumObjectID("123456789012345678901234"),
			CourseID:    "course-123",
			AuthorID:    "author-123",
			Title:       "Test Question",
			Description: "Test Description",
			Tags:        []model.QuestionTag{model.QuestionTagGeneral},
			Status:      model.QuestionStatusOpen,
			Answers: []model.ForumAnswer{
				{
					ID:       "answer-123",
					AuthorID: "answer-author-123",
					Content:  "Test Answer",
				},
				{
					ID:       "answer-456",
					AuthorID: "answer-author-456",
					Content:  "Another Answer",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	if id == "question-author-123" {
		return &model.ForumQuestion{
			ID:          mustParseForumObjectID("123456789012345678901234"),
			CourseID:    "course-123",
			AuthorID:    "author-123",
			Title:       "Question by Author 123",
			Description: "Test Description",
			Tags:        []model.QuestionTag{model.QuestionTagGeneral},
			Status:      model.QuestionStatusOpen,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if id == "question-author-456" {
		return &model.ForumQuestion{
			ID:          mustParseForumObjectID("123456789012345678901235"),
			CourseID:    "course-123",
			AuthorID:    "author-456",
			Title:       "Question by Author 456",
			Description: "Test Description",
			Tags:        []model.QuestionTag{model.QuestionTagGeneral},
			Status:      model.QuestionStatusOpen,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	return &model.ForumQuestion{
		ID:          mustParseForumObjectID("123456789012345678901234"),
		CourseID:    "course-123",
		AuthorID:    "author-123",
		Title:       "Test Question",
		Description: "Test Description",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral},
		Status:      model.QuestionStatusOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *MockForumRepository) GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error) {
	if courseID == "error-course" {
		return nil, errors.New("failed to get questions")
	}
	if courseID == "empty-course" {
		return []model.ForumQuestion{}, nil
	}

	return []model.ForumQuestion{
		{
			ID:          mustParseForumObjectID("123456789012345678901234"),
			CourseID:    courseID,
			AuthorID:    "author-123",
			Title:       "Question 1",
			Description: "Description 1",
			Tags:        []model.QuestionTag{model.QuestionTagGeneral},
			Status:      model.QuestionStatusOpen,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          mustParseForumObjectID("123456789012345678901235"),
			CourseID:    courseID,
			AuthorID:    "author-456",
			Title:       "Question 2",
			Description: "Description 2",
			Tags:        []model.QuestionTag{model.QuestionTagTeoria},
			Status:      model.QuestionStatusResolved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

func (m *MockForumRepository) UpdateQuestion(id string, question model.ForumQuestion) (*model.ForumQuestion, error) {
	if id == "non-existent-question" {
		return nil, errors.New("question not found")
	}

	updatedQuestion := &model.ForumQuestion{
		ID:          mustParseForumObjectID("123456789012345678901234"),
		CourseID:    "course-123",
		AuthorID:    "author-123",
		Title:       question.Title,
		Description: question.Description,
		Tags:        question.Tags,
		Status:      question.Status,
		UpdatedAt:   time.Now(),
	}

	return updatedQuestion, nil
}

func (m *MockForumRepository) DeleteQuestion(id string) error {
	if id == "non-existent-question" {
		return errors.New("question not found")
	}
	return nil
}

func (m *MockForumRepository) AddAnswer(questionID string, answer model.ForumAnswer) (*model.ForumAnswer, error) {
	if questionID == "non-existent-question" {
		return nil, errors.New("question not found")
	}

	answer.ID = "new-answer-id"
	answer.CreatedAt = time.Now()
	answer.UpdatedAt = time.Now()
	answer.Votes = []model.Vote{}
	answer.IsAccepted = false

	return &answer, nil
}

func (m *MockForumRepository) UpdateAnswer(questionID, answerID, content string) (*model.ForumAnswer, error) {
	if questionID == "non-existent-question" || answerID == "non-existent-answer" {
		return nil, errors.New("question or answer not found")
	}

	return &model.ForumAnswer{
		ID:        answerID,
		AuthorID:  "answer-author-123",
		Content:   content,
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockForumRepository) DeleteAnswer(questionID, answerID string) error {
	if questionID == "non-existent-question" || answerID == "non-existent-answer" {
		return errors.New("question or answer not found")
	}
	return nil
}

func (m *MockForumRepository) AcceptAnswer(questionID, answerID string) error {
	if questionID == "non-existent-question" || answerID == "non-existent-answer" {
		return errors.New("question or answer not found")
	}
	return nil
}

func (m *MockForumRepository) AddVoteToQuestion(questionID, userID string, voteType int) error {
	if questionID == "non-existent-question" {
		return errors.New("question not found")
	}
	return nil
}

func (m *MockForumRepository) AddVoteToAnswer(questionID, answerID, userID string, voteType int) error {
	if questionID == "non-existent-question" || answerID == "non-existent-answer" {
		return errors.New("question or answer not found")
	}
	return nil
}

func (m *MockForumRepository) RemoveVoteFromQuestion(questionID, userID string) error {
	if questionID == "non-existent-question" {
		return errors.New("question not found")
	}
	return nil
}

func (m *MockForumRepository) RemoveVoteFromAnswer(questionID, answerID, userID string) error {
	if questionID == "non-existent-question" || answerID == "non-existent-answer" {
		return errors.New("question or answer not found")
	}
	return nil
}

func (m *MockForumRepository) SearchQuestions(courseID, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error) {
	if courseID == "error-course" {
		return nil, errors.New("failed to search questions")
	}

	// Return filtered results based on parameters
	questions := []model.ForumQuestion{
		{
			ID:          mustParseForumObjectID("123456789012345678901234"),
			CourseID:    courseID,
			AuthorID:    "author-123",
			Title:       "Architecture Question",
			Description: "How to implement clean architecture?",
			Tags:        []model.QuestionTag{model.QuestionTagTeoria},
			Status:      model.QuestionStatusOpen,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          mustParseForumObjectID("123456789012345678901235"),
			CourseID:    courseID,
			AuthorID:    "author-456",
			Title:       "Database Question",
			Description: "Database best practices",
			Tags:        []model.QuestionTag{model.QuestionTagPractica},
			Status:      model.QuestionStatusResolved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Simple filtering for testing
	if query == "architecture" {
		return []model.ForumQuestion{questions[0]}, nil
	}
	if len(tags) > 0 && tags[0] == model.QuestionTagPractica {
		return []model.ForumQuestion{questions[1]}, nil
	}
	if status == model.QuestionStatusResolved {
		return []model.ForumQuestion{questions[1]}, nil
	}

	return questions, nil
}

// Backoffice statistics methods for MockForumRepository
func (m *MockForumRepository) CountQuestions() (int64, error) {
	return 2, nil
}

func (m *MockForumRepository) CountQuestionsByStatus(status model.QuestionStatus) (int64, error) {
	if status == model.QuestionStatusOpen {
		return 1, nil
	}
	if status == model.QuestionStatusResolved {
		return 1, nil
	}
	return 0, nil
}

func (m *MockForumRepository) CountAnswers() (int64, error) {
	return 3, nil
}

// Mock Course Repository (reusing from existing tests)
type MockForumCourseRepository struct{}

func (m *MockForumCourseRepository) CreateCourse(c model.Course) (*model.Course, error) {
	return &model.Course{}, nil
}

func (m *MockForumCourseRepository) GetCourses() ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockForumCourseRepository) GetCourseById(id string) (*model.Course, error) {
	if id == "non-existent-course" {
		return nil, errors.New("course not found")
	}
	if id == "error-course" {
		return nil, errors.New("database error")
	}

	return &model.Course{
		ID:          mustParseForumObjectID("123456789012345678901234"),
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "teacher-123",
		Capacity:    30,
	}, nil
}

func (m *MockForumCourseRepository) DeleteCourse(id string) error {
	return nil
}

func (m *MockForumCourseRepository) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockForumCourseRepository) GetCourseByTitle(title string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockForumCourseRepository) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	return &model.Course{}, nil
}

func (m *MockForumCourseRepository) UpdateStudentsAmount(courseID string, newStudentsAmount int) error {
	return nil
}

func (m *MockForumCourseRepository) CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error) {
	return &model.CourseFeedback{}, nil
}

func (m *MockForumCourseRepository) GetCourseFeedback(courseID string, request schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	return []*model.CourseFeedback{}, nil
}

func (m *MockForumCourseRepository) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockForumCourseRepository) AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

func (m *MockForumCourseRepository) RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

func (m *MockForumCourseRepository) GetCoursesByAuxTeacherId(auxTeacherId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

// Backoffice statistics methods for MockForumCourseRepository
func (m *MockForumCourseRepository) CountCourses() (int64, error) {
	return 2, nil
}

func (m *MockForumCourseRepository) CountActiveCourses() (int64, error) {
	return 1, nil
}

func (m *MockForumCourseRepository) CountFinishedCourses() (int64, error) {
	return 1, nil
}

func (m *MockForumCourseRepository) CountCoursesCreatedThisMonth() (int64, error) {
	return 2, nil
}

func (m *MockForumCourseRepository) CountUniqueTeachers() (int64, error) {
	return 2, nil
}

func (m *MockForumCourseRepository) CountUniqueAuxTeachers() (int64, error) {
	return 3, nil
}

func (m *MockForumCourseRepository) GetTopTeachersByCourseCount(limit int) ([]schemas.CourseDistributionByTeacher, error) {
	return []schemas.CourseDistributionByTeacher{
		{TeacherID: "teacher-1", TeacherName: "Teacher One", CourseCount: 2},
		{TeacherID: "teacher-2", TeacherName: "Teacher Two", CourseCount: 1},
	}, nil
}

func (m *MockForumCourseRepository) GetRecentCourses(limit int) ([]schemas.CourseBasicInfo, error) {
	return []schemas.CourseBasicInfo{
		{ID: "course1", Title: "Test Course 1", TeacherName: "Teacher One", StudentsAmount: 15, Capacity: 20},
		{ID: "course2", Title: "Test Course 2", TeacherName: "Teacher Two", StudentsAmount: 10, Capacity: 15},
	}, nil
}

// Helper function
func mustParseForumObjectID(id string) primitive.ObjectID {
	if len(id) == 24 {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			return objID
		}
	}
	return primitive.NewObjectID()
}

// Test functions

func TestCreateQuestion(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("course-123", "author-123", "Test Title", "Test Description", []model.QuestionTag{model.QuestionTagGeneral})

	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, "course-123", question.CourseID)
	assert.Equal(t, "author-123", question.AuthorID)
	assert.Equal(t, "Test Title", question.Title)
	assert.Equal(t, "Test Description", question.Description)
	assert.Equal(t, []model.QuestionTag{model.QuestionTagGeneral}, question.Tags)
	assert.Equal(t, model.QuestionStatusOpen, question.Status)
}

func TestCreateQuestionWithEmptyCourseID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("", "author-123", "Test Title", "Test Description", []model.QuestionTag{model.QuestionTagGeneral})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "course ID is required", err.Error())
}

func TestCreateQuestionWithEmptyAuthorID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("course-123", "", "Test Title", "Test Description", []model.QuestionTag{model.QuestionTagGeneral})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "author ID is required", err.Error())
}

func TestCreateQuestionWithEmptyTitle(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("course-123", "author-123", "", "Test Description", []model.QuestionTag{model.QuestionTagGeneral})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "title is required", err.Error())
}

func TestCreateQuestionWithEmptyDescription(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("course-123", "author-123", "Test Title", "", []model.QuestionTag{model.QuestionTagGeneral})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "description is required", err.Error())
}

func TestCreateQuestionWithNonExistentCourse(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("non-existent-course", "author-123", "Test Title", "Test Description", []model.QuestionTag{model.QuestionTagGeneral})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "course not found", err.Error())
}

func TestCreateQuestionWithInvalidTags(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.CreateQuestion("course-123", "author-123", "Test Title", "Test Description", []model.QuestionTag{"invalid-tag"})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Contains(t, err.Error(), "invalid tag")
}

func TestGetQuestionById(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.GetQuestionById("question-123")

	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, "Test Question", question.Title)
}

func TestGetQuestionByIdWithEmptyID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.GetQuestionById("")

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestGetQuestionByIdWithNonExistentID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.GetQuestionById("non-existent-question")

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "question not found", err.Error())
}

func TestGetQuestionsByCourseId(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.GetQuestionsByCourseId("course-123")

	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 2)
	assert.Equal(t, "Question 1", questions[0].Title)
	assert.Equal(t, "Question 2", questions[1].Title)
}

func TestGetQuestionsByCourseIdWithEmptyID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.GetQuestionsByCourseId("")

	assert.Error(t, err)
	assert.Nil(t, questions)
	assert.Equal(t, "course ID is required", err.Error())
}

func TestGetQuestionsByCourseIdWithNonExistentCourse(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.GetQuestionsByCourseId("non-existent-course")

	assert.Error(t, err)
	assert.Nil(t, questions)
	assert.Equal(t, "course not found", err.Error())
}

func TestUpdateQuestion(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.UpdateQuestion("question-123", "Updated Title", "Updated Description", []model.QuestionTag{model.QuestionTagPractica})

	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, "Updated Title", question.Title)
	assert.Equal(t, "Updated Description", question.Description)
	assert.Equal(t, []model.QuestionTag{model.QuestionTagPractica}, question.Tags)
}

func TestUpdateQuestionWithEmptyID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.UpdateQuestion("", "Updated Title", "Updated Description", []model.QuestionTag{model.QuestionTagPractica})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestUpdateQuestionWithNoFields(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	question, err := forumService.UpdateQuestion("question-123", "", "", []model.QuestionTag{})

	assert.Error(t, err)
	assert.Nil(t, question)
	assert.Equal(t, "at least one field must be provided for update", err.Error())
}

func TestDeleteQuestion(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteQuestion("question-author-123", "author-123")

	assert.NoError(t, err)
}

func TestDeleteQuestionWithEmptyID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteQuestion("", "author-123")

	assert.Error(t, err)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestDeleteQuestionWithEmptyAuthorID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteQuestion("question-123", "")

	assert.Error(t, err)
	assert.Equal(t, "author ID is required", err.Error())
}

func TestDeleteQuestionWithWrongAuthor(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteQuestion("question-author-123", "wrong-author")

	assert.Error(t, err)
	assert.Equal(t, "you can only delete your own questions", err.Error())
}

func TestAddAnswer(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.AddAnswer("question-123", "author-123", "Test answer content")

	assert.NoError(t, err)
	assert.NotNil(t, answer)
	assert.Equal(t, "author-123", answer.AuthorID)
	assert.Equal(t, "Test answer content", answer.Content)
	assert.Equal(t, "new-answer-id", answer.ID)
}

func TestAddAnswerWithEmptyQuestionID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.AddAnswer("", "author-123", "Test answer content")

	assert.Error(t, err)
	assert.Nil(t, answer)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestAddAnswerWithEmptyAuthorID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.AddAnswer("question-123", "", "Test answer content")

	assert.Error(t, err)
	assert.Nil(t, answer)
	assert.Equal(t, "author ID is required", err.Error())
}

func TestAddAnswerWithEmptyContent(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.AddAnswer("question-123", "author-123", "")

	assert.Error(t, err)
	assert.Nil(t, answer)
	assert.Equal(t, "content is required", err.Error())
}

func TestAddAnswerWithNonExistentQuestion(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.AddAnswer("non-existent-question", "author-123", "Test answer content")

	assert.Error(t, err)
	assert.Nil(t, answer)
	assert.Equal(t, "question not found", err.Error())
}

func TestUpdateAnswer(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.UpdateAnswer("question-with-answers", "answer-123", "answer-author-123", "Updated content")

	assert.NoError(t, err)
	assert.NotNil(t, answer)
	assert.Equal(t, "Updated content", answer.Content)
}

func TestUpdateAnswerWithWrongAuthor(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	answer, err := forumService.UpdateAnswer("question-with-answers", "answer-123", "wrong-author", "Updated content")

	assert.Error(t, err)
	assert.Nil(t, answer)
	assert.Equal(t, "you can only update your own answers", err.Error())
}

func TestDeleteAnswer(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteAnswer("question-with-answers", "answer-123", "answer-author-123")

	assert.NoError(t, err)
}

func TestDeleteAnswerWithWrongAuthor(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.DeleteAnswer("question-with-answers", "answer-123", "wrong-author")

	assert.Error(t, err)
	assert.Equal(t, "you can only delete your own answers", err.Error())
}

func TestAcceptAnswer(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.AcceptAnswer("question-with-answers", "answer-123", "author-123")

	assert.NoError(t, err)
}

func TestAcceptAnswerWithWrongAuthor(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.AcceptAnswer("question-with-answers", "answer-123", "wrong-author")

	assert.Error(t, err)
	assert.Equal(t, "only the question author can accept answers", err.Error())
}

func TestVoteQuestion(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.VoteQuestion("question-123", "voter-123", model.VoteTypeUp)

	assert.NoError(t, err)
}

func TestVoteQuestionSelfVote(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.VoteQuestion("question-123", "author-123", model.VoteTypeUp)

	assert.Error(t, err)
	assert.Equal(t, "you cannot vote on your own question", err.Error())
}

func TestVoteAnswer(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.VoteAnswer("question-with-answers", "answer-123", "voter-123", model.VoteTypeUp)

	assert.NoError(t, err)
}

func TestVoteAnswerSelfVote(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	err := forumService.VoteAnswer("question-with-answers", "answer-123", "answer-author-123", model.VoteTypeUp)

	assert.Error(t, err)
	assert.Equal(t, "you cannot vote on your own answer", err.Error())
}

func TestSearchQuestions(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("course-123", "", []model.QuestionTag{}, "")

	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 2)
}

func TestSearchQuestionsWithQuery(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("course-123", "architecture", []model.QuestionTag{}, "")

	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 1)
	assert.Equal(t, "Architecture Question", questions[0].Title)
}

func TestSearchQuestionsWithTags(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("course-123", "", []model.QuestionTag{model.QuestionTagPractica}, "")

	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 1)
	assert.Equal(t, "Database Question", questions[0].Title)
}

func TestSearchQuestionsWithStatus(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("course-123", "", []model.QuestionTag{}, model.QuestionStatusResolved)

	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 1)
	assert.Equal(t, "Database Question", questions[0].Title)
}

func TestSearchQuestionsWithEmptyCourseID(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("", "", []model.QuestionTag{}, "")

	assert.Error(t, err)
	assert.Nil(t, questions)
	assert.Equal(t, "course ID is required", err.Error())
}

func TestSearchQuestionsWithNonExistentCourse(t *testing.T) {
	forumRepo := &MockForumRepository{}
	courseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(forumRepo, courseRepo)

	questions, err := forumService.SearchQuestions("non-existent-course", "", []model.QuestionTag{}, "")

	assert.Error(t, err)
	assert.Nil(t, questions)
	assert.Equal(t, "course not found", err.Error())
}

// Tests for RemoveVoteFromQuestion
func TestRemoveVoteFromQuestion(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromQuestion("valid-question", "user-123")
	assert.NoError(t, err)
}

func TestRemoveVoteFromQuestionWithEmptyQuestionID(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromQuestion("", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestRemoveVoteFromQuestionWithEmptyUserID(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromQuestion("valid-question", "")
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())
}

func TestRemoveVoteFromQuestionWithNonExistentQuestion(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromQuestion("non-existent-question", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "question not found", err.Error())
}

// Tests for RemoveVoteFromAnswer
func TestRemoveVoteFromAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("question-with-answers", "answer-123", "user-123")
	assert.NoError(t, err)
}

func TestRemoveVoteFromAnswerWithEmptyQuestionID(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("", "answer-123", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "question ID is required", err.Error())
}

func TestRemoveVoteFromAnswerWithEmptyAnswerID(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("valid-question", "", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "answer ID is required", err.Error())
}

func TestRemoveVoteFromAnswerWithEmptyUserID(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("valid-question", "answer-123", "")
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())
}

func TestRemoveVoteFromAnswerWithNonExistentQuestion(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("non-existent-question", "answer-123", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "question not found", err.Error())
}

func TestRemoveVoteFromAnswerWithNonExistentAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.RemoveVoteFromAnswer("question-with-answers", "non-existent-answer", "user-123")
	assert.Error(t, err)
	assert.Equal(t, "answer not found", err.Error())
}

// Additional tests for better coverage of existing functions
func TestVoteQuestionWithInvalidVoteType(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.VoteQuestion("valid-question", "user-123", 99)
	assert.Error(t, err)
	assert.Equal(t, "invalid vote type", err.Error())
}

func TestVoteAnswerWithInvalidVoteType(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.VoteAnswer("question-with-answers", "answer-123", "user-123", 99)
	assert.Error(t, err)
	assert.Equal(t, "invalid vote type", err.Error())
}

func TestVoteAnswerWithNonExistentAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.VoteAnswer("question-with-answers", "non-existent-answer", "user-123", model.VoteTypeUp)
	assert.Error(t, err)
	assert.Equal(t, "answer not found", err.Error())
}

func TestUpdateQuestionWithInvalidTags(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	invalidTags := []model.QuestionTag{"invalid-tag"}
	_, err := forumService.UpdateQuestion("valid-question", "New Title", "New Description", invalidTags)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag")
}

func TestUpdateAnswerWithNonExistentAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	_, err := forumService.UpdateAnswer("question-with-answers", "non-existent-answer", "author-123", "Updated content")
	assert.Error(t, err)
	assert.Equal(t, "answer not found", err.Error())
}

func TestDeleteAnswerWithNonExistentAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.DeleteAnswer("question-with-answers", "non-existent-answer", "author-123")
	assert.Error(t, err)
	assert.Equal(t, "answer not found", err.Error())
}

func TestAcceptAnswerWithNonExistentAnswer(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	err := forumService.AcceptAnswer("question-with-answers", "non-existent-answer", "author-123")
	assert.Error(t, err)
	assert.Equal(t, "answer not found", err.Error())
}

func TestSearchQuestionsWithInvalidTags(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	invalidTags := []model.QuestionTag{"invalid-tag"}
	_, err := forumService.SearchQuestions("course-123", "", invalidTags, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag")
}

func TestSearchQuestionsWithInvalidStatus(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	_, err := forumService.SearchQuestions("course-123", "", nil, "invalid-status")
	assert.Error(t, err)
	assert.Equal(t, "invalid question status", err.Error())
}

func TestValidateTagsWithEmptyTags(t *testing.T) {
	mockForumRepo := &MockForumRepository{}
	mockCourseRepo := &MockForumCourseRepository{}
	forumService := service.NewForumService(mockForumRepo, mockCourseRepo)

	// Test validateTags with empty tags by calling CreateQuestion with empty tags
	_, err := forumService.CreateQuestion("course-123", "author-123", "Title", "Description", []model.QuestionTag{})
	assert.NoError(t, err)
}
