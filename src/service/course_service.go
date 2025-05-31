package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"errors"
	"slices"
	"time"
)

type CourseService struct {
	courseRepository     repository.CourseRepositoryInterface
	enrollmentRepository repository.EnrollmentRepositoryInterface
}

func NewCourseService(courseRepository repository.CourseRepositoryInterface, enrollmentRepository repository.EnrollmentRepositoryInterface) *CourseService {
	return &CourseService{courseRepository: courseRepository, enrollmentRepository: enrollmentRepository}
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
		UpdatedAt:   time.Now(),
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
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

func (s *CourseService) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
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

func (s *CourseService) AddAuxTeacherToCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	course, err := s.courseRepository.GetCourseById(id)
	if err != nil {
		return nil, err
	}
	if course.TeacherUUID != titularTeacherId {
		return nil, errors.New("the teacher trying to add an aux teacher is not the owner of the course")
	}
	if course.TeacherUUID == auxTeacherId {
		return nil, errors.New("the titular teacher cannot be an aux teacher for his own course")
	}
	if slices.Contains(course.AuxTeachers, auxTeacherId) {
		return nil, errors.New("aux teacher already exists")
	}
	enrolled, err := s.enrollmentRepository.IsEnrolled(auxTeacherId, id)
	if err != nil {
		return nil, err
	}
	if enrolled {
		return nil, errors.New("the aux teacher is already enrolled in the course")
	}
	return s.courseRepository.AddAuxTeacherToCourse(course, auxTeacherId)
}

func (s *CourseService) RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	course, err := s.courseRepository.GetCourseById(id)
	if err != nil {
		return nil, err
	}
	if course.TeacherUUID != titularTeacherId {
		return nil, errors.New("the teacher trying to remove an aux teacher is not the owner of the course")
	}
	if course.TeacherUUID == auxTeacherId {
		return nil, errors.New("the titular teacher cannot be removed as aux teacher from his own course")
	}
	if !slices.Contains(course.AuxTeachers, auxTeacherId) {
		return nil, errors.New("aux teacher is not assigned to this course")
	}
	enrolled, err := s.enrollmentRepository.IsEnrolled(auxTeacherId, id)
	if err != nil {
		return nil, err
	}
	if enrolled {
		return nil, errors.New("the aux teacher is already enrolled in the course")
	}
	return s.courseRepository.RemoveAuxTeacherFromCourse(course, auxTeacherId)
}
