package model

import "time"

type Module struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Order       int       `json:"order" bson:"order"`
	Content     string    `json:"content" bson:"content"` // TODO change this with media in the future
	CourseID    string    `json:"course_id" bson:"course_id"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
