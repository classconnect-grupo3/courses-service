package repository

import (
	"context"
	"courses-service/src/model"
)

type CourseRepositoryInterface interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c model.Course) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCoursesByStudentId(studentId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error)
	AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error)
	RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error)
}

type AssignmentRepositoryInterface interface {
	CreateAssignment(assignment model.Assignment) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetByID(ctx context.Context, id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error)
	DeleteAssignment(id string) error
}

type EnrollmentRepositoryInterface interface {
	CreateEnrollment(enrollment model.Enrollment, course *model.Course) error
	IsEnrolled(studentID, courseID string) (bool, error)
	DeleteEnrollment(studentID string, course *model.Course) error
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
}
