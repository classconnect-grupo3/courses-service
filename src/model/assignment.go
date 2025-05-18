package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title" bson:"title"`
	Description string            `json:"description" bson:"description"`
	Type        string            `json:"type" bson:"type"`
	CourseID    string            `json:"course_id" bson:"course_id"`
	DueDate     time.Time         `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
} 