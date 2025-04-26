package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
)

type CourseService struct {
	courseRepository *repository.CourseRepository
}

func NewCourseService(courseRepository *repository.CourseRepository) *CourseService {
	return &CourseService{courseRepository: courseRepository}
}

func (s *CourseService) GetCourses() ([]*model.Course, error) {
	return s.courseRepository.GetCourses()
}

func (s *CourseService) CreateCourse(course model.Course) (*model.Course, error) {
	return s.courseRepository.CreateCourse(course)
}

func (s *CourseService) GetCourseById(id string) (*model.Course, error) {
	return s.courseRepository.GetCourseById(id)
}

func (s *CourseService) DeleteCourse(id string) error {
	return s.courseRepository.DeleteCourse(id)
}
