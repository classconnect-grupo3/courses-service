package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"fmt"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentService struct {
	enrollmentRepository repository.EnrollmentRepositoryInterface
	courseRepository     repository.CourseRepositoryInterface
}

func NewEnrollmentService(enrollmentRepository repository.EnrollmentRepositoryInterface, courseRepository repository.CourseRepositoryInterface) *EnrollmentService {
	return &EnrollmentService{enrollmentRepository: enrollmentRepository, courseRepository: courseRepository}
}

func (s *EnrollmentService) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	if courseID == "" {
		return nil, fmt.Errorf("course ID is required")
	}

	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return nil, fmt.Errorf("course %s not found", courseID)
	}

	if course.StudentsAmount <= 0 {
		return []*model.Enrollment{}, nil
	}

	enrollments, err := s.enrollmentRepository.GetEnrollmentsByCourseId(courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting enrollments by course ID: %v", err)
	}

	return enrollments, nil
}

func (s *EnrollmentService) EnrollStudent(studentID, courseID string) error {
	// First check if course exists
	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return fmt.Errorf("course %s not found for enrollment", courseID)
	}

	// Then check if the course has the capacity to enroll more students
	if course.StudentsAmount >= course.Capacity {
		return fmt.Errorf("course %s is full", courseID)
	}

	if course.TeacherUUID == studentID {
		return fmt.Errorf("teacher %s cannot enroll in course %s", studentID, courseID)
	}

	// Then check if the student is already enrolled in the course
	enrolled, err := s.enrollmentRepository.IsEnrolled(studentID, courseID)
	if err != nil {
		return fmt.Errorf("error checking if student %s is enrolled in course %s", studentID, courseID)
	}
	if enrolled {
		return fmt.Errorf("student %s is already enrolled in course %s", studentID, courseID)
	}

	// Then create the enrollment
	enrollment := model.Enrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
		Feedback:   []model.StudentFeedback{},
	}

	err = s.enrollmentRepository.CreateEnrollment(enrollment, course)
	if err != nil {
		return fmt.Errorf("error creating enrollment for student %s in course %s", studentID, courseID)
	}

	return nil
}

func (s *EnrollmentService) UnenrollStudent(studentID, courseID string) error {
	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return fmt.Errorf("course %s not found for unenrollment", courseID)
	}

	if course.StudentsAmount <= 0 {
		return fmt.Errorf("course %s is empty", courseID)
	}

	if course.TeacherUUID == studentID {
		return fmt.Errorf("teacher %s cannot unenroll from course %s", studentID, courseID)
	}

	enrolled, err := s.enrollmentRepository.IsEnrolled(studentID, courseID)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking if student %s is enrolled in course %s", studentID, courseID)
	}
	if !enrolled {
		return fmt.Errorf("student %s is not enrolled in course %s", studentID, courseID)
	}

	err = s.enrollmentRepository.DeleteEnrollment(studentID, course)
	if err != nil {
		return fmt.Errorf("error deleting enrollment for student %s in course %s", studentID, courseID)
	}

	return nil
}

func (s *EnrollmentService) SetFavouriteCourse(studentID, courseID string) error {
	if studentID == "" || courseID == "" {
		return fmt.Errorf("student ID and course ID are required")
	}

	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return fmt.Errorf("course %s not found for favourite course", courseID)
	}

	if course.TeacherUUID == studentID {
		return fmt.Errorf("teacher %s cannot set favourite course %s", studentID, courseID)
	}

	enrolled, err := s.enrollmentRepository.IsEnrolled(studentID, courseID)
	if err != nil {
		return fmt.Errorf("error checking if student %s is enrolled in course %s", studentID, courseID)
	}
	if !enrolled {
		return fmt.Errorf("student %s is not enrolled in course %s", studentID, courseID)
	}

	err = s.enrollmentRepository.SetFavouriteCourse(studentID, courseID)
	if err != nil {
		return fmt.Errorf("error setting favourite course for student %s in course %s", studentID, courseID)
	}

	return nil
}

func (s *EnrollmentService) UnsetFavouriteCourse(studentID, courseID string) error {
	if studentID == "" || courseID == "" {
		return fmt.Errorf("student ID and course ID are required")
	}

	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return fmt.Errorf("course %s not found for unset favourite course", courseID)
	}

	if course.TeacherUUID == studentID {
		return fmt.Errorf("teacher %s cannot unset favourite course %s", studentID, courseID)
	}

	enrolled, err := s.enrollmentRepository.IsEnrolled(studentID, courseID)
	if err != nil {
		return fmt.Errorf("error checking if student %s is enrolled in course %s", studentID, courseID)
	}
	if !enrolled {
		return fmt.Errorf("student %s is not enrolled in course %s", studentID, courseID)
	}

	err = s.enrollmentRepository.UnsetFavouriteCourse(studentID, courseID)
	if err != nil {
		return fmt.Errorf("error unsetting favourite course for student %s in course %s", studentID, courseID)
	}

	return nil
}

func (s *EnrollmentService) GetEnrollmentByStudentIdAndCourseId(studentID, courseID string) (*model.Enrollment, error) {
	if studentID == "" || courseID == "" {
		return nil, fmt.Errorf("student ID and course ID are required")
	}

	enrollment, err := s.enrollmentRepository.GetEnrollmentByStudentIdAndCourseId(studentID, courseID)
	if err != nil {
		return nil, fmt.Errorf("error getting enrollment by student ID and course ID: %v", err)
	}

	return enrollment, nil
}

func (s *EnrollmentService) CreateStudentFeedback(feedbackRequest schemas.CreateStudentFeedbackRequest) error {
	if feedbackRequest.Score < 1 || feedbackRequest.Score > 5 {
		return fmt.Errorf("score must be between 1 and 5, not %d", feedbackRequest.Score)
	}

	enrollment, err := s.GetEnrollmentByStudentIdAndCourseId(feedbackRequest.StudentUUID, feedbackRequest.CourseID)
	if err != nil {
		return err
	}

	course, err := s.courseRepository.GetCourseById(feedbackRequest.CourseID)
	if err != nil {
		return fmt.Errorf("error getting course by ID: %v", err)
	}

	if course.TeacherUUID != feedbackRequest.TeacherUUID && !slices.Contains(course.AuxTeachers, feedbackRequest.TeacherUUID) {
		return fmt.Errorf("teacher %s is not the teacher or aux teacher of course %s", feedbackRequest.TeacherUUID, feedbackRequest.CourseID)
	}

	feedback := model.StudentFeedback{
		StudentUUID:  feedbackRequest.StudentUUID,
		TeacherUUID:  feedbackRequest.TeacherUUID,
		FeedbackType: feedbackRequest.FeedbackType,
		Score:        feedbackRequest.Score,
		Feedback:     feedbackRequest.Feedback,
		CreatedAt:    time.Now(),
	}

	err = s.enrollmentRepository.CreateStudentFeedback(feedback, enrollment.ID.Hex())
	if err != nil {
		return fmt.Errorf("error creating student feedback: %v", err)
	}

	return nil
}

func (s *EnrollmentService) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	if studentID == "" {
		return nil, fmt.Errorf("student ID is required")
	}

	feedback, err := s.enrollmentRepository.GetFeedbackByStudentId(studentID, getFeedbackByStudentIdRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting feedback by student ID: %v", err)
	}

	return feedback, nil
}
