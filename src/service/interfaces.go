package service

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"time"
)

// CourseServiceInterface define los métodos que debe implementar un servicio de cursos
type CourseServiceInterface interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string, teacherId string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCoursesByStudentId(studentId string) ([]*model.Course, error)
	GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error)
	AddAuxTeacherToCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error)
	RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error)
	GetFavouriteCourses(studentId string) ([]*model.Course, error)
	CreateCourseFeedback(courseId string, feedbackRequest schemas.CreateCourseFeedbackRequest) (*model.CourseFeedback, error)
	GetCourseFeedback(courseId string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error)
	GetCourseMembers(courseId string) (*schemas.CourseMembersResponse, error)
}

type ModuleServiceInterface interface {
	CreateModule(module schemas.CreateModuleRequest) (*model.Module, error)
	GetModuleById(id string) (*model.Module, error)
	GetModulesByCourseId(courseId string) ([]model.Module, error)
	UpdateModule(id string, module model.Module) (*model.Module, error)
	DeleteModule(id string) error
}

// EnrollmentServiceInterface define los métodos que debe implementar un servicio de enrollment
type EnrollmentServiceInterface interface {
	GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error)
	EnrollStudent(studentID, courseID string) error
	UnenrollStudent(studentID, courseID string) error
	SetFavouriteCourse(studentID, courseID string) error
	UnsetFavouriteCourse(studentID, courseID string) error
	CreateStudentFeedback(feedbackRequest schemas.CreateStudentFeedbackRequest) error
	GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error)
	ApproveStudent(studentID, courseID string) error
	DisapproveStudent(studentID, courseID, reason string) error
}

type AssignmentServiceInterface interface {
	CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetAssignmentById(id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error)
	DeleteAssignment(id string) error
}

type SubmissionServiceInterface interface {
	CreateSubmission(ctx context.Context, submission *model.Submission) error
	UpdateSubmission(ctx context.Context, submission *model.Submission) error
	SubmitSubmission(ctx context.Context, submissionID string) error
	GetSubmission(ctx context.Context, id string) (*model.Submission, error)
	GetSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error)
	GetSubmissionsByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error)
	GetOrCreateSubmission(ctx context.Context, assignmentID, studentUUID, studentName string) (*model.Submission, error)
	GradeSubmission(ctx context.Context, submissionID string, score *float64, feedback string) (*model.Submission, error)
	ValidateTeacherPermissions(ctx context.Context, assignmentID, teacherUUID string) error
	GenerateFeedbackSummary(ctx context.Context, submissionID string) (*schemas.AiSummaryResponse, error)
	AutoCorrectSubmission(ctx context.Context, submissionID string) error
}

type ForumServiceInterface interface {
	// Question operations
	CreateQuestion(courseID, authorID, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error)
	GetQuestionById(id string) (*model.ForumQuestion, error)
	GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error)
	UpdateQuestion(id, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error)
	DeleteQuestion(id, authorID string) error

	// Answer operations
	AddAnswer(questionID, authorID, content string) (*model.ForumAnswer, error)
	UpdateAnswer(questionID, answerID, authorID, content string) (*model.ForumAnswer, error)
	DeleteAnswer(questionID, answerID, authorID string) error
	AcceptAnswer(questionID, answerID, authorID string) error

	// Vote operations
	VoteQuestion(questionID, userID string, voteType int) error
	VoteAnswer(questionID, answerID, userID string, voteType int) error
	RemoveVoteFromQuestion(questionID, userID string) error
	RemoveVoteFromAnswer(questionID, answerID, userID string) error

	// Search and filter operations
	SearchQuestions(courseID, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error)

	// Forum participants operations
	GetForumParticipants(courseID string) ([]string, error)
}

// StatisticsServiceInterface define los métodos que debe implementar un servicio de estadísticas
type StatisticsServiceInterface interface {

	// ExportCourseStatsCSV genera un CSV con las estadísticas del curso
	ExportCourseStatsCSV(ctx context.Context, courseID string, from, to time.Time) ([]byte, string, error)

	// ExportStudentStatsCSV genera un CSV con las estadísticas del estudiante
	ExportStudentStatsCSV(ctx context.Context, studentID string, courseID string, from, to time.Time) ([]byte, string, error)

	ExportTeacherCoursesStatsCSV(ctx context.Context, teacherID string, from, to time.Time) ([]byte, string, error)

	// Backoffice statistics methods
	GetBackofficeStatistics(ctx context.Context) (*schemas.BackofficeStatisticsResponse, error)
	GetBackofficeCoursesStats(ctx context.Context) (*schemas.BackofficeCoursesStatsResponse, error)
	GetBackofficeAssignmentsStats(ctx context.Context) (*schemas.BackofficeAssignmentsStatsResponse, error)
}

type TeacherActivityServiceInterface interface {
	LogActivityIfAuxTeacher(courseID, teacherUUID, activityType, description string)
	GetCourseActivityLogs(courseID string) ([]*model.TeacherActivityLog, error)
}
