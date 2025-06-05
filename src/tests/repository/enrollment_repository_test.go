package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnrollment(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create enrollment
	enrollment := model.Enrollment{
		StudentID:  "123e4567-e89b-12d3-a456-426614174000",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Verify enrollment was created and course capacity updated
	enrolled, err := enrollmentRepository.IsEnrolled(enrollment.StudentID, enrollment.CourseID)
	assert.NoError(t, err)
	assert.True(t, enrolled)

	// Verify course capacity was updated
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedCourse.StudentsAmount)
}

func TestIsEnrolled(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Test not enrolled initially
	enrolled, err := enrollmentRepository.IsEnrolled("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, enrolled)

	// Create enrollment
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Test enrolled after creation
	enrolled, err = enrollmentRepository.IsEnrolled("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.True(t, enrolled)
}

func TestIsEnrolledWithNonExistentEnrollment(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	enrolled, err := enrollmentRepository.IsEnrolled("non-existent-student", "non-existent-course")
	assert.NoError(t, err)
	assert.False(t, enrolled)
}

func TestDeleteEnrollment(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create enrollment
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Verify enrollment exists
	enrolled, err := enrollmentRepository.IsEnrolled("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.True(t, enrolled)

	// Get updated course for deletion (with correct StudentsAmount)
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Delete enrollment
	err = enrollmentRepository.DeleteEnrollment("student-123", updatedCourse)
	assert.NoError(t, err)

	// Verify enrollment was deleted
	enrolled, err = enrollmentRepository.IsEnrolled("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, enrolled)

	// Verify course capacity was updated
	finalCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 0, finalCourse.StudentsAmount)
}

func TestDeleteNonExistentEnrollment(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 1,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Try to delete non-existent enrollment (should not error)
	err = enrollmentRepository.DeleteEnrollment("non-existent-student", createdCourse)
	assert.NoError(t, err)

	// Verify course capacity was NOT changed (should remain 1)
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedCourse.StudentsAmount)
}

func TestMultipleEnrollments(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create multiple enrollments
	students := []string{"student-1", "student-2", "student-3"}

	for _, studentID := range students {
		enrollment := model.Enrollment{
			StudentID:  studentID,
			CourseID:   createdCourse.ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     model.EnrollmentStatusActive,
			UpdatedAt:  time.Now(),
		}

		// Get the current course state before each enrollment
		currentCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
		assert.NoError(t, err)

		err = enrollmentRepository.CreateEnrollment(enrollment, currentCourse)
		assert.NoError(t, err)

		// Verify each enrollment
		enrolled, err := enrollmentRepository.IsEnrolled(studentID, createdCourse.ID.Hex())
		assert.NoError(t, err)
		assert.True(t, enrolled)
	}

	// Verify final course capacity
	finalCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 3, finalCourse.StudentsAmount)
}

func TestGetEnrollmentsByCourseId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create multiple enrollments
	students := []string{"student-1", "student-2", "student-3"}

	for _, studentID := range students {
		enrollment := model.Enrollment{
			StudentID:  studentID,
			CourseID:   createdCourse.ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     model.EnrollmentStatusActive,
			UpdatedAt:  time.Now(),
		}

		// Get the current course state before each enrollment
		currentCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
		assert.NoError(t, err)

		err = enrollmentRepository.CreateEnrollment(enrollment, currentCourse)
		assert.NoError(t, err)
	}

	// Test GetEnrollmentsByCourseId
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 3, len(enrollments))

	// Verify all students are in the list
	studentIDs := make([]string, len(enrollments))
	for i, enrollment := range enrollments {
		studentIDs[i] = enrollment.StudentID
		assert.Equal(t, createdCourse.ID.Hex(), enrollment.CourseID)
		assert.Equal(t, model.EnrollmentStatusActive, enrollment.Status)
	}

	for _, expectedStudent := range students {
		assert.Contains(t, studentIDs, expectedStudent)
	}
}

func TestGetEnrollmentsByCourseIdEmpty(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course but no enrollments
	course := model.Course{
		Title:          "Empty Course",
		Description:    "Course with no enrollments",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Test GetEnrollmentsByCourseId with empty course
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 0, len(enrollments))
}

func TestGetEnrollmentsByCourseIdWithNonExistentCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Test with non-existent course ID
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId("663463666666666666666666")
	assert.NoError(t, err)
	assert.NotNil(t, enrollments)
	assert.Equal(t, 0, len(enrollments))
}

func TestSetFavouriteCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create enrollment
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		Favourite:  false, // Initially not favourite
		UpdatedAt:  time.Now(),
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Set the course as favourite
	err = enrollmentRepository.SetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Verify the course is now marked as favourite
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(enrollments))
	assert.True(t, enrollments[0].Favourite)
	assert.Equal(t, "student-123", enrollments[0].StudentID)
}

func TestSetFavouriteCourseWithNonExistentEnrollment(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course but no enrollment
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Try to set favourite for non-existent enrollment
	err = enrollmentRepository.SetFavouriteCourse("non-existent-student", createdCourse.ID.Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrollment not found for student")
}

func TestSetFavouriteMultipleTimes(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course first
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create enrollment
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		Favourite:  false,
		UpdatedAt:  time.Now(),
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Set favourite multiple times (should not error)
	err = enrollmentRepository.SetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	err = enrollmentRepository.SetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Verify the course is still marked as favourite
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(enrollments))
	assert.True(t, enrollments[0].Favourite)
}
