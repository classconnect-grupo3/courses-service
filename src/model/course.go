package model

import "time"

type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherUUID string    `json:"teacher_uuid"`
	CreatedAt   time.Time `json:"created_at"`
	// TODO: add reamining fields
}
