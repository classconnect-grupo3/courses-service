package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackType string

const (
	FeedbackTypePositive FeedbackType = "POSITIVO"
	FeedbackTypeNegative FeedbackType = "NEGATIVO"
	FeedbackTypeNeutral  FeedbackType = "NEUTRO"
)

var FeedbackTypes = []FeedbackType{FeedbackTypePositive, FeedbackTypeNegative, FeedbackTypeNeutral}

type StudentFeedback struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentUUID  string             `json:"student_uuid" bson:"student_uuid"`
	TeacherUUID  string             `json:"teacher_uuid" bson:"teacher_uuid"`
	CourseID     string             `json:"course_id" bson:"course_id"`
	FeedbackType FeedbackType       `json:"feedback_type" bson:"feedback_type"`
	Score        int                `json:"score" bson:"score"`
	Feedback     string             `json:"feedback" bson:"feedback"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type CourseFeedback struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentUUID  string             `json:"student_uuid" bson:"student_uuid"`
	FeedbackType FeedbackType       `json:"feedback_type" bson:"feedback_type"`
	Score        int                `json:"score" bson:"score"`
	Feedback     string             `json:"feedback" bson:"feedback"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at,omitempty"`
}
