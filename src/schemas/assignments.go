package schemas

import (
	"courses-service/src/model"
	"time"
)

type CreateAssignmentRequest struct {
	Title           string          `json:"title" binding:"required"`
	Description     string          `json:"description" binding:"required"`
	Instructions    string          `json:"instructions" binding:"required"`
	Type            string          `json:"type" binding:"required"`
	CourseID        string          `json:"course_id" binding:"required"`
	DueDate         time.Time       `json:"due_date" binding:"required"`
	GracePeriod     int            `json:"grace_period" binding:"required"`
	Status          string          `json:"status" binding:"required"`
	Questions       []model.Question `json:"questions" binding:"required"`
	TotalPoints     float64         `json:"total_points" binding:"required"`
	PassingScore    float64         `json:"passing_score" binding:"required"`
}

type UpdateAssignmentRequest struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Instructions    string          `json:"instructions"`
	Type            string          `json:"type"`
	DueDate         time.Time       `json:"due_date"`
	GracePeriod     int            `json:"grace_period"`
	Status          string          `json:"status"`
	Questions       []model.Question `json:"questions"`
	TotalPoints     float64         `json:"total_points"`
	PassingScore    float64         `json:"passing_score"`
} 