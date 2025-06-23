package controller

import (
	"log/slog"
	"net/http"

	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"

	"github.com/gin-gonic/gin"
)

type ForumController struct {
	service service.ForumServiceInterface
}

func NewForumController(service service.ForumServiceInterface) *ForumController {
	return &ForumController{service: service}
}

// Question endpoints

// @Summary Create a new question
// @Description Create a new question in the forum for a specific course
// @Tags forum
// @Accept json
// @Produce json
// @Param question body schemas.CreateQuestionRequest true "Question to create"
// @Success 201 {object} schemas.QuestionDetailResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions [post]
func (c *ForumController) CreateQuestion(ctx *gin.Context) {
	slog.Debug("Creating forum question")

	var request schemas.CreateQuestionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	question, err := c.service.CreateQuestion(
		request.CourseID,
		request.AuthorID,
		request.Title,
		request.Description,
		request.Tags,
	)
	if err != nil {
		slog.Error("Error creating question", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	response := c.mapQuestionToDetailResponse(question)
	slog.Debug("Question created", "question_id", question.ID.Hex())
	ctx.JSON(http.StatusCreated, response)
}

// @Summary Get question by ID
// @Description Get a specific question by its ID with all answers
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Success 200 {object} schemas.QuestionDetailResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId} [get]
func (c *ForumController) GetQuestionById(ctx *gin.Context) {
	slog.Debug("Getting question by ID")

	id := ctx.Param("questionId")
	question, err := c.service.GetQuestionById(id)
	if err != nil {
		slog.Error("Error getting question by ID", "error", err)
		ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	response := c.mapQuestionToDetailResponse(question)
	slog.Debug("Question retrieved", "question_id", id)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Get questions by course ID
// @Description Get all questions for a specific course
// @Tags forum
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Success 200 {array} schemas.QuestionResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/courses/{courseId}/questions [get]
func (c *ForumController) GetQuestionsByCourseId(ctx *gin.Context) {
	slog.Debug("Getting questions by course ID")

	courseID := ctx.Param("courseId")
	questions, err := c.service.GetQuestionsByCourseId(courseID)
	if err != nil {
		slog.Error("Error getting questions by course ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	var responses []schemas.QuestionResponse
	for _, question := range questions {
		responses = append(responses, c.mapQuestionToResponse(&question))
	}

	slog.Debug("Questions retrieved", "course_id", courseID, "count", len(responses))
	ctx.JSON(http.StatusOK, responses)
}

// @Summary Update a question
// @Description Update a question's title, description, or tags
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param question body schemas.UpdateQuestionRequest true "Question update data"
// @Success 200 {object} schemas.QuestionDetailResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId} [put]
func (c *ForumController) UpdateQuestion(ctx *gin.Context) {
	slog.Debug("Updating question")

	id := ctx.Param("questionId")
	var request schemas.UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	question, err := c.service.UpdateQuestion(id, request.Title, request.Description, request.Tags)
	if err != nil {
		slog.Error("Error updating question", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	response := c.mapQuestionToDetailResponse(question)
	slog.Debug("Question updated", "question_id", id)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Delete a question
// @Description Delete a question (only by the author)
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param authorId query string true "Author ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId} [delete]
func (c *ForumController) DeleteQuestion(ctx *gin.Context) {
	slog.Debug("Deleting question")

	id := ctx.Param("questionId")
	authorID := ctx.Query("authorId")

	if authorID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "authorId query parameter is required"})
		return
	}

	err := c.service.DeleteQuestion(id, authorID)
	if err != nil {
		slog.Error("Error deleting question", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Debug("Question deleted", "question_id", id)
	ctx.JSON(http.StatusOK, schemas.MessageResponse{Message: "Question deleted successfully"})
}

// Answer endpoints

// @Summary Add an answer to a question
// @Description Add a new answer to a specific question
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answer body schemas.CreateAnswerRequest true "Answer to create"
// @Success 201 {object} schemas.AnswerResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers [post]
func (c *ForumController) AddAnswer(ctx *gin.Context) {
	slog.Debug("Adding answer to question")

	questionID := ctx.Param("questionId")
	var request schemas.CreateAnswerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	answer, err := c.service.AddAnswer(questionID, request.AuthorID, request.Content)
	if err != nil {
		slog.Error("Error adding answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	response := c.mapAnswerToResponse(answer)
	slog.Debug("Answer added", "question_id", questionID, "answer_id", answer.ID)
	ctx.JSON(http.StatusCreated, response)
}

// @Summary Update an answer
// @Description Update an answer's content (only by the author)
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answerId path string true "Answer ID"
// @Param answer body schemas.UpdateAnswerRequest true "Answer update data"
// @Param authorId query string true "Author ID"
// @Success 200 {object} schemas.AnswerResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers/{answerId} [put]
func (c *ForumController) UpdateAnswer(ctx *gin.Context) {
	slog.Debug("Updating answer")

	questionID := ctx.Param("questionId")
	answerID := ctx.Param("answerId")
	authorID := ctx.Query("authorId")

	if authorID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "authorId query parameter is required"})
		return
	}

	var request schemas.UpdateAnswerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	answer, err := c.service.UpdateAnswer(questionID, answerID, authorID, request.Content)
	if err != nil {
		slog.Error("Error updating answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	response := c.mapAnswerToResponse(answer)
	slog.Debug("Answer updated", "question_id", questionID, "answer_id", answerID)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Delete an answer
// @Description Delete an answer (only by the author)
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answerId path string true "Answer ID"
// @Param authorId query string true "Author ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers/{answerId} [delete]
func (c *ForumController) DeleteAnswer(ctx *gin.Context) {
	slog.Debug("Deleting answer")

	questionID := ctx.Param("questionId")
	answerID := ctx.Param("answerId")
	authorID := ctx.Query("authorId")

	if authorID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "authorId query parameter is required"})
		return
	}

	err := c.service.DeleteAnswer(questionID, answerID, authorID)
	if err != nil {
		slog.Error("Error deleting answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Debug("Answer deleted", "question_id", questionID, "answer_id", answerID)
	ctx.JSON(http.StatusOK, schemas.MessageResponse{Message: "Answer deleted successfully"})
}

// @Summary Accept an answer
// @Description Accept an answer as the solution (only by the question author)
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answerId path string true "Answer ID"
// @Param authorId query string true "Question Author ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers/{answerId}/accept [post]
func (c *ForumController) AcceptAnswer(ctx *gin.Context) {
	slog.Debug("Accepting answer")

	questionID := ctx.Param("questionId")
	answerID := ctx.Param("answerId")
	authorID := ctx.Query("authorId")

	if authorID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "authorId query parameter is required"})
		return
	}

	err := c.service.AcceptAnswer(questionID, answerID, authorID)
	if err != nil {
		slog.Error("Error accepting answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Debug("Answer accepted", "question_id", questionID, "answer_id", answerID)
	ctx.JSON(http.StatusOK, schemas.MessageResponse{Message: "Answer accepted successfully"})
}

// Vote endpoints

// @Summary Vote on a question
// @Description Vote up or down on a question
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param vote body schemas.VoteRequest true "Vote data"
// @Success 200 {object} schemas.VoteResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/vote [post]
func (c *ForumController) VoteQuestion(ctx *gin.Context) {
	slog.Debug("Voting on question")

	questionID := ctx.Param("questionId")
	var request schemas.VoteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	err := c.service.VoteQuestion(questionID, request.UserID, request.VoteType)
	if err != nil {
		slog.Error("Error voting on question", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	voteTypeStr := "up"
	if request.VoteType == model.VoteTypeDown {
		voteTypeStr = "down"
	}

	slog.Debug("Vote registered", "question_id", questionID, "vote_type", voteTypeStr)
	ctx.JSON(http.StatusOK, schemas.VoteResponse{Message: "Vote registered successfully"})
}

// @Summary Vote on an answer
// @Description Vote up or down on an answer
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answerId path string true "Answer ID"
// @Param vote body schemas.VoteRequest true "Vote data"
// @Success 200 {object} schemas.VoteResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 403 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers/{answerId}/vote [post]
func (c *ForumController) VoteAnswer(ctx *gin.Context) {
	slog.Debug("Voting on answer")

	questionID := ctx.Param("questionId")
	answerID := ctx.Param("answerId")
	var request schemas.VoteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	err := c.service.VoteAnswer(questionID, answerID, request.UserID, request.VoteType)
	if err != nil {
		slog.Error("Error voting on answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	voteTypeStr := "up"
	if request.VoteType == model.VoteTypeDown {
		voteTypeStr = "down"
	}

	slog.Debug("Vote registered", "question_id", questionID, "answer_id", answerID, "vote_type", voteTypeStr)
	ctx.JSON(http.StatusOK, schemas.VoteResponse{Message: "Vote registered successfully"})
}

// @Summary Remove vote from question
// @Description Remove a user's vote from a question
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param userId query string true "User ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/vote [delete]
func (c *ForumController) RemoveVoteFromQuestion(ctx *gin.Context) {
	slog.Debug("Removing vote from question")

	questionID := ctx.Param("questionId")
	userID := ctx.Query("userId")

	if userID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "userId query parameter is required"})
		return
	}

	err := c.service.RemoveVoteFromQuestion(questionID, userID)
	if err != nil {
		slog.Error("Error removing vote from question", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Debug("Vote removed", "question_id", questionID, "user_id", userID)
	ctx.JSON(http.StatusOK, schemas.MessageResponse{Message: "Vote removed successfully"})
}

// @Summary Remove vote from answer
// @Description Remove a user's vote from an answer
// @Tags forum
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID"
// @Param answerId path string true "Answer ID"
// @Param userId query string true "User ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/questions/{questionId}/answers/{answerId}/vote [delete]
func (c *ForumController) RemoveVoteFromAnswer(ctx *gin.Context) {
	slog.Debug("Removing vote from answer")

	questionID := ctx.Param("questionId")
	answerID := ctx.Param("answerId")
	userID := ctx.Query("userId")

	if userID == "" {
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: "userId query parameter is required"})
		return
	}

	err := c.service.RemoveVoteFromAnswer(questionID, answerID, userID)
	if err != nil {
		slog.Error("Error removing vote from answer", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Debug("Vote removed", "question_id", questionID, "answer_id", answerID, "user_id", userID)
	ctx.JSON(http.StatusOK, schemas.MessageResponse{Message: "Vote removed successfully"})
}

// Search endpoints

// @Summary Search questions
// @Description Search questions in a course with optional filters
// @Tags forum
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Param query query string false "Search query"
// @Param tags query []string false "Filter by tags"
// @Param status query string false "Filter by status"
// @Success 200 {object} schemas.SearchQuestionsResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /forum/courses/{courseId}/search [get]
func (c *ForumController) SearchQuestions(ctx *gin.Context) {
	slog.Debug("Searching questions")

	courseID := ctx.Param("courseId")

	var request schemas.SearchQuestionsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		slog.Error("Error binding query parameters", "error", err)
		ctx.JSON(http.StatusBadRequest, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	questions, err := c.service.SearchQuestions(courseID, request.Query, request.Tags, request.Status)
	if err != nil {
		slog.Error("Error searching questions", "error", err)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{Error: err.Error()})
		return
	}

	var questionResponses []schemas.QuestionResponse
	for _, question := range questions {
		questionResponses = append(questionResponses, c.mapQuestionToResponse(&question))
	}

	response := schemas.SearchQuestionsResponse{
		Questions: questionResponses,
		Total:     len(questionResponses),
	}

	slog.Debug("Questions searched", "course_id", courseID, "total", response.Total)
	ctx.JSON(http.StatusOK, response)
}

// Helper methods for mapping models to responses

func (c *ForumController) mapQuestionToResponse(question *model.ForumQuestion) schemas.QuestionResponse {
	voteCount := c.calculateVoteCount(question.Votes)
	answerCount := len(question.Answers)

	return schemas.QuestionResponse{
		ID:               question.ID.Hex(),
		CourseID:         question.CourseID,
		AuthorID:         question.AuthorID,
		Title:            question.Title,
		Description:      question.Description,
		Tags:             question.Tags,
		Votes:            question.Votes,
		VoteCount:        voteCount,
		AnswerCount:      answerCount,
		Status:           question.Status,
		AcceptedAnswerID: question.AcceptedAnswerID,
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}
}

func (c *ForumController) mapQuestionToDetailResponse(question *model.ForumQuestion) schemas.QuestionDetailResponse {
	voteCount := c.calculateVoteCount(question.Votes)

	var answers []schemas.AnswerResponse
	for _, answer := range question.Answers {
		answers = append(answers, c.mapAnswerToResponse(&answer))
	}

	return schemas.QuestionDetailResponse{
		ID:               question.ID.Hex(),
		CourseID:         question.CourseID,
		AuthorID:         question.AuthorID,
		Title:            question.Title,
		Description:      question.Description,
		Tags:             question.Tags,
		Votes:            question.Votes,
		VoteCount:        voteCount,
		Answers:          answers,
		Status:           question.Status,
		AcceptedAnswerID: question.AcceptedAnswerID,
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}
}

func (c *ForumController) mapAnswerToResponse(answer *model.ForumAnswer) schemas.AnswerResponse {
	voteCount := c.calculateVoteCount(answer.Votes)

	return schemas.AnswerResponse{
		ID:         answer.ID,
		AuthorID:   answer.AuthorID,
		Content:    answer.Content,
		Votes:      answer.Votes,
		VoteCount:  voteCount,
		IsAccepted: answer.IsAccepted,
		CreatedAt:  answer.CreatedAt,
		UpdatedAt:  answer.UpdatedAt,
	}
}

func (c *ForumController) calculateVoteCount(votes []model.Vote) int {
	count := 0
	for _, vote := range votes {
		count += vote.VoteType
	}
	return count
}
