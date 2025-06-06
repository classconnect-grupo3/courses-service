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

type MockEnrollmentRepository struct{}

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
			StudentID: "student-1",
			CourseID:  courseID,
			Favourite: true,
		},
	}, nil
}

func (m *MockEnrollmentRepository) IsEnrolled(studentID, courseID string) (bool, error) {
	// Return true for specific cases to test enrolled scenarios
	if studentID == "enrolled-teacher" {
		return true, nil
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
	if id == "course-with-owner" {
		return &model.Course{
			ID:          primitive.NewObjectID(),
			Title:       "Owner Course",
			Description: "Course with owner",
			TeacherUUID: "owner-teacher",
			AuxTeachers: []string{"aux-teacher-1", "enrolled-teacher"},
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
	err := courseService.DeleteCourse("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
}

func TestDeleteCourseWithEmptyId(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	err := courseService.DeleteCourse("")
	assert.Error(t, err)
}

func TestUpdateCourse(t *testing.T) {
	courseService := service.NewCourseService(&MockCourseRepository{}, &MockEnrollmentRepository{})
	course, err := courseService.UpdateCourse("123e4567-e89b-12d3-a456-426614174000", schemas.UpdateCourseRequest{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
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
		TeacherID:   "123e4567-e89b-12d3-a456-426614174000",
		Capacity:    10,
	})
	assert.Error(t, err)
	assert.Nil(t, course)
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
