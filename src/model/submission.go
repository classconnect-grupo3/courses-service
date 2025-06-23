package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmissionStatus string

const (
	SubmissionStatusDraft     SubmissionStatus = "draft"
	SubmissionStatusSubmitted SubmissionStatus = "submitted"
	SubmissionStatusLate      SubmissionStatus = "late"
)

type Answer struct {
	QuestionID string      `json:"question_id" bson:"question_id"`
	Content    interface{} `json:"content" bson:"content"` // Can be string, []string for multiple choice, or file URL
	Type       string      `json:"type" bson:"type"`       // text, multiple_choice, file
}

type Submission struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AssignmentID      string             `json:"assignment_id" bson:"assignment_id"`
	StudentUUID       string             `json:"student_uuid" bson:"student_uuid"`
	StudentName       string             `json:"student_name" bson:"student_name"`
	Status            SubmissionStatus   `json:"status" bson:"status"`
	Answers           []Answer           `json:"answers" bson:"answers"`
	Score             *float64           `json:"score,omitempty" bson:"score,omitempty"`
	Feedback          string             `json:"feedback,omitempty" bson:"feedback,omitempty"`
	AIScore           *float64           `json:"ai_score,omitempty" bson:"ai_score,omitempty"`
	AIFeedback        string             `json:"ai_feedback,omitempty" bson:"ai_feedback,omitempty"`
	NeedsManualReview *bool              `json:"needs_manual_review,omitempty" bson:"needs_manual_review,omitempty"`
	SubmittedAt       *time.Time         `json:"submitted_at,omitempty" bson:"submitted_at,omitempty"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}
