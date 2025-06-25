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

// BackofficeStatisticsResponse represents general system statistics for backoffice
type BackofficeStatisticsResponse struct {
	// General counts
	TotalCourses        int `json:"total_courses"`
	TotalAssignments    int `json:"total_assignments"`
	TotalSubmissions    int `json:"total_submissions"`
	TotalEnrollments    int `json:"total_enrollments"`
	TotalForumQuestions int `json:"total_forum_questions"`
	TotalForumAnswers   int `json:"total_forum_answers"`

	// Course statistics
	ActiveCourses   int `json:"active_courses"`
	FinishedCourses int `json:"finished_courses"`

	// Assignment statistics by type
	TotalExams     int `json:"total_exams"`
	TotalHomeworks int `json:"total_homeworks"`
	TotalQuizzes   int `json:"total_quizzes"`

	// Submission statistics by status
	DraftSubmissions     int `json:"draft_submissions"`
	SubmittedSubmissions int `json:"submitted_submissions"`
	LateSubmissions      int `json:"late_submissions"`

	// Enrollment statistics by status
	ActiveEnrollments    int `json:"active_enrollments"`
	DroppedEnrollments   int `json:"dropped_enrollments"`
	CompletedEnrollments int `json:"completed_enrollments"`

	// Forum statistics by status
	OpenForumQuestions     int `json:"open_forum_questions"`
	ResolvedForumQuestions int `json:"resolved_forum_questions"`
	ClosedForumQuestions   int `json:"closed_forum_questions"`

	// Teacher statistics
	TotalUniqueTeachers    int `json:"total_unique_teachers"`
	TotalUniqueAuxTeachers int `json:"total_unique_aux_teachers"`

	// Student statistics
	TotalUniqueStudents int `json:"total_unique_students"`

	// Average statistics
	AverageStudentsPerCourse        float64 `json:"average_students_per_course"`
	AverageAssignmentsPerCourse     float64 `json:"average_assignments_per_course"`
	AverageSubmissionsPerAssignment float64 `json:"average_submissions_per_assignment"`

	// Date-based statistics
	CoursesCreatedThisMonth     int `json:"courses_created_this_month"`
	AssignmentsCreatedThisMonth int `json:"assignments_created_this_month"`
	SubmissionsThisMonth        int `json:"submissions_this_month"`
	EnrollmentsThisMonth        int `json:"enrollments_this_month"`
}

// BackofficeCoursesStatsResponse represents detailed course statistics for backoffice
type BackofficeCoursesStatsResponse struct {
	TotalCourses    int                `json:"total_courses"`
	CoursesByStatus map[string]int     `json:"courses_by_status"`
	RecentCourses   []CourseBasicInfo  `json:"recent_courses"`
}

// CourseBasicInfo represents basic course information
type CourseBasicInfo struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	TeacherName    string    `json:"teacher_name"`
	StudentsAmount int       `json:"students_amount"`
	Capacity       int       `json:"capacity"`
	CreatedAt      time.Time `json:"created_at"`
}

// AssignmentDistribution represents assignment distribution by type and status
type AssignmentDistribution struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// BackofficeAssignmentsStatsResponse represents detailed assignment statistics for backoffice
type BackofficeAssignmentsStatsResponse struct {
	TotalAssignments       int                      `json:"total_assignments"`
	AssignmentsByType      map[string]int           `json:"assignments_by_type"`
	AssignmentsByStatus    map[string]int           `json:"assignments_by_status"`
	AssignmentDistribution []AssignmentDistribution `json:"assignment_distribution"`
	RecentAssignments      []AssignmentBasicInfo    `json:"recent_assignments"`
}

// AssignmentBasicInfo represents basic assignment information
type AssignmentBasicInfo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CourseID  string    `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	DueDate   time.Time `json:"due_date"`
}
