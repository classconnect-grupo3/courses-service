package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enrollment struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID     string             `json:"student_id" bson:"student_id"`
	CourseID      string             `json:"course_id" bson:"course_id"`
	EnrolledAt    time.Time          `json:"enrolled_at" bson:"enrolled_at"`
	CompletedDate time.Time          `json:"completed_date" bson:"completed_date"`
	Status        EnrollmentStatus   `json:"status" bson:"status"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

type EnrollmentStatus string

const (
	EnrollmentStatusActive    EnrollmentStatus = "active"
	EnrollmentStatusDropped   EnrollmentStatus = "dropped"
	EnrollmentStatusCompleted EnrollmentStatus = "completed"
)
