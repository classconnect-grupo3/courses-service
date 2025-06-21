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
	EventType  string `json:"event_type"`
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
	TeacherID  string `json:"teacher_id"`
}

func NewAddedAuxTeacherToCourseMessage(courseID string, courseName string, teacherID string) *AddedAuxTeacherToCourseMessage {
	return &AddedAuxTeacherToCourseMessage{
		EventType:  "aux_teacher.added",
		CourseID:   courseID,
		CourseName: courseName,
		TeacherID:  teacherID,
	}
}

func (m *AddedAuxTeacherToCourseMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type":  m.EventType,
		"course_id":   m.CourseID,
		"course_name": m.CourseName,
		"teacher_id":  m.TeacherID,
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

type ForumActivityMessage struct {
	EventType     string    `json:"event_type"`
	CourseID      string    `json:"course_id"`
	StudentID     string    `json:"student_id"`
	PostID        string    `json:"post_id"`
	PostText      string    `json:"post_text"`
	PostCreatedAt time.Time `json:"post_created_at"`
}

func NewForumActivityMessage(courseID string, studentID string, postID string, postText string, postCreatedAt time.Time) *ForumActivityMessage {
	return &ForumActivityMessage{
		EventType:     "forum.activity",
		CourseID:      courseID,
		StudentID:     studentID,
		PostID:        postID,
		PostText:      postText,
		PostCreatedAt: postCreatedAt,
	}
}

func (m *ForumActivityMessage) Encode() (map[string]any, error) {
	return map[string]any{
		"event_type":      m.EventType,
		"course_id":       m.CourseID,
		"student_id":      m.StudentID,
		"post_id":         m.PostID,
		"post_text":       m.PostText,
		"post_created_at": m.PostCreatedAt.Format(time.RFC3339),
	}, nil
}

type SubmissionCorrectedMessage struct {
	EventType         string    `json:"event_type"`
	CourseID          string    `json:"course_id"`
	AssignmentID      string    `json:"assignment_id"`
	SubmissionID      string    `json:"submission_id"`
	StudentID         string    `json:"student_id"`
	Score             *float64  `json:"score,omitempty"`
	Feedback          string    `json:"feedback"`
	CorrectionType    string    `json:"correction_type"` // "automatic", "needs_manual_review"
	NeedsManualReview bool      `json:"needs_manual_review"`
	CorrectedAt       time.Time `json:"corrected_at"`
}

func NewSubmissionCorrectedMessage(
	assignmentID string,
	submissionID string,
	studentID string,
	score *float64,
	feedback string,
	correctionType string,
	needsManualReview bool,
) *SubmissionCorrectedMessage {
	return &SubmissionCorrectedMessage{
		EventType:         "submission.corrected",
		AssignmentID:      assignmentID,
		SubmissionID:      submissionID,
		StudentID:         studentID,
		Score:             score,
		Feedback:          feedback,
		CorrectionType:    correctionType,
		NeedsManualReview: needsManualReview,
		CorrectedAt:       time.Now(),
	}
}

func (m *SubmissionCorrectedMessage) Encode() (map[string]any, error) {
	encoded := map[string]any{
		"event_type":          m.EventType,
		"assignment_id":       m.AssignmentID,
		"submission_id":       m.SubmissionID,
		"student_id":          m.StudentID,
		"feedback":            m.Feedback,
		"correction_type":     m.CorrectionType,
		"needs_manual_review": m.NeedsManualReview,
		"corrected_at":        m.CorrectedAt.Format(time.RFC3339),
	}

	if m.Score != nil {
		encoded["score"] = *m.Score
	}

	return encoded, nil
}
