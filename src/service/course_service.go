package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"errors"
	"fmt"
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
		Modules:     []model.Module{},
		AuxTeachers: []string{},
		Feedback:    []model.CourseFeedback{},
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

func (s *CourseService) DeleteCourse(id string, teacherId string) error {
	if id == "" {
		return errors.New("id is required")
	}

	course, err := s.courseRepository.GetCourseById(id)
	if err != nil {
		return err
	}

	if course.TeacherUUID != teacherId {
		return errors.New("the user trying to delete the course is not the owner of the course")
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

	fmt.Printf("ID: %v\n", userId)
	auxTeacherCourses, err := s.courseRepository.GetCoursesByAuxTeacherId(userId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Aux Teacher Courses: %v\n", auxTeacherCourses)
	fmt.Printf("Teacher Courses: %v\n", teacherCourses)
	fmt.Printf("Student Courses: %v\n", studentCourses)

	result.Student = studentCourses
	result.Teacher = teacherCourses
	result.AuxTeacher = auxTeacherCourses

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

	course, err := s.courseRepository.GetCourseById(id)
	if err != nil {
		return nil, err
	}
	if course.TeacherUUID != updateCourseRequest.TeacherID {
		return nil, errors.New("the user trying to update the course is not the owner of the course")
	}
	courseToUpdate := model.Course{
		Title:       updateCourseRequest.Title,
		Description: updateCourseRequest.Description,
		TeacherUUID: updateCourseRequest.TeacherID,
		Capacity:    updateCourseRequest.Capacity,
		UpdatedAt:   time.Now(),
	}
	return s.courseRepository.UpdateCourse(id, courseToUpdate)
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

func (s *CourseService) GetFavouriteCourses(studentId string) ([]*model.Course, error) {
	if studentId == "" {
		return nil, errors.New("studentId is required")
	}

	enrollments, err := s.enrollmentRepository.GetEnrollmentsByStudentId(studentId)
	if err != nil {
		return nil, err
	}

	courses, err := s.courseRepository.GetCoursesByStudentId(studentId)
	if err != nil {
		return nil, err
	}

	favouriteCourses := make([]*model.Course, 0)

	for _, course := range courses {
		for _, enrollment := range enrollments {
			if enrollment.CourseID == course.ID.Hex() && enrollment.Favourite {
				favouriteCourses = append(favouriteCourses, course)
			}
		}
	}
	return favouriteCourses, nil
}

func (s *CourseService) CreateCourseFeedback(courseId string, feedbackRequest schemas.CreateCourseFeedbackRequest) (*model.CourseFeedback, error) {
	course, err := s.courseRepository.GetCourseById(courseId)
	if err != nil {
		return nil, err
	}

	if feedbackRequest.Score < 1 || feedbackRequest.Score > 5 {
		return nil, errors.New("score must be between 1 and 5")
	}

	// Check if the student is the teacher or an aux teacher (should not be allowed to give feedback)
	if course.TeacherUUID == feedbackRequest.StudentUUID || slices.Contains(course.AuxTeachers, feedbackRequest.StudentUUID) {
		return nil, errors.New("the teacher cannot give feedback to his own course")
	}

	if enrolled, err := s.enrollmentRepository.IsEnrolled(feedbackRequest.StudentUUID, courseId); err != nil {
		return nil, err
	} else if !enrolled {
		return nil, errors.New("the student is not enrolled in the course")
	}

	feedback := model.CourseFeedback{
		StudentUUID:  feedbackRequest.StudentUUID,
		FeedbackType: feedbackRequest.FeedbackType,
		Score:        feedbackRequest.Score,
		Feedback:     feedbackRequest.Feedback,
		CreatedAt:    time.Now(),
	}

	return s.courseRepository.CreateCourseFeedback(courseId, feedback)
}

func (s *CourseService) GetCourseFeedback(courseId string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	_, err := s.courseRepository.GetCourseById(courseId)
	if err != nil {
		return nil, errors.New("course not found: " + err.Error())
	}

	feedback, err := s.courseRepository.GetCourseFeedback(courseId, getCourseFeedbackRequest)
	if err != nil {
		return nil, errors.New("error getting course feedback: " + err.Error())
	}

	return feedback, nil
}

func (s *CourseService) GetCourseMembers(courseId string) (*schemas.CourseMembersResponse, error) {
	if courseId == "" {
		return nil, errors.New("courseId is required")
	}

	// Get course to get teacher and aux teachers
	course, err := s.courseRepository.GetCourseById(courseId)
	if err != nil {
		return nil, err
	}

	// Get enrolled students
	enrollments, err := s.enrollmentRepository.GetEnrollmentsByCourseId(courseId)
	if err != nil {
		return nil, err
	}

	// Extract student IDs from enrollments
	var studentIDs []string
	for _, enrollment := range enrollments {
		studentIDs = append(studentIDs, enrollment.StudentID)
	}

	// Build response
	response := &schemas.CourseMembersResponse{
		TeacherID:      course.TeacherUUID,
		AuxTeachersIDs: course.AuxTeachers,
		StudentsIDs:    studentIDs,
	}

	return response, nil
}
