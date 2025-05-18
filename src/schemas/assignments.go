package schemas

import "time"

type CreateAssignmentRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Type        string    `json:"type" binding:"required,oneof=ASSIGNMENT EXAM"`
	CourseID    string    `json:"course_id" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
}

type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type" binding:"omitempty,oneof=ASSIGNMENT EXAM"`
	DueDate     time.Time `json:"due_date"`
} 