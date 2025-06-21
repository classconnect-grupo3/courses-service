package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Module struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Order       int                `json:"order" bson:"order"`
	Resources   []ModuleResource   `json:"resources" bson:"resources"`
	CourseID    string             `json:"course_id" bson:"course_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ModuleResource struct {
	Id   uint64 `json:"id" bson:"id"` // this comes from the frontend so it's a number
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}
