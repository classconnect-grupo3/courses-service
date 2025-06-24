package service_test

import (
	"context"
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
	if courseID == "enrollment-repo-error-course" {
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
	// For most test cases, return false to allow enrollment
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
	if studentID == "error-student" && courseID == "course-with-enrollment" {
		// Return a valid enrollment so we can test the DisapproveStudent repository error
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	if studentID == "error-student" || courseID == "error-course" {
		return nil, errors.New("Error getting enrollment from repository")
	}
	if studentID == "dropped-student" && courseID == "valid-course" {
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusDropped,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	if studentID == "already-enrolled-student" && courseID == "valid-course" {
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	if studentID == "completed-student" && courseID == "valid-course" {
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusCompleted,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	if studentID == "error-deleting-student" && courseID == "valid-course" {
		return &model.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	// Handle feedback test cases
	if studentID == "student-with-enrollment" && courseID == "course-with-enrollment" {
		return &model.Enrollment{
			ID:        primitive.NewObjectID(),
			StudentID: studentID,
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Favourite: false,
			Feedback:  []model.StudentFeedback{},
		}, nil
	}
	// Special case for specific test scenarios
	if studentID == "valid-student" && courseID == "valid-course" {
		// For the main GetEnrollmentByStudentIdAndCourseId test, return an enrollment
		// But for enrollment tests, this student should not have an existing enrollment
		// We'll differentiate based on method call context - enrollment flow expects nil for new enrollments
		return nil, nil
	}
	if studentID == "non-existent-student" || courseID == "non-existent-course" {
		return nil, errors.New("enrollment not found")
	}
	if studentID == "error-checking-student" {
		return nil, errors.New("Error checking enrollment")
	}
	// For students that aren't enrolled or don't exist, return nil (no enrollment found)
	return nil, nil
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
				CourseID:     "course-123",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent work!",
			},
			{
				StudentUUID:  studentID,
				TeacherUUID:  "teacher-456",
				CourseID:     "course-456",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Good effort",
			},
		}, nil
	}
	return []*model.StudentFeedback{}, nil
}

// ApproveStudent implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) ApproveStudent(studentID, courseID string) error {
	if studentID == "error-student" {
		return errors.New("error approving student")
	}
	if courseID == "error-course-repo" {
		return errors.New("error approving student")
	}
	return nil
}

// DisapproveStudent implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) DisapproveStudent(studentID, courseID, reason string) error {
	if studentID == "error-student" {
		return errors.New("error disapproving student")
	}
	if courseID == "error-course-repo" {
		return errors.New("error disapproving student")
	}
	if studentID == "error-deleting-student" {
		return errors.New("error deleting enrollment")
	}
	return nil
}

// ReactivateDroppedEnrollment implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepositoryForEnrollmentService) ReactivateDroppedEnrollment(studentID, courseID string) error {
	if studentID == "error-student" {
		return errors.New("error reactivating enrollment")
	}
	if courseID == "error-course-repo" {
		return errors.New("error reactivating enrollment")
	}
	return nil
}

type MockCourseRepositoryForEnrollment struct{}

// GetCourseFeedback implements repository.CourseRepositoryInterface.
func (m *MockCourseRepositoryForEnrollment) GetCourseFeedback(courseID string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	panic("unimplemented")
}

// CreateCourseFeedback implements repository.CourseRepositoryInterface.
func (m *MockCourseRepositoryForEnrollment) CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error) {
	return nil, nil
}

func (m *MockCourseRepositoryForEnrollment) GetCourseById(id string) (*model.Course, error) {
	if id == "non-existent-course" {
		return nil, errors.New("Course not found")
	}
	if id == "error-course" {
		return nil, errors.New("Course not found")
	}
	if id == "enrollment-repo-error-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Enrollment Repo Error Course",
			Capacity:       10,
			StudentsAmount: 1,
			TeacherUUID:    "teacher-123",
		}, nil
	}
	if id == "full-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Full Course",
			Capacity:       10,
			StudentsAmount: 10,
			TeacherUUID:    "teacher-123",
		}, nil
	}
	if id == "teacher-course" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Teacher Course",
			Capacity:       10,
			StudentsAmount: 5,
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
	if id == "course-with-enrollment" {
		return &model.Course{
			ID:             primitive.NewObjectID(),
			Title:          "Course with Enrollment",
			Capacity:       20,
			StudentsAmount: 5,
			TeacherUUID:    "teacher-123",
			AuxTeachers:    []string{"aux-teacher-123"},
		}, nil
	}
	// Default case for valid courses
	return &model.Course{
		ID:             primitive.NewObjectID(),
		Title:          "Valid Course",
		Capacity:       10,
		StudentsAmount: 5,
		TeacherUUID:    "teacher-123",
		AuxTeachers:    []string{"aux-teacher-123"},
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
func (m *MockCourseRepositoryForEnrollment) GetCoursesByAuxTeacherId(auxTeacherId string) ([]*model.Course, error) {
	return nil, nil
}

// MockSubmissionRepositoryForEnrollmentService for testing enrollment service
type MockSubmissionRepositoryForEnrollmentService struct{}

func (m *MockSubmissionRepositoryForEnrollmentService) Create(ctx context.Context, submission *model.Submission) error {
	return nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) Update(ctx context.Context, submission *model.Submission) error {
	return nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) GetByID(ctx context.Context, id string) (*model.Submission, error) {
	return nil, nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error) {
	return nil, nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return []model.Submission{}, nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	return []model.Submission{}, nil
}

func (m *MockSubmissionRepositoryForEnrollmentService) DeleteByStudentAndCourse(ctx context.Context, studentUUID, courseID string) error {
	if studentUUID == "error-student" || courseID == "error-course" {
		return errors.New("error deleting submissions")
	}
	return nil
}

// Helper function to create enrollment service with proper dependencies
func createEnrollmentServiceForTests() *service.EnrollmentService {
	enrollmentRepo := &MockEnrollmentRepositoryForEnrollmentService{}
	courseRepo := &MockCourseRepositoryForEnrollment{}
	submissionRepo := &MockSubmissionRepositoryForEnrollmentService{}

	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo, submissionRepo)

	return enrollmentService
}

func TestEnrollStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("valid-student", "valid-course")
	assert.NoError(t, err)
}

func TestEnrollStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for enrollment")
}

func TestEnrollStudentWithFullCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("valid-student", "full-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course full-course is full")
}

func TestEnrollStudentAsTeacher(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot enroll in course teacher-course")
}

func TestEnrollStudentAlreadyEnrolled(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("already-enrolled-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student already-enrolled-student is already enrolled in course valid-course")
}

func TestEnrollStudentWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("error-checking-student", "valid-course")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking existing enrollment for student error-checking-student in course valid-course")
}

func TestEnrollStudentWithErrorCreatingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("error-creating-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating enrollment for student error-creating-student in course valid-course")
}

func TestUnenrollStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestUnenrollStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for unenrollment")
}

func TestUnenrollStudentFromEmptyCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("valid-student", "empty-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course empty-course is empty")
}

func TestUnenrollTeacherFromCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot unenroll from course teacher-course")
}

func TestUnenrollStudentNotEnrolled(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled in course valid-course")
}

func TestUnenrollStudentWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestUnenrollStudentWithErrorDeletingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnenrollStudent("error-deleting-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting enrollment")
}

func TestGetEnrollmentsByCourseId(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("valid-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 2, len(enrollments))
	assert.Equal(t, "student-1", enrollments[0].StudentID)
	assert.Equal(t, "student-2", enrollments[1].StudentID)
}

func TestGetEnrollmentsByCourseIdWithEmptyId(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "course ID is required")
}

func TestGetEnrollmentsByCourseIdWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("non-existent-course")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "course non-existent-course not found")
}

func TestGetEnrollmentsByCourseIdWithEmptyCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("empty-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 0, len(enrollments))
}

func TestGetEnrollmentsByCourseIdWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollments, err := enrollmentService.GetEnrollmentsByCourseId("enrollment-repo-error-course")
	assert.Error(t, err)
	assert.Nil(t, enrollments)
	assert.Contains(t, err.Error(), "error getting enrollments by course ID")
}

func TestSetFavouriteCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestSetFavouriteCourseWithEmptyStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestSetFavouriteCourseWithEmptyCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("valid-student", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestSetFavouriteCourseWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found for favourite course")
}

func TestSetFavouriteCourseAsTeacher(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot set favourite course teacher-course")
}

func TestSetFavouriteCourseNotEnrolled(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled in course valid-course")
}

func TestSetFavouriteCourseWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled in course valid-course")
}

func TestSetFavouriteCourseWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.SetFavouriteCourse("error-setting-favourite-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student error-setting-favourite-student is not enrolled in course valid-course")
}

func TestUnsetFavouriteCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
}

func TestUnsetFavouriteCourseWithEmptyStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestUnsetFavouriteCourseWithEmptyCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestUnsetFavouriteCourseWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course non-existent-course not found")
}

func TestUnsetFavouriteCourseAsTeacher(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("teacher-student", "teacher-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teacher teacher-student cannot unset favourite course")
}

func TestUnsetFavouriteCourseNotEnrolled(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("valid-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student valid-student is not enrolled")
}

func TestUnsetFavouriteCourseWithErrorCheckingEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("error-checking-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error checking if student error-checking-student is enrolled")
}

func TestUnsetFavouriteCourseWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.UnsetFavouriteCourse("error-unsetting-favourite-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student error-unsetting-favourite-student is not enrolled")
}

func TestCreateStudentFeedback(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "student-with-enrollment",
		TeacherUUID:  "teacher-123",
		CourseID:     "course-with-enrollment",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Excellent work!",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)

	assert.NoError(t, err)
}

func TestCreateStudentFeedbackWithNonExistentEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

	feedbackRequest := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "student-with-enrollment",
		TeacherUUID:  "aux-teacher-123",
		CourseID:     "course-with-enrollment",
		FeedbackType: model.FeedbackTypeNeutral,
		Score:        3,
		Feedback:     "Good participation as aux teacher",
	}

	err := enrollmentService.CreateStudentFeedback(feedbackRequest)

	assert.NoError(t, err)
}

func TestCreateStudentFeedbackWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

	testCases := []struct {
		name         string
		feedbackType model.FeedbackType
		score        int
		feedback     string
	}{
		{
			name:         "Positive Feedback",
			feedbackType: model.FeedbackTypePositive,
			score:        5,
			feedback:     "Excellent work!",
		},
		{
			name:         "Negative Feedback",
			feedbackType: model.FeedbackTypeNegative,
			score:        1,
			feedback:     "Needs improvement",
		},
		{
			name:         "Neutral Feedback",
			feedbackType: model.FeedbackTypeNeutral,
			score:        3,
			feedback:     "Average performance",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			feedbackRequest := schemas.CreateStudentFeedbackRequest{
				StudentUUID:  "student-with-enrollment",
				TeacherUUID:  "teacher-123",
				CourseID:     "course-with-enrollment",
				FeedbackType: tc.feedbackType,
				Score:        tc.score,
				Feedback:     tc.feedback,
			}

			err := enrollmentService.CreateStudentFeedback(feedbackRequest)
			assert.NoError(t, err)
		})
	}
}

func TestCreateStudentFeedbackWithInvalidScoreTooLow(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

	// Test minimum valid score (1)
	feedbackRequest1 := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "student-with-enrollment",
		TeacherUUID:  "teacher-123",
		CourseID:     "course-with-enrollment",
		FeedbackType: model.FeedbackTypeNegative,
		Score:        1,
		Feedback:     "Minimum score feedback",
	}

	err1 := enrollmentService.CreateStudentFeedback(feedbackRequest1)
	assert.NoError(t, err1)

	// Test maximum valid score (5)
	feedbackRequest5 := schemas.CreateStudentFeedbackRequest{
		StudentUUID:  "student-with-enrollment",
		TeacherUUID:  "teacher-123",
		CourseID:     "course-with-enrollment",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Maximum score feedback",
	}

	err5 := enrollmentService.CreateStudentFeedback(feedbackRequest5)
	assert.NoError(t, err5)
}

func TestGetFeedbackByStudentId(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

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
	enrollmentService := createEnrollmentServiceForTests()

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestGetFeedbackByStudentIdWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("error-student", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "error getting feedback by student ID")
}

func TestGetFeedbackByStudentIdWithNoFeedback(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	getFeedbackRequest := schemas.GetFeedbackByStudentIdRequest{
		CourseID: "valid-course",
	}

	feedback, err := enrollmentService.GetFeedbackByStudentId("student-without-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 0, len(feedback))
}

func TestGetFeedbackByStudentIdWithFilters(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

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

// Tests for ApproveStudent
func TestApproveStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "valid-course")
	assert.NoError(t, err)
}

func TestApproveStudentWithEmptyStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestApproveStudentWithEmptyCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course ID is required")
}

func TestApproveStudentWithBothEmptyIDs(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestApproveStudentWithWhitespaceStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("   ", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestApproveStudentWithWhitespaceCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "   ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course ID is required")
}

func TestApproveStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course not found")
}

func TestApproveStudentWithCourseRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "error-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course not found")
}

func TestApproveStudentWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("error-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error approving student")
}

func TestApproveStudentWithRepositoryErrorFromCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.ApproveStudent("valid-student", "error-course-repo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error approving student")
}

func TestApproveStudentWithValidUUIDs(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	testCases := []struct {
		name      string
		studentID string
		courseID  string
		shouldErr bool
	}{
		{
			name:      "Valid UUID format student and course",
			studentID: "123e4567-e89b-12d3-a456-426614174000",
			courseID:  "987fcdeb-51c2-43d4-b567-531028391849",
			shouldErr: false,
		},
		{
			name:      "Valid standard string IDs",
			studentID: "student-123",
			courseID:  "course-456",
			shouldErr: false,
		},
		{
			name:      "Mixed valid formats",
			studentID: "123e4567-e89b-12d3-a456-426614174000",
			courseID:  "course-456",
			shouldErr: false,
		},
		{
			name:      "Very long valid IDs",
			studentID: "very-long-student-id-with-many-characters-that-should-still-work",
			courseID:  "very-long-course-id-with-many-characters-that-should-still-work",
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := enrollmentService.ApproveStudent(tc.studentID, tc.courseID)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApproveStudentWithSpecialCharacters(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	testCases := []struct {
		name      string
		studentID string
		courseID  string
		shouldErr bool
	}{
		{
			name:      "IDs with underscores",
			studentID: "student_123",
			courseID:  "course_456",
			shouldErr: false,
		},
		{
			name:      "IDs with hyphens",
			studentID: "student-123",
			courseID:  "course-456",
			shouldErr: false,
		},
		{
			name:      "IDs with dots",
			studentID: "student.123",
			courseID:  "course.456",
			shouldErr: false,
		},
		{
			name:      "IDs with numbers",
			studentID: "123456",
			courseID:  "789012",
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := enrollmentService.ApproveStudent(tc.studentID, tc.courseID)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApproveStudentCourseValidation(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	// Test that course validation happens before repository call
	err := enrollmentService.ApproveStudent("valid-student", "non-existent-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course not found")
}

func TestApproveStudentErrorPropagation(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	// Test that repository errors are properly wrapped
	err := enrollmentService.ApproveStudent("error-student", "valid-course")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error approving student")
}

func TestApproveStudentServiceLayerValidation(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	// Test input sanitization - trimming whitespace
	err := enrollmentService.ApproveStudent("  valid-student  ", "  valid-course  ")
	assert.NoError(t, err) // Should succeed since our implementation now trims whitespace
}

// Tests for DisapproveStudent functionality
func TestDisapproveStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "course-with-enrollment", "Did not meet course requirements")

	assert.NoError(t, err)
}

func TestDisapproveStudentWithEmptyStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("", "valid-course", "Valid reason")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student ID is required")
}

func TestDisapproveStudentWithEmptyCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "", "Valid reason")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course ID is required")
}

func TestDisapproveStudentWithEmptyReason(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "course-with-enrollment", "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reason is required")
}

func TestDisapproveStudentWithWhitespaceReason(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "course-with-enrollment", "   ")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reason is required")
}

func TestDisapproveStudentWithNonExistentCourse(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "non-existent-course", "Valid reason")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course not found")
}

func TestDisapproveStudentWithCourseRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.DisapproveStudent("student-with-enrollment", "error-course", "Valid reason")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "course not found")
}

func TestDisapproveStudentWithSpecialCharactersInReason(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	specialReason := "Student failed to meet requirements: @#$%^&*()_+{}|:<>?[]\\;',./"
	err := enrollmentService.DisapproveStudent("student-with-enrollment", "course-with-enrollment", specialReason)

	assert.NoError(t, err)
}

// Test for re-enrollment of dropped students
func TestEnrollDroppedStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("dropped-student", "valid-course")

	assert.NoError(t, err) // Should succeed by reactivating the dropped enrollment
}

// Test for enrolling a student who completed the course (should fail)
func TestEnrollCompletedStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("completed-student", "valid-course")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has already completed course")
}

// Test with a student that doesn't exist (different from error-checking-student)
func TestEnrollStudentWithNewStudent(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	err := enrollmentService.EnrollStudent("new-student", "valid-course")

	assert.NoError(t, err) // Should succeed creating a new enrollment
}

func TestGetEnrollmentByStudentIdAndCourseId(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("already-enrolled-student", "valid-course")
	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, "already-enrolled-student", enrollment.StudentID)
	assert.Equal(t, "valid-course", enrollment.CourseID)
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyStudentID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("", "valid-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyCourseID(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("valid-student", "")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithEmptyBothIDs(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("", "")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "student ID and course ID are required")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithNonExistentEnrollment(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("non-existent-student", "non-existent-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "enrollment not found")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithRepositoryError(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("error-student", "error-course")
	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "Error getting enrollment from repository")
}

func TestGetEnrollmentByStudentIdAndCourseIdWithDefaultCase(t *testing.T) {
	enrollmentService := createEnrollmentServiceForTests()

	enrollment, err := enrollmentService.GetEnrollmentByStudentIdAndCourseId("default-student", "default-course")
	assert.NoError(t, err)
	assert.Nil(t, enrollment) // Should return nil for non-existent enrollments
}
