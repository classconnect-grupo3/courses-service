package service_test

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockEnrollmentRepositoryForEnrollmentService struct{}

// GetEnrollmentsByStudentId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) GetEnrollmentsByStudentId(studentID string) ([]*model.Enrollment, error) {
	return nil, nil
}

// GetEnrollmentsByCourseId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	if courseID == "error-course" {
		return nil, errors.New("Error getting enrollments from repository")
	}
	if courseID == "valid-course" {
		return []*model.Enrollment{
			{
				StudentID: "student-1",
				CourseID:  courseID,
				Feedback:  []model.StudentFeedback{},
			},
			{
				StudentID: "student-2",
				CourseID:  courseID,
				Feedback:  []model.StudentFeedback{},
			},
		}, nil
	}
	return []*model.Enrollment{}, nil
}

func (m *MockEnrollmentRepositoryForEnrollmentService) IsEnrolled(studentID, courseID string) (bool, error) {
	// Return specific cases for testing
	if studentID == "already-enrolled-student" {
		return true, nil
	}
	if studentID == "error-checking-student" {
		return false, errors.New("Error checking enrollment")
	}
	if studentID == "error-deleting-student" {
		return true, nil // Make sure this student appears as enrolled so we reach the deletion error
	}
	return false, nil
}

func (m *MockEnrollmentRepositoryForEnrollmentService) CreateEnrollment(enrollment model.Enrollment, course *model.Course) error {
	if enrollment.StudentID == "error-creating-student" {
		return errors.New("Error creating enrollment")
	}
	return nil
}

func (m *MockEnrollmentRepositoryForEnrollmentService) DeleteEnrollment(studentID string, course *model.Course) error {
	if studentID == "error-deleting-student" {
		return errors.New("Error deleting enrollment")
	}
	return nil
}

func (m *MockEnrollmentRepositoryForEnrollmentService) SetFavouriteCourse(studentID, courseID string) error {
	if studentID == "error-setting-favourite-student" {
		return errors.New("Error setting favourite course")
	}
	if studentID == "non-enrolled-student" && courseID == "valid-course" {
		return errors.New("Error setting favourite course for student non-enrolled-student in course valid-course")
	}
	return nil
}

func (m *MockEnrollmentRepositoryForEnrollmentService) UnsetFavouriteCourse(studentID, courseID string) error {
	if studentID == "error-unsetting-favourite-student" {
		return errors.New("Error unsetting favourite course")
	}
	if studentID == "non-enrolled-student" && courseID == "valid-course" {
		return errors.New("Error unsetting favourite course for student non-enrolled-student in course valid-course")
	}
	return nil
}

// GetEnrollmentByStudentIdAndCourseId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) GetEnrollmentByStudentIdAndCourseId(studentID, courseID string) (*model.Enrollment, error) {
	if studentID == "error-student" || courseID == "error-course" {
		return nil, errors.New("Error getting enrollment from repository")
	}
	if studentID == "valid-student" && courseID == "valid-course" {
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Favourite: true,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	if studentID == "non-existent-student" || courseID == "non-existent-course" {
		return nil, errors.New("enrollment not found")
	}
	return &model.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    model.EnrollmentStatusActive,
		Favourite: false,
		Feedback:  []model.StudentFeedback{},
	}, nil
}

// CreateStudentFeedback implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) CreateStudentFeedback(feedback model.StudentFeedback, enrollmentID string) error {
	if feedback.StudentUUID == "error-student" || feedback.TeacherUUID == "error-teacher" {
		return errors.New("Error creating student feedback")
	}
	if enrollmentID == "invalid-enrollment-id" {
		return errors.New("Invalid enrollment ID")
	}
	return nil
}

// GetFeedbackByStudentId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	if studentID == "error-student" {
		return nil, errors.New("Error getting feedback from repository")
	}
	if studentID == "student-with-feedback" {
		return []*model.StudentFeedback{
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-123",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent work!",
			},
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-456",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Good effort",
			},
		}, nil
	}
	return []*model.StudentFeedback{}, nil
}

type MockCourseRepositoryForEnrollment struct{}

func (m *MockCourseRepositoryForEnrollment) GetCourseById(id string) (*model.Course, error) {
	if id == "non-existent-course" {
		return nil, errors.New("Course not found")
	}
	if id == "full-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Full Course",
			Capacity:       2,
			StudentsAmount: 2,
			TeacherUUID:    "teacher-123",
		}, nil
	}
	if id == "teacher-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Teacher Course",
			Capacity:       10,
			StudentsAmount: 1,
			TeacherUUID:    "teacher-student",
		}, nil
	}
	if id == "empty-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Empty Course",
			Capacity:       10,
			StudentsAmount: 0,
			TeacherUUID:    "teacher-123",
		}, nil
	}
	return &model.Course{
		ID:             primitive.NewObjectID(),
		Title:          "Valid Course",
		Capacity:       10,
		StudentsAmount: 1,
		TeacherUUID:    "teacher-123",
		AuxTeachers:    []string{"aux-teacher-456", "aux-teacher-789"},
	}, nil
}

// Mock implementations for other methods (not used in enrollment service)
func (m *MockCourseRepositoryForEnrollment) CreateCourse(c model.Course) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) GetCourses() ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) DeleteCourse(id string) error { return nil }
func (m *MockCourseRepositoryForEnrollment) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) GetCourseByTitle(title string) ([]*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return nil, nil
}
func (m *MockCourseRepositoryForEnrollment) UpdateStudentsAmount(courseID string, newStudentsAmount int) error {
	return nil
}

func TestEnrollStudent(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("valid-student", "valid-course")
	assert.NoError(t, err)
}

func TestEnrollStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for enrollment")
}

func TestEnrollStudentWithFullCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("valid-student", "full-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course full-course is full")
}

func TestEnrollStudentAsTeacher(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot enroll in course teacher-course")
}

func TestEnrollStudentAlreadyEnrolled(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("already-enrolled-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student already-enrolled-student is already enrolled in course valid-course")
}

func TestEnrollStudentWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestEnrollStudentWithErrorCreatingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.EnrollStudent("error-creating-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating enrollment for student error-creating-student in course valid-course")
}

func TestUnenrollStudent(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestUnenrollStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for unenrollment")
}

func TestUnenrollStudentFromEmptyCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("valid-student", "empty-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course empty-course is empty")
}

func TestUnenrollTeacherFromCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot unenroll from course teacher-course")
}

func TestUnenrollStudentNotEnrolled(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled in course valid-course")
}

func TestUnenrollStudentWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestUnenrollStudentWithErrorDeletingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnenrollStudent("error-deleting-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting enrollment for student error-deleting-student in course valid-course")
}

func TestGetEnrollmentsByCourseId(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("valid-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 2, len(enrollments))
	assert.Equal(t, "student-1", enrollments[0].StudentID)
	assert.Equal(t, "student-2", enrollments[1].StudentID)
}

func TestGetEnrollmentsByCourseIdWithEmptyId(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "course ID is required")
}

func TestGetEnrollmentsByCourseIdWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("non-existent-course")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "course non-existent-course not found")
}

func TestGetEnrollmentsByCourseIdWithEmptyCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("empty-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 0, len(enrollments))
}

func TestGetEnrollmentsByCourseIdWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("error-course")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "error getting enrollments by course ID")
}

func TestSetFavouriteCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestSetFavouriteCourseWithEmptyStudentID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestSetFavouriteCourseWithEmptyCourseID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("valid-student", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestSetFavouriteCourseWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for favourite course")
}

func TestSetFavouriteCourseAsTeacher(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot set favourite course teacher-course")
}

func TestSetFavouriteCourseNotEnrolled(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled in course valid-course")
}

func TestSetFavouriteCourseWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestSetFavouriteCourseWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.SetFavouriteCourse("error-setting-favourite-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student error-setting-favourite-student is not enrolled in course valid-course")
}

func TestUnsetFavouriteCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestUnsetFavouriteCourseWithEmptyStudentID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestUnsetFavouriteCourseWithEmptyCourseID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestUnsetFavouriteCourseWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for unset favourite course")
}

func TestUnsetFavouriteCourseAsTeacher(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot unset favourite course teacher-course")
}

func TestUnsetFavouriteCourseNotEnrolled(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled in course valid-course")
}

func TestUnsetFavouriteCourseWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestUnsetFavouriteCourseWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	err := enrollmentService.UnsetFavouriteCourse("error-unsetting-favourite-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student error-unsetting-favourite-student is not enrolled in course valid-course")
}

func TestGetEnrollmentByStudentIdAndCourseId(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("valid-student", "valid-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, "valid-student", enrollment.StudentID)
	assert.Equal(t, "valid-course", enrollment.CourseID)
	assert.Equal(t, model.EnrollmentStatusActive, enrollment.Status)
	assert.True(t, enrollment.Favourite)
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyStudentID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("", "valid-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyCourseID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("valid-student", "")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyBothIDs(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("", "")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithNonExistentEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("non-existent-student", "valid-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "error getting enrollment by student ID and course ID")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("error-student", "valid-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "error getting enrollment by student ID and course ID")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithDefaultCase(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("default-student", "default-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, "default-student", enrollment.StudentID)
	assert.Equal(t, "default-course", enrollment.CourseID)
	assert.Equal(t, model.EnrollmentStatusActive, enrollment.Status)
	assert.False(t, enrollment.Favourite)
}

func TestCreateStudentFeedback(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Great job on the assignment!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.NoError(t, err)
}

func TestCreateStudentFeedbackWithNonExistentEnrollment(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "non-existent-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Great job!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting enrollment by student ID and course ID")
}

func TestCreateStudentFeedbackWithNonExistentCourse(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "non-existent-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Great job!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrollment not found")
}

func TestCreateStudentFeedbackWithUnauthorizedTeacher(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "unauthorized-teacher",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Great job!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher unauthorized-teacher is not the teacher or aux teacher of course valid-course")
}

func TestCreateStudentFeedbackWithTeacherAsAuxTeacher(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "aux-teacher-456", // This teacher is in AuxTeachers array in mock
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypeNeutral,
		Score:        3,
		Feedback:     "Good effort!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.NoError(t, err)
}

func TestCreateStudentFeedbackWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "error-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypeNegative,
		Score:        1,
		Feedback:     "Needs improvement",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting enrollment by student ID and course ID")
}

func TestCreateStudentFeedbackWithDifferentFeedbackTypes(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	testCases := []struct {
		feedbackType model.FeedbackType
		score        int
		feedback     string
	}{
		{model.FeedbackTypePositive, 5, "Excellent work!"},
		{model.FeedbackTypeNeutral, 3, "Good effort"},
		{model.FeedbackTypeNegative, 1, "Needs improvement"},
	}

	for _, tc := range testCases {
		feedbackRequest := schemas.CreateStudentFeedbackRequest{
			StudentUUID:  "valid-student",
			TeacherUUID:  "teacher-123",
			CourseID:     "valid-course",
			FeedbackType: tc.feedbackType,
			Score:        tc.score,
			Feedback:     tc.feedback,
		}

		err := enrollmentService.CreateStudentFeedback(feedbackRequest)
		assert.NoError(t, err)
	}
}

func TestCreateStudentFeedbackWithInvalidScoreTooLow(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        0, // Too low
		Feedback:     "Great work!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "score must be between 1 and 5")
}

func TestCreateStudentFeedbackWithInvalidScoreTooHigh(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        6, // Too high
		Feedback:     "Great work!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "score must be between 1 and 5")
}

func TestCreateStudentFeedbackWithValidScoreBoundaries(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	// Test lower boundary (1)
	feedbackRequest1 := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypeNegative,
		Score:        1,
		Feedback:     "Needs significant improvement",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest1)
	assert.NoError(t, err)

	// Test upper boundary (5)
	feedbackRequest2 := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "valid-student",
		TeacherUUID:  "teacher-123",
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Outstanding work!",
	}

	err = enrollmentService.CreateStudentFeedback(feedbackRequest2)
	assert.NoError(t, err)
}

func TestGetFeedbackByStudentId(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID:     "valid-course",
		FeedbackType: model.FeedbackTypePositive,
		StartScore:   1,
		EndScore:     5,
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("student-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 2, len(feedback))
	assert.Equal(t, "student-with-feedback", feedback[0].StudentUUID)
	assert.Equal(t, "teacher-123", feedback[0].TeacherUUID)
	assert.Equal(t, model.FeedbackTypePositive, feedback[0].FeedbackType)
	assert.Equal(t, 5, feedback[0].Score)
	assert.Equal(t, "Excellent work!", feedback[0].Feedback)
}

func TestGetFeedbackByStudentIdWithEmptyStudentID(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestGetFeedbackByStudentIdWithRepositoryError(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("error-student", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "error getting feedback by student ID")
}

func TestGetFeedbackByStudentIdWithNoFeedback(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("student-without-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 0, len(feedback))
}

func TestGetFeedbackByStudentIdWithFilters(t *testing.T) {
	enrollmentService := service.NewEnrollmentService(&MockEnrollmentRepositoryForEnrollmentService{}, &MockCourseRepositoryForEnrollment{})

	// Test with different filters
	testCases := []struct {
		name            string
		studentID       string
		request         schemas.GetFeedbackByStudentIdRequest
		expectedResults int
	}{
		{
			name:      "Filter by course ID",
			studentID: "student-with-feedback",
			request: schemas.GetFeedbackByStudentIdRequest{
				CourseID: "specific-course",
			},
			expectedResults: 2,
		},
		{
			name:      "Filter by feedback type",
			studentID: "student-with-feedback",
			request: schemas.GetFeedbackByStudentIdRequest{
				FeedbackType: model.FeedbackTypePositive,
			},
			expectedResults: 2,
		},
		{
			name:      "Filter by score range",
			studentID: "student-with-feedback",
			request: schemas.GetFeedbackByStudentIdRequest{
				StartScore: 4,
				EndScore:   5,
			},
			expectedResults: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			feedback, err := enrollmentService.GetFeedbackByStudentId(tc.studentID, tc.request)
			assert.NoError(t, err)
			assert.NotNil(t, feedback)
			assert.Equal(t, tc.expectedResults, len(feedback))
		})
	}
}
