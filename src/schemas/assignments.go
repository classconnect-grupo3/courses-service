package schemas

import "time"

type CreateAssignmentRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Instructions    string    `json:"instructions" binding:"required"`
	Type            string    `json:"type" binding:"required,oneof=ASSIGNMENT EXAM"`
	CourseID        string    `json:"course_id" binding:"required"`
	DueDate         time.Time `json:"due_date" binding:"required"`
	GracePeriod     int       `json:"grace_period" binding:"min=0"`                    // Optional, defaults to 0
	Status          string    `json:"status" binding:"required,oneof=draft published"` // Required status
	SubmissionRules []string  `json:"submission_rules"`                                // Optional rules
}

type UpdateAssignmentRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Instructions    string    `json:"instructions"`
	Type            string    `json:"type" binding:"omitempty,oneof=ASSIGNMENT EXAM"`
	DueDate         time.Time `json:"due_date"`
	GracePeriod     int       `json:"grace_period" binding:"omitempty,min=0"`
	Status          string    `json:"status" binding:"omitempty,oneof=draft published"`
	SubmissionRules []string  `json:"submission_rules"`
} 