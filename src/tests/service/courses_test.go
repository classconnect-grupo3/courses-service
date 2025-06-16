package service_test

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"errors"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockEnrollmentRepository struct{}

// CreateStudentFeedback implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) CreateStudentFeedback(feedbackRequest model.StudentFeedback, enrollmentID string) error {
	return nil
}

// GetFeedbackByStudentId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	return []*model.StudentFeedback{}, nil
}

// GetEnrollmentByStudentIdAndCourseId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) GetEnrollmentByStudentIdAndCourseId(studentID string, courseID string) (*model.Enrollment, error) {
	return &model.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
	}, nil
}

// GetEnrollmentsByStudentId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) GetEnrollmentsByStudentId(studentID string) ([]*model.Enrollment, error) {
	if studentID == "student-with-favourites" {
		return []*model.Enrollment{
			{
				StudentID: studentID,
				CourseID:  "123456789012345678901234", // course-favourite-1
				Favourite: true,
			},
			{
				StudentID: studentID,
				CourseID:  "123456789012345678901236", // course-not-favourite
				Favourite: false,
			},
			{
				StudentID: studentID,
				CourseID:  "123456789012345678901235", // course-favourite-2
				Favourite: true,
			},
		}, nil
	}
	if studentID == "student-no-favourites" {
		return []*model.Enrollment{
			{
				StudentID: studentID,
				CourseID:  "123456789012345678901237", // course-not-favourite-1
				Favourite: false,
			},
			{
				StudentID: studentID,
				CourseID:  "123456789012345678901238", // course-not-favourite-2
				Favourite: false,
			},
		}, nil
	}
	if studentID == "student-no-enrollments" {
		return []*model.Enrollment{}, nil
	}
	if studentID == "error-getting-enrollments" {
		return nil, errors.New("Error getting enrollments")
	}
	return []*model.Enrollment{}, nil
}

// SetFavouriteCourse implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) SetFavouriteCourse(studentID string, courseID string) error {
	return nil
}

// UnsetFavouriteCourse implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) UnsetFavouriteCourse(studentID string, courseID string) error {
	return nil
}

// GetEnrollmentsByCourseId implements repository.EnrollmentRepositoryInterface.
func (m *MockEnrollmentRepository) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	return []*model.Enrollment{
		{
			ID:        primitive.NewObjectID(),
			StudentID: "student-123",
			CourseID:  courseID,
			Status:    model.EnrollmentStatusActive,
			Feedback:  []model.StudentFeedback{}, // Initialize as empty slice
		},
	}, nil
}

func (m *MockEnrollmentRepository) IsEnrolled(studentID, courseID string) (bool, error) {
	// Return true for specific cases to test enrolled scenarios
	if studentID == "enrolled-student" {
		return true, nil
	}
	if studentID == "enrolled-teacher" {
		return true, nil
	}
	if studentID == "error-checking-student" {
		return false, errors.New("Error checking enrollment")
	}
	if studentID == "non-enrolled-student" {
		return false, nil
	}
	return false, nil
}

func (m *MockEnrollmentRepository) CreateEnrollment(enrollment model.Enrollment, course *model.Course) error {
	return nil
}

func (m *MockEnrollmentRepository) DeleteEnrollment(studentID string, course *model.Course) error {
	return nil
}

type MockCourseRepository struct{}

// RemoveAuxTeacherFromCourse implements repository.CourseRepositoryInterface.
func (m *MockCourseRepository) RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

// AddAuxTeacherToCourse implements service.CourseRepository.
func (m *MockCourseRepository) AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

// GetCoursesByStudentId implements service.CourseRepository.
func (m *MockCourseRepository) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	if studentId == "123e4567-e89b-12d3-a456-426614174000" {
		return []*model.Course{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Student Course",
				Description: "Course for student",
				TeacherUUID: "teacher-123",
				AuxTeachers: []string{"aux-teacher-1", "aux-teacher-2"},
				Capacity:    20,
			},
		}, nil
	}
	if studentId == "student-with-favourites" {
		return []*model.Course{
			{
				ID:          mustParseObjectID("course-favourite-1"),
				Title:       "Favourite Course 1",
				Description: "First favourite course",
				TeacherUUID: "teacher-1",
				Capacity:    20,
			},
			{
				ID:          mustParseObjectID("course-not-favourite"),
				Title:       "Not Favourite Course",
				Description: "Course not marked as favourite",
				TeacherUUID: "teacher-2",
				Capacity:    15,
			},
			{
				ID:          mustParseObjectID("course-favourite-2"),
				Title:       "Favourite Course 2",
				Description: "Second favourite course",
				TeacherUUID: "teacher-3",
				Capacity:    25,
			},
		}, nil
	}
	if studentId == "student-no-favourites" {
		return []*model.Course{
			{
				ID:          mustParseObjectID("course-not-favourite-1"),
				Title:       "Course 1",
				Description: "First course not favourite",
				TeacherUUID: "teacher-1",
				Capacity:    20,
			},
			{
				ID:          mustParseObjectID("course-not-favourite-2"),
				Title:       "Course 2",
				Description: "Second course not favourite",
				TeacherUUID: "teacher-2",
				Capacity:    15,
			},
		}, nil
	}
	if studentId == "student-no-enrollments" {
		return []*model.Course{}, nil
	}
	if studentId == "error-getting-courses" {
		return nil, errors.New("Error getting courses")
	}
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) CreateCourse(c model.Course) (*model.Course, error) {
	return &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       c.Title,
		Description: c.Description,
		TeacherUUID: c.TeacherUUID,
		Capacity:    c.Capacity,
	}, nil
}

func (m *MockCourseRepository) GetCourses() ([]*model.Course, error) {
	return []*model.Course{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Test Course 1",
			Description: "Test Description 1",
			TeacherUUID: "teacher-1",
			Capacity:    10,
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Test Course 2",
			Description: "Test Description 2",
			TeacherUUID: "teacher-2",
			Capacity:    15,
		},
	}, nil
}

func (m *MockCourseRepository) GetCourseById(id string) (*model.Course, error) {
	if id == "123e4567-e89b-12d3-a456-426614174000" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Test Course",
			Description: "Test Description",
			TeacherUUID: "titular-teacher",
			AuxTeachers: []string{"existing-aux-teacher"},
		}, nil
	}
	if id == "course-123" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Course 123",
			Description: "Test course for GetCourseMembers",
			TeacherUUID: "teacher-123",
			AuxTeachers: []string{"aux-teacher-1", "aux-teacher-2"},
		}, nil
	}
	if id == "course-with-owner" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Owner Course",
			Description: "Course with owner",
			TeacherUUID: "owner-teacher",
			AuxTeachers: []string{"aux-teacher-1", "enrolled-teacher"},
		}, nil
	}
	if id == "valid-course" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Valid Course",
			Description: "Course for feedback testing",
			TeacherUUID: "teacher-456",
			AuxTeachers: []string{"aux-teacher-2"},
		}, nil
	}
	if id == "course-with-feedback" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Feedback Course",
			Description: "Course with feedback for testing",
			TeacherUUID: "teacher-feedback",
			AuxTeachers: []string{},
		}, nil
	}
	if id == "course-no-feedback" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "No Feedback Course",
			Description: "Course without feedback",
			TeacherUUID: "teacher-no-feedback",
			AuxTeachers: []string{},
		}, nil
	}
	if id == "error-course" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Error Course",
			Description: "Course that causes repository errors",
			TeacherUUID: "teacher-error",
			AuxTeachers: []string{},
		}, nil
	}
	if id == "123e4567-e89b-12d3-a456-426614174001" {
		return nil, nil
	}
	return nil, errors.New("course not found")
}

func (m *MockCourseRepository) DeleteCourse(id string) error {
	return nil
}

func (m *MockCourseRepository) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	if teacherId == "123e4567-e89b-12d3-a456-426614174000" {
		return []*model.Course{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Test Course",
				Description: "Test Description",
				TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
				Capacity:    10,
			},
		}, nil
	}
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) GetCourseByTitle(title string) ([]*model.Course, error) {
	if title == "Test Course" {
		return []*model.Course{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Test Course",
				Description: "Test Description",
				TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
				Capacity:    10,
			},
		}, nil
	}
	return []*model.Course{}, nil
}

func (m *MockCourseRepository) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	return &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	}, nil
}

func (m *MockCourseRepository) UpdateStudentsAmount(courseID string, newStudentsAmount int) error {
	return nil
}

// CreateCourseFeedback implements repository.CourseRepositoryInterface.
func (m *MockCourseRepository) CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error) {
	if courseID == "non-existent-course" {
		return nil, errors.New("Course not found")
	}
	if feedback.StudentUUID == "error-student" {
		return nil, errors.New("Error creating course feedback")
	}

	// Simulate successful creation
	feedback.ID = primitive.NewObjectID()
	feedback.CreatedAt = time.Now()
	return &feedback, nil
}

// GetCourseFeedback implements repository.CourseRepositoryInterface.
func (m *MockCourseRepository) GetCourseFeedback(courseID string, request schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	if courseID == "course-with-feedback" {
		feedbacks := []*model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-1",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent course!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-2",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Average course",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-3",
				FeedbackType: model.FeedbackTypeNegative,
				Score:        2,
				Feedback:     "Needs improvement",
				CreatedAt:    time.Now(),
			},
		}

		// Apply filters if provided
		filteredFeedbacks := []*model.CourseFeedback{}
		for _, feedback := range feedbacks {
			// Filter by feedback type if specified
			if request.FeedbackType != "" && feedback.FeedbackType != request.FeedbackType {
				continue
			}
			// Filter by score range if specified
			if request.StartScore > 0 && feedback.Score < request.StartScore {
				continue
			}
			if request.EndScore > 0 && feedback.Score > request.EndScore {
				continue
			}
			// Filter by date range if specified
			if !request.StartDate.IsZero() && feedback.CreatedAt.Before(request.StartDate) {
				continue
			}
			if !request.EndDate.IsZero() && feedback.CreatedAt.After(request.EndDate) {
				continue
			}
			filteredFeedbacks = append(filteredFeedbacks, feedback)
		}
		return filteredFeedbacks, nil
	}
	if courseID == "course-no-feedback" {
		return []*model.CourseFeedback{}, nil
	}
	if courseID == "error-course" {
		return nil, errors.New("Error getting course feedback")
	}
	if courseID == "non-existent-course" {
		return nil, errors.New("Course not found")
	}
	return []*model.CourseFeedback{}, nil
}

// Helper function to create ObjectID from string
func mustParseObjectID(id string) primitive.ObjectID {
	// For testing purposes, we'll create consistent ObjectIDs
	// In real scenarios, these would be actual MongoDB ObjectIDs
	switch id {
	case "course-favourite-1":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
		return objectID
	case "course-favourite-2":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901235")
		return objectID
	case "course-not-favourite":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901236")
		return objectID
	case "course-not-favourite-1":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901237")
		return objectID
	case "course-not-favourite-2":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901238")
		return objectID
	default:
		return primitive.NewObjectID()
	}
}

func TestCreateCourseWithInvalidCapacity(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.CreateCourse(schemas.CreateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    0,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
}

func TestCreateCourseWithValidCapacity(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	_, err := courseService.CreateCourse(schemas.CreateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	})
	assert.NoError(t, err)
}

func TestGetCourseById(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.GetCourseById("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestGetCourseByIdWithNonExistentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.GetCourseById("123e4567-e89b-12d3-a456-426614174001")
	assert.NoError(t, err)
	assert.Nil(t, course)
}

func TestGetCourseByIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.GetCourseById("")
	assert.Error(t, err)
	assert.Nil(t, course)
}

func TestGetCourseByTeacherId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTeacherId("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(courses))
}

func TestGetCourseByTeacherIdWithNonExistentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTeacherId("123e4567-e89b-12d3-a456-426614174001")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestGetCourseByTeacherIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTeacherId("")
	assert.Error(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestGetCourseByTitle(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTitle("Test Course")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(courses))
}

func TestGetCourseByTitleWithNonExistentTitle(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTitle("Non Existent Title")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestDeleteCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	err := courseService.DeleteCourse("123e4567-e89b-12d3-a456-426614174000", "titular-teacher")
	assert.NoError(t, err)
}

func TestDeleteCourseWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	err := courseService.DeleteCourse("", "titular-teacher")
	assert.Error(t, err)
}

func TestDeleteCourseWithNonOwnerTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	err := courseService.DeleteCourse("123e4567-e89b-12d3-a456-426614174000", "non-owner-teacher")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the user trying to delete the course is not the owner of the course")
}

func TestUpdateCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.UpdateCourse("123e4567-e89b-12d3-a456-426614174000", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "titular-teacher",
		Capacity:    10,
	})
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestUpdateCourseWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.UpdateCourse("", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "titular-teacher",
		Capacity:    10,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
}

func TestUpdateCourseWithNonOwnerTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.UpdateCourse("123e4567-e89b-12d3-a456-426614174000", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "non-owner-teacher",
		Capacity:    10,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the user trying to update the course is not the owner of the course")
}

func TestGetCoursesByStudentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCoursesByStudentId("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(courses))
}

func TestGetCoursesByStudentIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCoursesByStudentId("")
	assert.Error(t, err)
	assert.Nil(t, courses)
}

func TestGetCourses(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourses()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(courses))
}

func TestGetCourseByTitleWithEmptyTitle(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetCourseByTitle("")
	assert.Error(t, err)
	assert.Nil(t, courses)
}

func TestGetCoursesByUserId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	response, err := courseService.GetCoursesByUserId("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Student))
	assert.Equal(t, 1, len(response.Teacher))
}

func TestGetCoursesByUserIdWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	response, err := courseService.GetCoursesByUserId("")
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestAddAuxTeacherToCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("course-with-owner", "owner-teacher", "new-aux-teacher")
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestAddAuxTeacherToCourseWithNonExistentCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("non-existent-course", "owner-teacher", "new-aux-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "course not found")
}

func TestAddAuxTeacherToCourseWithNonOwnerTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("course-with-owner", "non-owner-teacher", "new-aux-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the teacher trying to add an aux teacher is not the owner of the course")
}

func TestAddAuxTeacherToCourseWithTitularTeacherAsAux(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("course-with-owner", "owner-teacher", "owner-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the titular teacher cannot be an aux teacher for his own course")
}

func TestAddAuxTeacherToCourseWithExistingAuxTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("course-with-owner", "owner-teacher", "aux-teacher-1")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "aux teacher already exists")
}

func TestAddAuxTeacherToCourseWithEnrolledTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.AddAuxTeacherToCourse("course-with-owner", "owner-teacher", "enrolled-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "aux teacher already exists")
}

func TestRemoveAuxTeacherFromCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("course-with-owner", "owner-teacher", "aux-teacher-1")
	assert.NoError(t, err)
	assert.NotNil(t, course)
}

func TestRemoveAuxTeacherFromCourseWithNonExistentCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("non-existent-course", "owner-teacher", "aux-teacher-1")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "course not found")
}

func TestRemoveAuxTeacherFromCourseWithNonOwnerTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("course-with-owner", "non-owner-teacher", "aux-teacher-1")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the teacher trying to remove an aux teacher is not the owner of the course")
}

func TestRemoveAuxTeacherFromCourseWithTitularTeacherAsAux(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("course-with-owner", "owner-teacher", "owner-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the titular teacher cannot be removed as aux teacher from his own course")
}

func TestRemoveAuxTeacherFromCourseWithNonAssignedAuxTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("course-with-owner", "owner-teacher", "non-assigned-aux")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "aux teacher is not assigned to this course")
}

func TestRemoveAuxTeacherFromCourseWithEnrolledTeacher(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.RemoveAuxTeacherFromCourse("course-with-owner", "owner-teacher", "enrolled-teacher")
	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "the aux teacher is already enrolled in the course")
}

func TestGetFavouriteCourses(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("student-with-favourites")
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Equal(t, 2, len(courses))

	// Verify that only favourite courses are returned
	titles := make([]string, len(courses))
	for i, course := range courses {
		titles[i] = course.Title
	}
	assert.Contains(t, titles, "Favourite Course 1")
	assert.Contains(t, titles, "Favourite Course 2")
	assert.NotContains(t, titles, "Not Favourite Course")
}

func TestGetFavouriteCoursesWithEmptyStudentId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("")
	assert.Error(t, err)
	assert.Nil(t, courses)
	assert.Contains(t, err.Error(), "studentId is required")
}

func TestGetFavouriteCoursesWithNoFavourites(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("student-no-favourites")
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Equal(t, 0, len(courses))
}

func TestGetFavouriteCoursesWithNoEnrollments(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("student-no-enrollments")
	assert.NoError(t, err)
	assert.NotNil(t, courses)
	assert.Equal(t, 0, len(courses))
}

func TestGetFavouriteCoursesWithErrorGettingEnrollments(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("error-getting-enrollments")
	assert.Error(t, err)
	assert.Nil(t, courses)
	assert.Contains(t, err.Error(), "Error getting enrollments")
}

func TestGetFavouriteCoursesWithErrorGettingCourses(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	courses, err := courseService.GetFavouriteCourses("error-getting-courses")
	assert.Error(t, err)
	assert.Nil(t, courses)
	assert.Contains(t, err.Error(), "Error getting courses")
}

func (m *MockEnrollmentRepository) GetEnrollmentsByCourseID(courseID string) ([]*model.Enrollment, error) {
	return []*model.Enrollment{
		{
			ID:        primitive.NewObjectID(),
			StudentID: "student-1",
			CourseID:  courseID,
			Favourite: true,
			Feedback:  []model.StudentFeedback{}, // Initialize as empty slice
		},
	}, nil
}

func (m *MockEnrollmentRepository) GetStudentFavouriteCourses(studentUUID string) ([]*model.Enrollment, error) {
	return []*model.Enrollment{
		{
			ID:        primitive.NewObjectID(),
			StudentID: studentUUID,
			Favourite: true,
			Feedback:  []model.StudentFeedback{}, // Initialize as empty slice
		},
	}, nil
}

func TestCreateCourseFeedback(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        5,
		FeedbackType: model.FeedbackTypePositive,
		Feedback:     "Excellent course! Very informative.",
	}

	feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, "enrolled-student", feedback.StudentUUID)
	assert.Equal(t, model.FeedbackTypePositive, feedback.FeedbackType)
	assert.Equal(t, 5, feedback.Score)
	assert.Equal(t, "Excellent course! Very informative.", feedback.Feedback)
	assert.False(t, feedback.ID.IsZero())
	assert.False(t, feedback.CreatedAt.IsZero())
}

func TestCreateCourseFeedbackWithNonExistentCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        3,
		FeedbackType: model.FeedbackTypeNeutral,
		Feedback:     "Average course",
	}

	feedback, err := courseService.CreateCourseFeedback("non-existent-course", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "course not found")
}

func TestCreateCourseFeedbackWithInvalidScoreTooLow(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        0, // Too low
		FeedbackType: model.FeedbackTypeNegative,
		Feedback:     "Poor course",
	}

	feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "score must be between 1 and 5")
}

func TestCreateCourseFeedbackWithInvalidScoreTooHigh(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        6, // Too high
		FeedbackType: model.FeedbackTypePositive,
		Feedback:     "Great course",
	}

	feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "score must be between 1 and 5")
}

func TestCreateCourseFeedbackWithValidScoreBoundaries(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	// Test lower boundary (1)
	feedbackRequest1 := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        1,
		FeedbackType: model.FeedbackTypeNegative,
		Feedback:     "Needs significant improvement",
	}

	feedback1, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest1)
	assert.NoError(t, err)
	assert.NotNil(t, feedback1)
	assert.Equal(t, 1, feedback1.Score)

	// Test upper boundary (5)
	feedbackRequest2 := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "enrolled-student",
		Score:        5,
		FeedbackType: model.FeedbackTypePositive,
		Feedback:     "Outstanding course!",
	}

	feedback2, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest2)
	assert.NoError(t, err)
	assert.NotNil(t, feedback2)
	assert.Equal(t, 5, feedback2.Score)
}

func TestCreateCourseFeedbackWithTeacherAsStudent(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "titular-teacher", // This is the teacher UUID
		Score:        4,
		FeedbackType: model.FeedbackTypePositive,
		Feedback:     "Self feedback",
	}

	feedback, err := courseService.CreateCourseFeedback("123e4567-e89b-12d3-a456-426614174000", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "the teacher cannot give feedback to his own course")
}

func TestCreateCourseFeedbackWithAuxTeacherAsStudent(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "existing-aux-teacher", // This is an aux teacher
		Score:        3,
		FeedbackType: model.FeedbackTypeNeutral,
		Feedback:     "Aux teacher feedback",
	}

	feedback, err := courseService.CreateCourseFeedback("123e4567-e89b-12d3-a456-426614174000", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "the teacher cannot give feedback to his own course")
}

func TestCreateCourseFeedbackWithNonEnrolledStudent(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "non-enrolled-student",
		Score:        4,
		FeedbackType: model.FeedbackTypePositive,
		Feedback:     "Great course!",
	}

	feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "the student is not enrolled in the course")
}

func TestCreateCourseFeedbackWithEnrollmentCheckError(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	feedbackRequest := schemas.CreateCourseFeedbackRequest{
		StudentUUID:  "error-checking-student",
		Score:        3,
		FeedbackType: model.FeedbackTypeNeutral,
		Feedback:     "Average course",
	}

	feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "Error checking enrollment")
}

func TestCreateCourseFeedbackWithDifferentFeedbackTypes(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	testCases := []struct {
		feedbackType model.FeedbackType
		score        int
		feedback     string
	}{
		{model.FeedbackTypePositive, 5, "Excellent course!"},
		{model.FeedbackTypeNeutral, 3, "Average course"},
		{model.FeedbackTypeNegative, 1, "Poor course"},
	}

	for _, tc := range testCases {
		feedbackRequest := schemas.CreateCourseFeedbackRequest{
			StudentUUID:  "enrolled-student",
			Score:        tc.score,
			FeedbackType: tc.feedbackType,
			Feedback:     tc.feedback,
		}

		feedback, err := courseService.CreateCourseFeedback("valid-course", feedbackRequest)
		assert.NoError(t, err)
		assert.NotNil(t, feedback)
		assert.Equal(t, tc.feedbackType, feedback.FeedbackType)
		assert.Equal(t, tc.score, feedback.Score)
		assert.Equal(t, tc.feedback, feedback.Feedback)
	}
}

func TestGetCourseFeedback(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}

	feedback, err := courseService.GetCourseFeedback("course-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 3, len(feedback))

	// Verify feedback details
	assert.Equal(t, "student-1", feedback[0].StudentUUID)
	assert.Equal(t, model.FeedbackTypePositive, feedback[0].FeedbackType)
	assert.Equal(t, 5, feedback[0].Score)
	assert.Equal(t, "Excellent course!", feedback[0].Feedback)

	assert.Equal(t, "student-2", feedback[1].StudentUUID)
	assert.Equal(t, model.FeedbackTypeNeutral, feedback[1].FeedbackType)
	assert.Equal(t, 3, feedback[1].Score)
	assert.Equal(t, "Average course", feedback[1].Feedback)

	assert.Equal(t, "student-3", feedback[2].StudentUUID)
	assert.Equal(t, model.FeedbackTypeNegative, feedback[2].FeedbackType)
	assert.Equal(t, 2, feedback[2].Score)
	assert.Equal(t, "Needs improvement", feedback[2].Feedback)
}

func TestGetCourseFeedbackWithNonExistentCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}

	feedback, err := courseService.GetCourseFeedback("non-existent-course", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "course not found")
}

func TestGetCourseFeedbackWithNoFeedback(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}

	feedback, err := courseService.GetCourseFeedback("course-no-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 0, len(feedback))
}

func TestGetCourseFeedbackWithFeedbackTypeFilter(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		FeedbackType: model.FeedbackTypePositive,
	}

	feedback, err := courseService.GetCourseFeedback("course-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 1, len(feedback))
	assert.Equal(t, model.FeedbackTypePositive, feedback[0].FeedbackType)
	assert.Equal(t, "student-1", feedback[0].StudentUUID)
}

func TestGetCourseFeedbackWithScoreRangeFilter(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	// Test filtering for high scores (4-5)
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		StartScore: 4,
		EndScore:   5,
	}

	feedback, err := courseService.GetCourseFeedback("course-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 1, len(feedback))
	assert.Equal(t, 5, feedback[0].Score)
	assert.Equal(t, "student-1", feedback[0].StudentUUID)
}

func TestGetCourseFeedbackWithLowScoreFilter(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	// Test filtering for low scores (1-2)
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		StartScore: 1,
		EndScore:   2,
	}

	feedback, err := courseService.GetCourseFeedback("course-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 1, len(feedback))
	assert.Equal(t, 2, feedback[0].Score)
	assert.Equal(t, "student-3", feedback[0].StudentUUID)
}

func TestGetCourseFeedbackWithCombinedFilters(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	// Test combining feedback type and score filters
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		FeedbackType: model.FeedbackTypeNegative,
		StartScore:   1,
		EndScore:     3,
	}

	feedback, err := courseService.GetCourseFeedback("course-with-feedback", getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedback)
	assert.Equal(t, 1, len(feedback))
	assert.Equal(t, model.FeedbackTypeNegative, feedback[0].FeedbackType)
	assert.Equal(t, 2, feedback[0].Score)
	assert.Equal(t, "student-3", feedback[0].StudentUUID)
}

func TestGetCourseFeedbackWithRepositoryError(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})

	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}

	feedback, err := courseService.GetCourseFeedback("error-course", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedback)
	assert.Contains(t, err.Error(), "error getting course feedback")
}

// Tests for GetCourseMembers
func TestGetCourseMembers(t *testing.T) {
	courseRepo := &MockCourseRepository{}
	enrollmentRepo := &MockEnrollmentRepository{}
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	members, err := courseService.GetCourseMembers("course-123")

	assert.NoError(t, err)
	assert.NotNil(t, members)
	assert.Equal(t, "teacher-123", members.TeacherID)
	assert.Contains(t, members.AuxTeachersIDs, "aux-teacher-1")
	assert.Contains(t, members.AuxTeachersIDs, "aux-teacher-2")
	assert.Contains(t, members.StudentsIDs, "student-123")
	assert.Len(t, members.StudentsIDs, 1)
}

func TestGetCourseMembersWithEmptyCourseId(t *testing.T) {
	courseRepo := &MockCourseRepository{}
	enrollmentRepo := &MockEnrollmentRepository{}
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	members, err := courseService.GetCourseMembers("")

	assert.Error(t, err)
	assert.Nil(t, members)
	assert.Equal(t, "courseId is required", err.Error())
}

func TestGetCourseMembersWithNonExistentCourse(t *testing.T) {
	courseRepo := &MockCourseRepository{}
	enrollmentRepo := &MockEnrollmentRepository{}
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	members, err := courseService.GetCourseMembers("non-existent-course")

	assert.Error(t, err)
	assert.Nil(t, members)
	assert.Contains(t, err.Error(), "course not found")
}

func TestGetCourseMembersWithEnrollmentRepositoryError(t *testing.T) {
	courseRepo := &MockCourseRepository{}
	enrollmentRepo := &MockEnrollmentRepositoryWithError{}
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	members, err := courseService.GetCourseMembers("course-123")

	assert.Error(t, err)
	assert.Nil(t, members)
	assert.Contains(t, err.Error(), "error getting enrollments")
}

func TestGetCourseMembersWithCourseRepositoryError(t *testing.T) {
	courseRepo := &MockCourseRepositoryWithError{}
	enrollmentRepo := &MockEnrollmentRepository{}
	courseService := service.NewCourseService(courseRepo, enrollmentRepo)

	members, err := courseService.GetCourseMembers("error-course")

	assert.Error(t, err)
	assert.Nil(t, members)
	assert.Contains(t, err.Error(), "error getting course")
}

// Mock repository with errors for testing error scenarios
type MockEnrollmentRepositoryWithError struct{}

func (m *MockEnrollmentRepositoryWithError) CreateStudentFeedback(feedbackRequest model.StudentFeedback, enrollmentID string) error {
	return errors.New("error creating feedback")
}

func (m *MockEnrollmentRepositoryWithError) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	return nil, errors.New("error getting feedback")
}

func (m *MockEnrollmentRepositoryWithError) GetEnrollmentByStudentIdAndCourseId(studentID string, courseID string) (*model.Enrollment, error) {
	return nil, errors.New("error getting enrollment")
}

func (m *MockEnrollmentRepositoryWithError) GetEnrollmentsByStudentId(studentID string) ([]*model.Enrollment, error) {
	return nil, errors.New("error getting enrollments")
}

func (m *MockEnrollmentRepositoryWithError) SetFavouriteCourse(studentID string, courseID string) error {
	return errors.New("error setting favourite")
}

func (m *MockEnrollmentRepositoryWithError) UnsetFavouriteCourse(studentID string, courseID string) error {
	return errors.New("error unsetting favourite")
}

func (m *MockEnrollmentRepositoryWithError) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	return nil, errors.New("error getting enrollments")
}

func (m *MockEnrollmentRepositoryWithError) IsEnrolled(studentID, courseID string) (bool, error) {
	return false, errors.New("error checking enrollment")
}

func (m *MockEnrollmentRepositoryWithError) CreateEnrollment(enrollment model.Enrollment, course *model.Course) error {
	return errors.New("error creating enrollment")
}

func (m *MockEnrollmentRepositoryWithError) DeleteEnrollment(studentID string, course *model.Course) error {
	return errors.New("error deleting enrollment")
}

func (m *MockEnrollmentRepositoryWithError) GetEnrollmentsByCourseID(courseID string) ([]*model.Enrollment, error) {
	return nil, errors.New("error getting enrollments")
}

func (m *MockEnrollmentRepositoryWithError) GetStudentFavouriteCourses(studentUUID string) ([]*model.Enrollment, error) {
	return nil, errors.New("error getting favourite courses")
}

// Mock course repository with errors for testing error scenarios
type MockCourseRepositoryWithError struct{}

func (m *MockCourseRepositoryWithError) CreateCourse(c model.Course) (*model.Course, error) {
	return nil, errors.New("error creating course")
}

func (m *MockCourseRepositoryWithError) GetCourses() ([]*model.Course, error) {
	return nil, errors.New("error getting courses")
}

func (m *MockCourseRepositoryWithError) GetCourseById(id string) (*model.Course, error) {
	return nil, errors.New("error getting course")
}

func (m *MockCourseRepositoryWithError) DeleteCourse(id string) error {
	return errors.New("error deleting course")
}

func (m *MockCourseRepositoryWithError) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return nil, errors.New("error getting courses by teacher")
}

func (m *MockCourseRepositoryWithError) GetCourseByTitle(title string) ([]*model.Course, error) {
	return nil, errors.New("error getting courses by title")
}

func (m *MockCourseRepositoryWithError) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	return nil, errors.New("error updating course")
}

func (m *MockCourseRepositoryWithError) UpdateStudentsAmount(courseID string, newStudentsAmount int) error {
	return errors.New("error updating students amount")
}

func (m *MockCourseRepositoryWithError) CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error) {
	return nil, errors.New("error creating feedback")
}

func (m *MockCourseRepositoryWithError) GetCourseFeedback(courseID string, request schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	return nil, errors.New("error getting course feedback")
}

func (m *MockCourseRepositoryWithError) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return nil, errors.New("error getting courses by student")
}

func (m *MockCourseRepositoryWithError) AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return nil, errors.New("error adding aux teacher")
}

func (m *MockCourseRepositoryWithError) RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	return nil, errors.New("error removing aux teacher")
}
