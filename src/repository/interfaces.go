package repository

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
)

type CourseRepositoryInterface interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c model.Course) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCoursesByStudentId(studentId string) ([]*model.Course, error)
	GetCoursesByAuxTeacherId(auxTeacherId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error)
	AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error)
	RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error)
	UpdateStudentsAmount(courseID string, newStudentsAmount int) error
	CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error)
	GetCourseFeedback(courseID string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error)

	// Backoffice statistics methods
	CountCourses() (int64, error)
	CountActiveCourses() (int64, error)
	CountFinishedCourses() (int64, error)
	CountCoursesCreatedThisMonth() (int64, error)
	CountUniqueTeachers() (int64, error)
	CountUniqueAuxTeachers() (int64, error)
	GetTopTeachersByCourseCount(limit int) ([]schemas.CourseDistributionByTeacher, error)
	GetRecentCourses(limit int) ([]schemas.CourseBasicInfo, error)
}

type AssignmentRepositoryInterface interface {
	CreateAssignment(assignment model.Assignment) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetByID(ctx context.Context, id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error)
	DeleteAssignment(id string) error

	// Backoffice statistics methods
	CountAssignments() (int64, error)
	CountAssignmentsByType(assignmentType string) (int64, error)
	CountAssignmentsByStatus(status string) (int64, error)
	CountAssignmentsCreatedThisMonth() (int64, error)
	GetAssignmentDistribution() ([]schemas.AssignmentDistribution, error)
	GetRecentAssignments(limit int) ([]schemas.AssignmentBasicInfo, error)
}

type EnrollmentRepositoryInterface interface {
	CreateEnrollment(enrollment model.Enrollment, course *model.Course) error
	IsEnrolled(studentID, courseID string) (bool, error)
	DeleteEnrollment(studentID string, course *model.Course) error
	GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error)
	SetFavouriteCourse(studentID, courseID string) error
	UnsetFavouriteCourse(studentID, courseID string) error
	GetEnrollmentsByStudentId(studentID string) ([]*model.Enrollment, error)
	GetEnrollmentByStudentIdAndCourseId(studentID, courseID string) (*model.Enrollment, error)
	CreateStudentFeedback(feedbackRequest model.StudentFeedback, enrollmentID string) error
	GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error)
	ApproveStudent(studentID, courseID string) error
	DisapproveStudent(studentID, courseID, reason string) error
	ReactivateDroppedEnrollment(studentID, courseID string) error

	// Backoffice statistics methods
	CountEnrollments() (int64, error)
	CountEnrollmentsByStatus(status model.EnrollmentStatus) (int64, error)
	CountEnrollmentsThisMonth() (int64, error)
	CountUniqueStudents() (int64, error)
}

type ModuleRepositoryInterface interface {
	GetNextModuleOrder(courseID string) (int, error)
	CreateModule(courseID string, module model.Module) (*model.Module, error)
	GetModuleById(id string) (*model.Module, error)
	UpdateModule(id string, module model.Module) (*model.Module, error)
	DeleteModule(id string) error
	GetModulesByCourseId(courseId string) ([]model.Module, error)
	GetModuleByName(courseID string, moduleName string) (*model.Module, error)
	GetModuleByOrder(courseID string, order int) (*model.Module, error)
}

type SubmissionRepositoryInterface interface {
	Create(ctx context.Context, submission *model.Submission) error
	Update(ctx context.Context, submission *model.Submission) error
	GetByID(ctx context.Context, id string) (*model.Submission, error)
	GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error)
	GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error)
	GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error)
	DeleteByStudentAndCourse(ctx context.Context, studentUUID, courseID string) error

	// Backoffice statistics methods
	CountSubmissions(ctx context.Context) (int64, error)
	CountSubmissionsByStatus(ctx context.Context, status model.SubmissionStatus) (int64, error)
	CountSubmissionsThisMonth(ctx context.Context) (int64, error)
}

type ForumRepositoryInterface interface {
	// Question operations
	CreateQuestion(question model.ForumQuestion) (*model.ForumQuestion, error)
	GetQuestionById(id string) (*model.ForumQuestion, error)
	GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error)
	UpdateQuestion(id string, question model.ForumQuestion) (*model.ForumQuestion, error)
	DeleteQuestion(id string) error

	// Answer operations
	AddAnswer(questionID string, answer model.ForumAnswer) (*model.ForumAnswer, error)
	UpdateAnswer(questionID string, answerID string, content string) (*model.ForumAnswer, error)
	DeleteAnswer(questionID string, answerID string) error
	AcceptAnswer(questionID string, answerID string) error

	// Vote operations
	AddVoteToQuestion(questionID string, userID string, voteType int) error
	AddVoteToAnswer(questionID string, answerID string, userID string, voteType int) error
	RemoveVoteFromQuestion(questionID string, userID string) error
	RemoveVoteFromAnswer(questionID string, answerID string, userID string) error

	// Search and filter operations
	SearchQuestions(courseID string, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error)

	// Backoffice statistics methods
	CountQuestions() (int64, error)
	CountQuestionsByStatus(status model.QuestionStatus) (int64, error)
	CountAnswers() (int64, error)
}
