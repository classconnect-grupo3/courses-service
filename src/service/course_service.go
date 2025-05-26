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
	GetCoursesByStudentId(studentId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error)
}

type CourseService interface {
	GetCourses() ([]*model.Course, error)
	CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error)
	GetCourseById(id string) (*model.Course, error)
	DeleteCourse(id string) error
	GetCourseByTeacherId(teacherId string) ([]*model.Course, error)
	GetCourseByTitle(title string) ([]*model.Course, error)
	UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error)
}

type CourseServiceImpl struct {
	courseRepository CourseRepository
}

func NewCourseService(courseRepository CourseRepository) CourseService {
	return &CourseServiceImpl{courseRepository: courseRepository}
}

func (s *CourseServiceImpl) GetCourses() ([]*model.Course, error) {
	return s.courseRepository.GetCourses()
}

func (s *CourseServiceImpl) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
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
		UpdatedAt:   time.Now(),
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
	}
	return s.courseRepository.CreateCourse(course)
}

func (s *CourseServiceImpl) GetCourseById(id string) (*model.Course, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.courseRepository.GetCourseById(id)
}

func (s *CourseServiceImpl) DeleteCourse(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.courseRepository.DeleteCourse(id)
}

func (s *CourseServiceImpl) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	if teacherId == "" {
		return nil, errors.New("teacherId is required")
	}
	return s.courseRepository.GetCourseByTeacherId(teacherId)
}

func (s *CourseService) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	if studentId == "" {
		return nil, errors.New("studentId is required")
	}
	return s.courseRepository.GetCoursesByStudentId(studentId)
}

func (s *CourseService) GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error) {
	if userId == "" {
		return nil, errors.New("userId is required")
	}
	result := schemas.GetCoursesByUserIdResponse{}

	studentCourses, err := s.courseRepository.GetCoursesByStudentId(userId)
	if err != nil {
		return nil, err
	}

	teacherCourses, err := s.courseRepository.GetCourseByTeacherId(userId)
	if err != nil {
		return nil, err
	}

	result.Student = studentCourses
	result.Teacher = teacherCourses

	return &result, nil
}

func (s *CourseService) GetCourseByTitle(title string) ([]*model.Course, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	return s.courseRepository.GetCourseByTitle(title)
}

func (s *CourseServiceImpl) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	course := model.Course{
		Title:       updateCourseRequest.Title,
		Description: updateCourseRequest.Description,
		TeacherUUID: updateCourseRequest.TeacherID,
		Capacity:    updateCourseRequest.Capacity,
		UpdatedAt:   time.Now(),
	}
	return s.courseRepository.UpdateCourse(id, course)
}
