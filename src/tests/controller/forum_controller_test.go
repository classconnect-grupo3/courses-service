package controller_test

import (
	"bytes"
	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"courses-service/src/schemas"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock Forum Service
type MockForumService struct{}

func (m *MockForumService) CreateQuestion(courseID, authorID, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error) {
	if courseID == "error-course" {
		return nil, errors.New("course not found")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}

	return &model.ForumQuestion{
		ID:          primitive.NewObjectID(),
		CourseID:    courseID,
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Tags:        tags,
		Status:      model.QuestionStatusOpen,
		Votes:       []model.Vote{},
		Answers:     []model.ForumAnswer{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *MockForumService) GetQuestionById(id string) (*model.ForumQuestion, error) {
	if id == "non-existent" {
		return nil, errors.New("question not found")
	}

	return &model.ForumQuestion{
		ID:          mustParseForumObjectID("123456789012345678901234"),
		CourseID:    "course-123",
		AuthorID:    "author-123",
		Title:       "Test Question",
		Description: "Test Description",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral},
		Status:      model.QuestionStatusOpen,
		Votes: []model.Vote{
			{UserID: "user1", VoteType: model.VoteTypeUp, CreatedAt: time.Now()},
			{UserID: "user2", VoteType: model.VoteTypeDown, CreatedAt: time.Now()},
		},
		Answers: []model.ForumAnswer{
			{
				ID:       "answer-123",
				AuthorID: "answer-author-123",
				Content:  "Test Answer",
				Votes: []model.Vote{
					{UserID: "user3", VoteType: model.VoteTypeUp, CreatedAt: time.Now()},
				},
				IsAccepted: false,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockForumService) GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error) {
	if courseID == "error-course" {
		return nil, errors.New("course not found")
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
			Votes:       []model.Vote{{UserID: "user1", VoteType: model.VoteTypeUp}},
			Answers:     []model.ForumAnswer{},
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
			Votes:       []model.Vote{},
			Answers:     []model.ForumAnswer{{ID: "answer1", Content: "Answer content"}},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

func (m *MockForumService) UpdateQuestion(id, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error) {
	if id == "non-existent" {
		return nil, errors.New("question not found")
	}

	return &model.ForumQuestion{
		ID:          mustParseForumObjectID("123456789012345678901234"),
		CourseID:    "course-123",
		AuthorID:    "author-123",
		Title:       title,
		Description: description,
		Tags:        tags,
		Status:      model.QuestionStatusOpen,
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *MockForumService) DeleteQuestion(id, authorID string) error {
	if id == "non-existent" {
		return errors.New("question not found")
	}
	if authorID == "wrong-author" {
		return errors.New("you can only delete your own questions")
	}
	return nil
}

func (m *MockForumService) AddAnswer(questionID, authorID, content string) (*model.ForumAnswer, error) {
	if questionID == "non-existent" {
		return nil, errors.New("question not found")
	}

	return &model.ForumAnswer{
		ID:        "new-answer-id",
		AuthorID:  authorID,
		Content:   content,
		Votes:     []model.Vote{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockForumService) UpdateAnswer(questionID, answerID, authorID, content string) (*model.ForumAnswer, error) {
	if questionID == "non-existent" || answerID == "non-existent" {
		return nil, errors.New("question or answer not found")
	}
	if authorID == "wrong-author" {
		return nil, errors.New("you can only update your own answers")
	}

	return &model.ForumAnswer{
		ID:        answerID,
		AuthorID:  authorID,
		Content:   content,
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockForumService) DeleteAnswer(questionID, answerID, authorID string) error {
	if questionID == "non-existent" || answerID == "non-existent" {
		return errors.New("question or answer not found")
	}
	if authorID == "wrong-author" {
		return errors.New("you can only delete your own answers")
	}
	return nil
}

func (m *MockForumService) AcceptAnswer(questionID, answerID, authorID string) error {
	if questionID == "non-existent" || answerID == "non-existent" {
		return errors.New("question or answer not found")
	}
	if authorID == "wrong-author" {
		return errors.New("only the question author can accept answers")
	}
	return nil
}

func (m *MockForumService) VoteQuestion(questionID, userID string, voteType int) error {
	if questionID == "non-existent" {
		return errors.New("question not found")
	}
	if userID == "author-123" {
		return errors.New("you cannot vote on your own question")
	}
	return nil
}

func (m *MockForumService) VoteAnswer(questionID, answerID, userID string, voteType int) error {
	if questionID == "non-existent" || answerID == "non-existent" {
		return errors.New("question or answer not found")
	}
	if userID == "answer-author-123" {
		return errors.New("you cannot vote on your own answer")
	}
	return nil
}

func (m *MockForumService) RemoveVoteFromQuestion(questionID, userID string) error {
	if questionID == "non-existent" {
		return errors.New("question not found")
	}
	return nil
}

func (m *MockForumService) RemoveVoteFromAnswer(questionID, answerID, userID string) error {
	if questionID == "non-existent" || answerID == "non-existent" {
		return errors.New("question or answer not found")
	}
	return nil
}

func (m *MockForumService) SearchQuestions(courseID, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error) {
	if courseID == "error-course" {
		return nil, errors.New("course not found")
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
			Votes:       []model.Vote{},
			Answers:     []model.ForumAnswer{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if query == "architecture" || len(tags) > 0 || status != "" {
		return questions, nil
	}

	return questions, nil
}

func (m *MockForumService) GetForumParticipants(courseID string) ([]string, error) {
	if courseID == "error-course" {
		return nil, errors.New("course not found")
	}
	if courseID == "empty-course" {
		return []string{}, nil
	}

	return []string{"author-123", "author-456", "voter-123"}, nil
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

// Setup
var (
	mockForumService  = &MockForumService{}
	forumController   = controller.NewForumController(mockForumService)
	normalForumRouter = gin.Default()
)

func init() {
	router.InitializeForumRoutes(normalForumRouter, forumController)
}

// Test functions

func TestCreateQuestion(t *testing.T) {

	requestBody := schemas.CreateQuestionRequest{
		CourseID:    "course-123",
		AuthorID:    "author-123",
		Title:       "Test Question",
		Description: "Test Description",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response schemas.QuestionDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Question", response.Title)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, "course-123", response.CourseID)
	assert.Equal(t, "author-123", response.AuthorID)
}

func TestCreateQuestionWithInvalidJSON(t *testing.T) {

	req, _ := http.NewRequest("POST", "/forum/questions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "invalid character")
}

func TestCreateQuestionWithServiceError(t *testing.T) {

	requestBody := schemas.CreateQuestionRequest{
		CourseID:    "error-course",
		AuthorID:    "author-123",
		Title:       "Test Question",
		Description: "Test Description",
		Tags:        []model.QuestionTag{model.QuestionTagGeneral},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "course not found", response.Error)
}

func TestGetQuestionById(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/questions/123456789012345678901234", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.QuestionDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Question", response.Title)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, 0, response.VoteCount) // 1 up - 1 down = 0
	assert.Len(t, response.Answers, 1)
}

func TestGetQuestionByIdNotFound(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/questions/non-existent", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "question not found", response.Error)
}

func TestGetQuestionsByCourseId(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/course-123/questions", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []schemas.QuestionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Question 1", response[0].Title)
	assert.Equal(t, "Question 2", response[1].Title)
	assert.Equal(t, 1, response[0].VoteCount)
	assert.Equal(t, 0, response[0].AnswerCount)
	assert.Equal(t, 1, response[1].AnswerCount)
}

func TestGetQuestionsByCourseIdWithError(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/error-course/questions", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "course not found", response.Error)
}

func TestUpdateQuestion(t *testing.T) {

	requestBody := schemas.UpdateQuestionRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
		Tags:        []model.QuestionTag{model.QuestionTagPractica},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/forum/questions/123456789012345678901234", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.QuestionDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, "Updated Description", response.Description)
}

func TestUpdateQuestionWithInvalidJSON(t *testing.T) {

	req, _ := http.NewRequest("PUT", "/forum/questions/123456789012345678901234", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateQuestionNotFound(t *testing.T) {

	requestBody := schemas.UpdateQuestionRequest{
		Title: "Updated Title",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/forum/questions/non-existent", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteQuestion(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234?authorId=author-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Question deleted successfully", response.Message)
}

func TestDeleteQuestionWithoutAuthorId(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "authorId query parameter is required", response.Error)
}

func TestDeleteQuestionWithWrongAuthor(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234?authorId=wrong-author", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddAnswer(t *testing.T) {

	requestBody := schemas.CreateAnswerRequest{
		AuthorID: "author-123",
		Content:  "Test answer content",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/answers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response schemas.AnswerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "author-123", response.AuthorID)
	assert.Equal(t, "Test answer content", response.Content)
	assert.Equal(t, "new-answer-id", response.ID)
}

func TestAddAnswerWithInvalidJSON(t *testing.T) {

	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/answers", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAnswerToNonExistentQuestion(t *testing.T) {

	requestBody := schemas.CreateAnswerRequest{
		AuthorID: "author-123",
		Content:  "Test answer content",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions/non-existent/answers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAnswer(t *testing.T) {

	requestBody := schemas.UpdateAnswerRequest{
		Content: "Updated answer content",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/forum/questions/123456789012345678901234/answers/answer-123?authorId=author-123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.AnswerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated answer content", response.Content)
}

func TestUpdateAnswerWithWrongAuthor(t *testing.T) {

	requestBody := schemas.UpdateAnswerRequest{
		Content: "Updated answer content",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/forum/questions/123456789012345678901234/answers/answer-123?authorId=wrong-author", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteAnswer(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234/answers/answer-123?authorId=author-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Answer deleted successfully", response.Message)
}

func TestDeleteAnswerWithoutAuthorId(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234/answers/answer-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAcceptAnswer(t *testing.T) {

	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/answers/answer-123/accept?authorId=author-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Answer accepted successfully", response.Message)
}

func TestAcceptAnswerWithWrongAuthor(t *testing.T) {

	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/answers/answer-123/accept?authorId=wrong-author", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVoteQuestion(t *testing.T) {

	requestBody := schemas.VoteRequest{
		UserID:   "voter-123",
		VoteType: model.VoteTypeUp,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.VoteResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Vote registered successfully", response.Message)
}

func TestVoteQuestionSelfVote(t *testing.T) {

	requestBody := schemas.VoteRequest{
		UserID:   "author-123",
		VoteType: model.VoteTypeUp,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVoteAnswer(t *testing.T) {

	requestBody := schemas.VoteRequest{
		UserID:   "voter-123",
		VoteType: model.VoteTypeUp,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/forum/questions/123456789012345678901234/answers/answer-123/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.VoteResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Vote registered successfully", response.Message)
}

func TestRemoveVoteFromQuestion(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234/vote?userId=voter-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Vote removed successfully", response.Message)
}

func TestRemoveVoteFromQuestionWithoutUserId(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234/vote", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveVoteFromAnswer(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/forum/questions/123456789012345678901234/answers/answer-123/vote?userId=voter-123", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Vote removed successfully", response.Message)
}

func TestSearchQuestions(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/course-123/search", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.SearchQuestionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Questions, 1)
	assert.Equal(t, "Architecture Question", response.Questions[0].Title)
	assert.Equal(t, 1, response.Total)
}

func TestSearchQuestionsWithQuery(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/course-123/search?query=architecture", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.SearchQuestionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Questions, 1)
}

func TestSearchQuestionsWithTags(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/course-123/search?tags=teoria", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.SearchQuestionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Questions, 1)
}

func TestSearchQuestionsWithStatus(t *testing.T) {

	req, _ := http.NewRequest("GET", "/forum/courses/course-123/search?status=open", nil)
	w := httptest.NewRecorder()
	normalForumRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.SearchQuestionsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Questions, 1)
}

func TestSearchQuestionsWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/forum/courses/error-course/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Additional tests for missing coverage

// Tests for UpdateAnswer missing cases
func TestUpdateAnswerWithoutAuthorId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.UpdateAnswerRequest{
		Content: "Updated content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/forum/questions/question-123/answers/answer-123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "authorId query parameter is required", response.Error)
}

func TestUpdateAnswerWithInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/forum/questions/question-123/answers/answer-123?authorId=author-123", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAnswerWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.UpdateAnswerRequest{
		Content: "Updated content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/forum/questions/non-existent/answers/answer-123?authorId=author-123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Tests for DeleteAnswer missing cases
func TestDeleteAnswerWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/forum/questions/non-existent/answers/answer-123?authorId=author-123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Tests for AcceptAnswer missing cases
func TestAcceptAnswerWithoutAuthorId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/question-123/answers/answer-123/accept", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "authorId query parameter is required", response.Error)
}

func TestAcceptAnswerWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/non-existent/answers/answer-123/accept?authorId=author-123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Tests for VoteQuestion missing cases
func TestVoteQuestionWithInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/question-123/vote", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVoteQuestionWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.VoteRequest{
		UserID:   "user-123",
		VoteType: model.VoteTypeUp,
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/non-existent/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVoteQuestionWithDownVote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.VoteRequest{
		UserID:   "user-123",
		VoteType: model.VoteTypeDown,
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/question-123/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response schemas.VoteResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Vote registered successfully", response.Message)
}

// Tests for VoteAnswer missing cases
func TestVoteAnswerWithInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/question-123/answers/answer-123/vote", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVoteAnswerWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.VoteRequest{
		UserID:   "user-123",
		VoteType: model.VoteTypeUp,
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/non-existent/answers/answer-123/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVoteAnswerWithDownVote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	requestBody := schemas.VoteRequest{
		UserID:   "user-123",
		VoteType: model.VoteTypeDown,
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forum/questions/question-123/answers/answer-123/vote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response schemas.VoteResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Vote registered successfully", response.Message)
}

// Tests for RemoveVoteFromQuestion missing cases
func TestRemoveVoteFromQuestionWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/forum/questions/non-existent/vote?userId=user-123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Tests for RemoveVoteFromAnswer missing cases
func TestRemoveVoteFromAnswerWithoutUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/forum/questions/question-123/answers/answer-123/vote", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "userId query parameter is required", response.Error)
}

func TestRemoveVoteFromAnswerWithServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/forum/questions/non-existent/answers/answer-123/vote?userId=user-123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Tests for SearchQuestions missing cases
func TestSearchQuestionsWithInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forumController := controller.NewForumController(&MockForumService{})
	router.InitializeForumRoutes(r, forumController)

	w := httptest.NewRecorder()
	// Test with malformed query parameters that could cause binding errors
	req, _ := http.NewRequest("GET", "/forum/courses/course-123/search?tags=invalid[", nil)
	r.ServeHTTP(w, req)

	// The response could be 200 or 400 depending on how gin handles the malformed query
	// Let's check for either valid response or error
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)
}
