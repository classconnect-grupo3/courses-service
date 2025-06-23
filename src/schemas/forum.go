package schemas

import (
	"courses-service/src/model"
	"time"
)

// Question schemas

type CreateQuestionRequest struct {
	CourseID    string              `json:"course_id" binding:"required"`
	AuthorID    string              `json:"author_id" binding:"required"`
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Tags        []model.QuestionTag `json:"tags"`
}

type UpdateQuestionRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Tags        []model.QuestionTag `json:"tags"`
}

type QuestionResponse struct {
	ID               string               `json:"id"`
	CourseID         string               `json:"course_id"`
	AuthorID         string               `json:"author_id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Tags             []model.QuestionTag  `json:"tags"`
	Votes            []model.Vote         `json:"votes"`
	VoteCount        int                  `json:"vote_count"`
	AnswerCount      int                  `json:"answer_count"`
	Status           model.QuestionStatus `json:"status"`
	AcceptedAnswerID *string              `json:"accepted_answer_id,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

type QuestionDetailResponse struct {
	ID               string               `json:"id"`
	CourseID         string               `json:"course_id"`
	AuthorID         string               `json:"author_id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Tags             []model.QuestionTag  `json:"tags"`
	Votes            []model.Vote         `json:"votes"`
	VoteCount        int                  `json:"vote_count"`
	Answers          []AnswerResponse     `json:"answers"`
	Status           model.QuestionStatus `json:"status"`
	AcceptedAnswerID *string              `json:"accepted_answer_id,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

// Answer schemas

type CreateAnswerRequest struct {
	AuthorID string `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdateAnswerRequest struct {
	Content string `json:"content" binding:"required"`
}

type AnswerResponse struct {
	ID         string       `json:"id"`
	AuthorID   string       `json:"author_id"`
	Content    string       `json:"content"`
	Votes      []model.Vote `json:"votes"`
	VoteCount  int          `json:"vote_count"`
	IsAccepted bool         `json:"is_accepted"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// Vote schemas

type VoteRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	VoteType int    `json:"vote_type" binding:"required"`
}

type VoteResponse struct {
	Message string `json:"message"`
}

// Search schemas

type SearchQuestionsRequest struct {
	Query  string               `form:"query"`
	Tags   []model.QuestionTag  `form:"tags"`
	Status model.QuestionStatus `form:"status"`
}

type SearchQuestionsResponse struct {
	Questions []QuestionResponse `json:"questions"`
	Total     int                `json:"total"`
}

// Generic response schemas

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
