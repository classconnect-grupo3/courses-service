package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	TeacherUUID    string             `json:"teacher_uuid" bson:"teacher_uuid"`
	TeacherName    string             `json:"teacher_name" bson:"teacher_name"`
	Capacity       int                `json:"capacity" bson:"capacity"`
	StudentsAmount int                `json:"students_amount" bson:"students_amount"`
	StartDate      time.Time          `json:"start_date" bson:"start_date"`
	EndDate        time.Time          `json:"end_date" bson:"end_date"`
	Modules        []Module           `json:"modules" bson:"modules"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
