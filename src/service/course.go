package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
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

func (s *CourseService) CreateCourse(course schemas.CreateCourseRequest) (*model.Course, error) {
	return s.courseRepository.CreateCourse(course.Title, course.Description, course.TeacherID, course.Capacity)
}

func (s *CourseService) GetCourseById(id string) (*model.Course, error) {
	return s.courseRepository.GetCourseById(id)
}

func (s *CourseService) DeleteCourse(id string) error {
	return s.courseRepository.DeleteCourse(id)
}

func (s *CourseService) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return s.courseRepository.GetCourseByTeacherId(teacherId)
}

func (s *CourseService) GetCourseByTitle(title string) ([]*model.Course, error) {
	return s.courseRepository.GetCourseByTitle(title)
}
