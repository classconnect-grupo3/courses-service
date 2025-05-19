package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionType string

const (
	QuestionTypeText           QuestionType = "text"
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeFile          QuestionType = "file"
)

type Question struct {
	ID              string      `json:"id" bson:"id"`
	Text            string      `json:"text" bson:"text"`
	Type            QuestionType `json:"type" bson:"type"`
	Options         []string    `json:"options,omitempty" bson:"options,omitempty"` // For multiple choice
	CorrectAnswers  []string    `json:"correct_answers,omitempty" bson:"correct_answers,omitempty"`
	Points          float64     `json:"points" bson:"points"`
	Order           int         `json:"order" bson:"order"`
}

type Assignment struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title           string            `json:"title" bson:"title"`
	Description     string            `json:"description" bson:"description"`
	Instructions    string            `json:"instructions" bson:"instructions"`
	Type            string            `json:"type" bson:"type"`                      // exam, homework, quiz
	CourseID        string            `json:"course_id" bson:"course_id"`
	DueDate         time.Time         `json:"due_date" bson:"due_date"`
	GracePeriod     int              `json:"grace_period" bson:"grace_period"`      // Minutes of tolerance after due_date
	Status          string            `json:"status" bson:"status"`                  // draft, published
	Questions       []Question        `json:"questions" bson:"questions"`
	TotalPoints     float64          `json:"total_points" bson:"total_points"`
	PassingScore    float64          `json:"passing_score" bson:"passing_score"`    // Minimum score to pass
	SubmissionRules []string         `json:"submission_rules" bson:"submission_rules"` // Array of rules for submission
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" bson:"updated_at"`
} 