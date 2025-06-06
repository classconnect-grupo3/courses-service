package service

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
)

// CourseServiceInterface define los métodos que debe implementar un servicio de cursos
type CourseServiceInterface interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCoursesByStudentId(studentId string) ([]*model.Course, error)
	GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error)
	AddAuxTeacherToCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error)
	RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error)
	GetFavouriteCourses(studentId string) ([]*model.Course, error)
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
}

type AssignmentServiceInterface interface {
	CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetAssignmentById(id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error)
	DeleteAssignment(id string) error
}
