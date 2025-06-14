package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionTag string

const (
	QuestionTagGeneral       QuestionTag = "general"
	QuestionTagTeoria        QuestionTag = "teoria"
	QuestionTagPractica      QuestionTag = "practica"
	QuestionTagNecesitoAyuda QuestionTag = "necesito-ayuda"
	QuestionTagInformacion   QuestionTag = "informacion"
	QuestionTagEjercitacion  QuestionTag = "ejercitacion"
	QuestionTagOtro          QuestionTag = "otro"
)

var QuestionTagValues = []QuestionTag{
	QuestionTagGeneral,
	QuestionTagTeoria,
	QuestionTagPractica,
	QuestionTagNecesitoAyuda,
	QuestionTagInformacion,
	QuestionTagEjercitacion,
	QuestionTagOtro,
}

type QuestionStatus string

const (
	QuestionStatusOpen     QuestionStatus = "open"
	QuestionStatusResolved QuestionStatus = "resolved"
	QuestionStatusClosed   QuestionStatus = "closed"
)

var QuestionStatusValues = []QuestionStatus{
	QuestionStatusOpen,
	QuestionStatusResolved,
	QuestionStatusClosed,
}

const (
	VoteTypeUp   = 1
	VoteTypeDown = -1
)

// ForumQuestion represents a forum question
type ForumQuestion struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	CourseID         string             `json:"course_id" bson:"course_id"`
	AuthorID         string             `json:"author_id" bson:"author_id"`
	Title            string             `json:"title" bson:"title"`
	Description      string             `json:"description" bson:"description"`
	Tags             []QuestionTag      `json:"tags" bson:"tags"`
	Votes            []Vote             `json:"votes" bson:"votes"`
	Answers          []ForumAnswer      `json:"answers" bson:"answers"`
	Status           QuestionStatus     `json:"status" bson:"status"`
	AcceptedAnswerID *string            `json:"accepted_answer_id,omitempty" bson:"accepted_answer_id,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// ForumAnswer represents a forum answer to a question
type ForumAnswer struct {
	ID         string    `json:"id" bson:"_id"`
	AuthorID   string    `json:"author_id" bson:"author_id"`
	Content    string    `json:"content" bson:"content"`
	Votes      []Vote    `json:"votes" bson:"votes"`
	IsAccepted bool      `json:"is_accepted" bson:"is_accepted"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

// Vote represents a vote on a question or answer
type Vote struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	VoteType  int       `json:"vote_type" bson:"vote_type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
