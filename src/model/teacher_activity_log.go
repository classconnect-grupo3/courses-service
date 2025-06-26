package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherActivityLog struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CourseID     string             `json:"course_id" bson:"course_id"`
	TeacherUUID  string             `json:"teacher_uuid" bson:"teacher_uuid"`
	ActivityType string             `json:"activity_type" bson:"activity_type"`
	Description  string             `json:"description" bson:"description"`
	Timestamp    time.Time          `json:"timestamp" bson:"timestamp"`
}
