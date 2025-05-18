package schemas

import "time"

type CreateAssignmentRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CourseID    string    `json:"course_id" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
}

type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
} 