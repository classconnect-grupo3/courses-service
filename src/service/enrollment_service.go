package service

import (
	"courses-service/src/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentRepository interface {
	CreateEnrollment(enrollment model.Enrollment, course *model.Course) error
	IsEnrolled(studentID, courseID string) (bool, error)
	DeleteEnrollment(studentID string, course *model.Course) error
}

type EnrollmentService struct {
	enrollmentRepository EnrollmentRepository
	courseRepository     CourseRepository
}

func NewEnrollmentService(enrollmentRepository EnrollmentRepository, courseRepository CourseRepository) *EnrollmentService {
	return &EnrollmentService{enrollmentRepository: enrollmentRepository, courseRepository: courseRepository}
}

func (s *EnrollmentService) EnrollStudent(studentID, courseID string) error {
	// First check if course exists
	course, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return fmt.Errorf("course %s not found for enrollment", courseID)
	}

	// Then check if the course has the capacity to enroll more students
	if course.Capacity <= 0 {
		return fmt.Errorf("course %s is full", courseID)
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
