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