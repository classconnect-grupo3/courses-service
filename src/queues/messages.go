package queues

import (
	"time"
)

type QueueMessage interface {
	Encode() (map[string]any, error)
}

type AssignmentCreatedMessage struct {
	EventType         string    `json:"event_type"`
	CourseID          string    `json:"course_id"`
	AssignmentID      string    `json:"assignment_id"`
	AssignmentTitle   string    `json:"assignment_title"`
	AssignmentDueDate time.Time `json:"assignment_due_date"`
}

func NewAssignmentCreatedMessage(
	courseID string,
	assignmentID string,
	assignmentTitle string,
	assignmentDueDate time.Time,
) *AssignmentCreatedMessage {
	return &AssignmentCreatedMessage{
		EventType:         "assignment.created",
		CourseID:          courseID,
		AssignmentID:      assignmentID,
		AssignmentTitle:   assignmentTitle,
		AssignmentDueDate: assignmentDueDate,
	}
}

func (m *AssignmentCreatedMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type":          m.EventType,
		"course_id":           m.CourseID,
		"assignment_id":       m.AssignmentID,
		"assignment_title":    m.AssignmentTitle,
		"assignment_due_date": m.AssignmentDueDate.Format(time.RFC3339),
	}, nil
}
