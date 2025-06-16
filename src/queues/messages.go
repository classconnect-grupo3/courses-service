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

type AddedAuxTeacherToCourseMessage struct {
	EventType string `json:"event_type"`
	CourseID  string `json:"course_id"`
	CourseName string `json:"course_name"`
	TeacherID string `json:"teacher_id"`
}

func NewAddedAuxTeacherToCourseMessage(courseID string, courseName string, teacherID string) *AddedAuxTeacherToCourseMessage {
	return &AddedAuxTeacherToCourseMessage{
		EventType: "aux_teacher.added",
		CourseID:  courseID,
		CourseName: courseName,
		TeacherID: teacherID,
	}
}	

func (m *AddedAuxTeacherToCourseMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type": m.EventType,
		"course_id":  m.CourseID,
		"course_name": m.CourseName,
		"teacher_id": m.TeacherID,
	}, nil
}

type RemoveAuxTeacherFromCourseMessage struct {
	EventType  string `json:"event_type"`
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
	TeacherID  string `json:"teacher_id"`
}

func NewRemoveAuxTeacherFromCourseMessage(courseID string, courseName string, teacherID string) *RemoveAuxTeacherFromCourseMessage {
	return &RemoveAuxTeacherFromCourseMessage{
		EventType:  "aux_teacher.removed",
		CourseID:   courseID,
		CourseName: courseName,
		TeacherID:  teacherID,
	}
}

func (m *RemoveAuxTeacherFromCourseMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type":  m.EventType,
		"course_id":   m.CourseID,
		"course_name": m.CourseName,
		"teacher_id":  m.TeacherID,
	}, nil
}

type FeedbackCreatedMessage struct {
	EventType         string    `json:"event_type"`
	CourseID          string    `json:"course_id"`
	FeedbackID        string    `json:"feedback_id"`
	FeedbackText      string    `json:"feedback_text"`
	FeedbackRating    int       `json:"feedback_rating"`
	FeedbackCreatedAt time.Time `json:"feedback_created_at"`
}

func NewFeedbackCreatedMessage(courseID string, feedbackID string, feedbackText string, feedbackRating int, feedbackCreatedAt time.Time) *FeedbackCreatedMessage {
	return &FeedbackCreatedMessage{
		EventType:         "feedback.created",
		CourseID:          courseID,
		FeedbackID:        feedbackID,
		FeedbackText:      feedbackText,
		FeedbackRating:    feedbackRating,
		FeedbackCreatedAt: feedbackCreatedAt,
	}
}

func (m *FeedbackCreatedMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type":          m.EventType,
		"course_id":           m.CourseID,
		"feedback_id":         m.FeedbackID,
		"feedback_text":       m.FeedbackText,
		"feedback_rating":     m.FeedbackRating,
		"feedback_created_at": m.FeedbackCreatedAt.Format(time.RFC3339),
	}, nil
}

type EnrolledStudentToCourseMessage struct {
	EventType string `json:"event_type"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

func NewEnrolledStudentToCourseMessage(courseID string, studentID string) *EnrolledStudentToCourseMessage {
	return &EnrolledStudentToCourseMessage{
		EventType: "student.enrolled",
		CourseID:  courseID,
		StudentID: studentID,
	}
}

func (m *EnrolledStudentToCourseMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type": m.EventType,	
		"course_id":  m.CourseID,
		"student_id": m.StudentID,
	}, nil
}

type UnenrolledStudentFromCourseMessage struct {
	EventType string `json:"event_type"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

type AICorrectionMessage struct {
	EventType string `json:"event_type"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
	SubmissionID string `json:"submission_id"`
	SubmissionText string `json:"submission_text"`
	SubmissionCreatedAt time.Time `json:"submission_created_at"`
	SubmissionFeedback string `json:"submission_feedback"`
}

func NewAICorrectionMessage(courseID string, studentID string, submissionID string, submissionText string, submissionCreatedAt time.Time) *AICorrectionMessage {
	return &AICorrectionMessage{
		EventType: "ai.correction",
		CourseID:  courseID,
		StudentID: studentID,
		SubmissionID: submissionID,
		SubmissionText: submissionText,
		SubmissionCreatedAt: submissionCreatedAt,
	}
}

func (m *AICorrectionMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type": m.EventType,
		"course_id":  m.CourseID,
		"student_id": m.StudentID,
		"submission_id": m.SubmissionID,
		"submission_text": m.SubmissionText,
		"submission_created_at": m.SubmissionCreatedAt.Format(time.RFC3339),
		"submission_feedback": m.SubmissionFeedback,
	}, nil
}

type ForumActivityMessage struct {
	EventType string `json:"event_type"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
	PostID string `json:"post_id"`
	PostText string `json:"post_text"`
	PostCreatedAt time.Time `json:"post_created_at"`
}

func NewForumActivityMessage(courseID string, studentID string, postID string, postText string, postCreatedAt time.Time) *ForumActivityMessage {
	return &ForumActivityMessage{
		EventType: "forum.activity",
		CourseID:  courseID,
		StudentID: studentID,
		PostID: postID,
		PostText: postText,
		PostCreatedAt: postCreatedAt,
	}
}

func (m *ForumActivityMessage) Encode() (map[string]any, error) {	
	return map[string]any{
		"event_type": m.EventType,
		"course_id":  m.CourseID,
		"student_id": m.StudentID,
		"post_id": m.PostID,
		"post_text": m.PostText,
		"post_created_at": m.PostCreatedAt.Format(time.RFC3339),
	}, nil
}