package schemas

import "time"

// Period represents a time range for filtering statistics
type Period struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// AssignmentStatus represents the completion status of an assignment
type AssignmentStatus struct {
	AssignmentID   string    `json:"assignment_id"`
	AssignmentName string    `json:"assignment_name"`
	DueDate        time.Time `json:"due_date"`
	Submitted      bool      `json:"submitted"`
	Score          float64   `json:"score"`
	MaxScore       float64   `json:"max_score"`
}

// ParticipationMetrics represents student participation in a course
type ParticipationMetrics struct {
	ForumPosts      int     `json:"forum_posts"`
	ForumResponses  int     `json:"forum_responses"`
	AssignmentRatio float64 `json:"assignment_ratio"` // Assignments completed / total assignments
	AttendanceRatio float64 `json:"attendance_ratio"` // Optional if your system tracks attendance
}

// StudentPerformanceSummary represents a summary of a student's performance in a course
type StudentPerformanceSummary struct {
	StudentID         string  `json:"student_id"`
	StudentName       string  `json:"student_name"`
	AverageScore      float64 `json:"average_score"`
	CompletionRate    float64 `json:"completion_rate"`
	ParticipationRate float64 `json:"participation_rate"`
}

// StudentStats contains all statistics for a single student
type StudentStats struct {
	PerformanceSummary   StudentPerformanceSummary
	StudentScore         float64
	CompletedAssignments int
	ExamScore            float64
	ExamCompleted        int
	HomeworkScore        float64
	HomeworkCompleted    int
	ForumPosts           int
	ForumParticipated    bool
	ForumQuestions       int
	ForumAnswers         int
}

// CourseStatisticsRequest represents a request for course statistics
type CourseStatisticsRequest struct {
	CourseID string    `json:"course_id" binding:"required"`
	From     time.Time `json:"from" form:"from"`
	To       time.Time `json:"to" form:"to"`
}

// StudentStatisticsRequest represents a request for student statistics
type StudentStatisticsRequest struct {
	StudentID string    `json:"student_id" binding:"required"`
	CourseID  string    `json:"course_id" form:"course_id"`
	From      time.Time `json:"from" form:"from"`
	To        time.Time `json:"to" form:"to"`
}

// CourseStatisticsResponse represents the response for course statistics
type CourseStatisticsResponse struct {
	CourseID                string  `json:"course_id"`
	CourseName              string  `json:"course_name"`
	Period                  Period  `json:"period"`
	AverageScore            float64 `json:"average_score"`              // Promedio general de las notas (examenes + homeworks)
	AssignmentCompletion    float64 `json:"assignment_completion_rate"` // % global de entrega de assignments
	ExamCompletionRate      float64 `json:"exam_completion_rate"`       // % de entrega de exámenes
	HomeworkCompletionRate  float64 `json:"homework_completion_rate"`   // % de entrega de homeworks
	ForumParticipationRate  float64 `json:"forum_participation_rate"`
	ExamScoreAverage        float64 `json:"exam_average"`     // Promedio de notas en los exámenes
	HomeworkScoreAverage    float64 `json:"homework_average"` // Promedio de notas en las homeworks
	TotalStudents           int     `json:"total_students"`
	TotalAssignments        int     `json:"total_assignments"`
	TotalAmountOfExams      int     `json:"total_amount_of_exams"`
	TotalAmountOfHomeworks  int     `json:"total_amount_of_hw"`
	ForumUniqueParticipants int     `json:"forum_unique_participants"`
}

// ExportFormat represents the format for exporting statistics
type ExportFormat string

const (
	CSV ExportFormat = "csv"
)
