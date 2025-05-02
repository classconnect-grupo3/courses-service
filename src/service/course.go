package service

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"errors"
	"time"
)

type CourseRepository interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c model.Course) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
}

type CourseService struct {
	courseRepository CourseRepository
}

func NewCourseService(courseRepository CourseRepository) *CourseService {
	return &CourseService{courseRepository: courseRepository}
}

func (s *CourseService) GetCourses() ([]*model.Course, error) {
	return s.courseRepository.GetCourses()
}

func (s *CourseService) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
	if c.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	//TODO: check teacher exists
	course := model.Course{
		Title:       c.Title,
		Description: c.Description,
		TeacherUUID: c.TeacherID,
		Capacity:    c.Capacity,
		CreatedAt:   time.Now(),
	}
	return s.courseRepository.CreateCourse(course)
}

func (s *CourseService) GetCourseById(id string) (*model.Course, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.courseRepository.GetCourseById(id)
}

func (s *CourseService) DeleteCourse(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.courseRepository.DeleteCourse(id)
}

func (s *CourseService) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	if teacherId == "" {
		return nil, errors.New("teacherId is required")
	}
	return s.courseRepository.GetCourseByTeacherId(teacherId)
}

func (s *CourseService) GetCourseByTitle(title string) ([]*model.Course, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	return s.courseRepository.GetCourseByTitle(title)
}
