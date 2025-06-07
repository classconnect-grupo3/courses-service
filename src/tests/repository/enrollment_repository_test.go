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
		Feedback:   []model.StudentFeedback{},
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
		Feedback:   []model.StudentFeedback{},
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
		Feedback:   []model.StudentFeedback{},
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
			Feedback:   []model.StudentFeedback{},
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
			Feedback:   []model.StudentFeedback{},
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
		Feedback:   []model.StudentFeedback{},
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
		Feedback:   []model.StudentFeedback{},
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

func TestUnsetFavouriteCourse(t *testing.T) {
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

	// Create enrollment with favourite set to true
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		Favourite:  true, // Initially favourite
		UpdatedAt:  time.Now(),
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Unset the course as favourite
	err = enrollmentRepository.UnsetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Verify the course is now not marked as favourite
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(enrollments))
	assert.False(t, enrollments[0].Favourite)
	assert.Equal(t, "student-123", enrollments[0].StudentID)
}

func TestUnsetFavouriteCourseWithNonExistentEnrollment(t *testing.T) {
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

	// Try to unset favourite for non-existent enrollment
	err = enrollmentRepository.UnsetFavouriteCourse("non-existent-student", createdCourse.ID.Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrollment not found for student")
}

func TestUnsetFavouriteMultipleTimes(t *testing.T) {
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

	// Create enrollment with favourite set to true
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		Favourite:  true,
		UpdatedAt:  time.Now(),
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Unset favourite multiple times (should not error)
	err = enrollmentRepository.UnsetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	err = enrollmentRepository.UnsetFavouriteCourse("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Verify the course is still not marked as favourite
	enrollments, err := enrollmentRepository.GetEnrollmentsByCourseId(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(enrollments))
	assert.False(t, enrollments[0].Favourite)
}

func TestGetEnrollmentsByStudentId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create multiple courses
	course1 := model.Course{
		Title:          "Course 1",
		Description:    "Test Description 1",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse1, err := courseRepository.CreateCourse(course1)
	assert.NoError(t, err)

	course2 := model.Course{
		Title:          "Course 2",
		Description:    "Test Description 2",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse2, err := courseRepository.CreateCourse(course2)
	assert.NoError(t, err)

	course3 := model.Course{
		Title:          "Course 3",
		Description:    "Test Description 3",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse3, err := courseRepository.CreateCourse(course3)
	assert.NoError(t, err)

	// Create enrollments for the same student in multiple courses
	studentID := "student-123"
	enrollments := []model.Enrollment{
		{
			StudentID:  studentID,
			CourseID:   createdCourse1.ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     model.EnrollmentStatusActive,
			Favourite:  true,
			UpdatedAt:  time.Now(),
			Feedback:   []model.StudentFeedback{},
		},
		{
			StudentID:  studentID,
			CourseID:   createdCourse2.ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     model.EnrollmentStatusActive,
			Favourite:  false,
			UpdatedAt:  time.Now(),
			Feedback:   []model.StudentFeedback{},
		},
		{
			StudentID:  studentID,
			CourseID:   createdCourse3.ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     model.EnrollmentStatusCompleted,
			Favourite:  true,
			UpdatedAt:  time.Now(),
			Feedback:   []model.StudentFeedback{},
		},
	}

	// Create the enrollments
	for i, enrollment := range enrollments {
		var course *model.Course
		switch i {
		case 0:
			course = createdCourse1
		case 1:
			course = createdCourse2
		case 2:
			course = createdCourse3
		}
		err = enrollmentRepository.CreateEnrollment(enrollment, course)
		assert.NoError(t, err)
	}

	// Test GetEnrollmentsByStudentId
	retrievedEnrollments, err := enrollmentRepository.GetEnrollmentsByStudentId(studentID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedEnrollments)
	assert.Equal(t, 3, len(retrievedEnrollments))

	// Verify all enrollments belong to the correct student
	courseIDs := make([]string, len(retrievedEnrollments))
	for i, enrollment := range retrievedEnrollments {
		assert.Equal(t, studentID, enrollment.StudentID)
		courseIDs[i] = enrollment.CourseID
	}

	// Verify all courses are present
	expectedCourseIDs := []string{
		createdCourse1.ID.Hex(),
		createdCourse2.ID.Hex(),
		createdCourse3.ID.Hex(),
	}
	for _, expectedCourseID := range expectedCourseIDs {
		assert.Contains(t, courseIDs, expectedCourseID)
	}
}

func TestGetEnrollmentsByStudentIdEmpty(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create a course and enroll a different student
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create enrollment for a different student
	enrollment := model.Enrollment{
		StudentID:  "other-student",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Test GetEnrollmentsByStudentId for a student with no enrollments
	enrollments, err := enrollmentRepository.GetEnrollmentsByStudentId("student-with-no-enrollments")
	assert.NoError(t, err)
	if enrollments != nil {
		assert.Equal(t, 0, len(enrollments))
	} else {
		// Accept nil as valid response when no enrollments exist
		assert.Nil(t, enrollments)
	}
}

func TestGetEnrollmentsByStudentIdWithNonExistentStudent(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Test with non-existent student ID
	enrollments, err := enrollmentRepository.GetEnrollmentsByStudentId("non-existent-student")
	assert.NoError(t, err)
	if enrollments != nil {
		assert.Equal(t, 0, len(enrollments))
	} else {
		// Accept nil as valid response when no enrollments exist
		assert.Nil(t, enrollments)
	}
}

func TestGetEnrollmentsByStudentIdWithDifferentStatuses(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create courses
	courses := []model.Course{
		{
			Title:          "Active Course",
			Description:    "Course with active enrollment",
			Capacity:       10,
			StudentsAmount: 0,
		},
		{
			Title:          "Completed Course",
			Description:    "Course with completed enrollment",
			Capacity:       10,
			StudentsAmount: 0,
		},
		{
			Title:          "Dropped Course",
			Description:    "Course with dropped enrollment",
			Capacity:       10,
			StudentsAmount: 0,
		},
	}

	var createdCourses []*model.Course
	for _, course := range courses {
		createdCourse, err := courseRepository.CreateCourse(course)
		assert.NoError(t, err)
		createdCourses = append(createdCourses, createdCourse)
	}

	// Create enrollments with different statuses
	studentID := "student-123"
	statuses := []model.EnrollmentStatus{
		model.EnrollmentStatusActive,
		model.EnrollmentStatusCompleted,
		model.EnrollmentStatusDropped,
	}

	for i, status := range statuses {
		enrollment := model.Enrollment{
			StudentID:  studentID,
			CourseID:   createdCourses[i].ID.Hex(),
			EnrolledAt: time.Now(),
			Status:     status,
			UpdatedAt:  time.Now(),
			Feedback:   []model.StudentFeedback{},
		}

		err := enrollmentRepository.CreateEnrollment(enrollment, createdCourses[i])
		assert.NoError(t, err)
	}

	// Test GetEnrollmentsByStudentId
	retrievedEnrollments, err := enrollmentRepository.GetEnrollmentsByStudentId(studentID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedEnrollments)
	assert.Equal(t, 3, len(retrievedEnrollments))

	// Verify all different statuses are present
	retrievedStatuses := make([]model.EnrollmentStatus, len(retrievedEnrollments))
	for i, enrollment := range retrievedEnrollments {
		assert.Equal(t, studentID, enrollment.StudentID)
		retrievedStatuses[i] = enrollment.Status
	}

	for _, expectedStatus := range statuses {
		assert.Contains(t, retrievedStatuses, expectedStatus)
	}
}

func TestGetEnrollmentByStudentIdAndCourseId(t *testing.T) {
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
	enrollmentTime := time.Now()
	enrollment := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse.ID.Hex(),
		EnrolledAt: enrollmentTime,
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  enrollmentTime,
		Favourite:  true,
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Test getting enrollment by student ID and course ID
	retrievedEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, retrievedEnrollment)
	assert.Equal(t, "student-123", retrievedEnrollment.StudentID)
	assert.Equal(t, createdCourse.ID.Hex(), retrievedEnrollment.CourseID)
	assert.Equal(t, model.EnrollmentStatusActive, retrievedEnrollment.Status)
	assert.True(t, retrievedEnrollment.Favourite)
	assert.WithinDuration(t, enrollmentTime, retrievedEnrollment.EnrolledAt, time.Second)
	assert.WithinDuration(t, enrollmentTime, retrievedEnrollment.UpdatedAt, time.Second)
}

func TestGetEnrollmentByStudentIdAndCourseIdWithNonExistentEnrollment(t *testing.T) {
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

	// Try to get non-existent enrollment
	retrievedEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("non-existent-student", createdCourse.ID.Hex())
	assert.Error(t, err)
	assert.Nil(t, retrievedEnrollment)
}

func TestGetEnrollmentByStudentIdAndCourseIdWithNonExistentCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Try to get enrollment with non-existent course
	retrievedEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", "non-existent-course")
	assert.Error(t, err)
	assert.Nil(t, retrievedEnrollment)
}

func TestGetEnrollmentByStudentIdAndCourseIdWithMultipleEnrollments(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create two courses
	course1 := model.Course{
		Title:          "Test Course 1",
		Description:    "Test Description 1",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse1, err := courseRepository.CreateCourse(course1)
	assert.NoError(t, err)

	course2 := model.Course{
		Title:          "Test Course 2",
		Description:    "Test Description 2",
		Capacity:       10,
		StudentsAmount: 0,
	}
	createdCourse2, err := courseRepository.CreateCourse(course2)
	assert.NoError(t, err)

	// Create enrollments for same student in different courses
	enrollment1 := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse1.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
		Favourite:  true,
		Feedback:   []model.StudentFeedback{},
	}

	enrollment2 := model.Enrollment{
		StudentID:  "student-123",
		CourseID:   createdCourse2.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
		Favourite:  false,
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment1, createdCourse1)
	assert.NoError(t, err)

	err = enrollmentRepository.CreateEnrollment(enrollment2, createdCourse2)
	assert.NoError(t, err)

	// Get enrollment for course 1
	retrievedEnrollment1, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse1.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, retrievedEnrollment1)
	assert.Equal(t, "student-123", retrievedEnrollment1.StudentID)
	assert.Equal(t, createdCourse1.ID.Hex(), retrievedEnrollment1.CourseID)
	assert.True(t, retrievedEnrollment1.Favourite)

	// Get enrollment for course 2
	retrievedEnrollment2, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse2.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, retrievedEnrollment2)
	assert.Equal(t, "student-123", retrievedEnrollment2.StudentID)
	assert.Equal(t, createdCourse2.ID.Hex(), retrievedEnrollment2.CourseID)
	assert.False(t, retrievedEnrollment2.Favourite)
}

func TestCreateStudentFeedback(t *testing.T) {
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
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Get the created enrollment to get its ID
	createdEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Create student feedback
	feedback := model.StudentFeedback{
		StudentUUID:  "student-123",
		TeacherUUID:  "teacher-456",
		FeedbackType: model.FeedbackTypePositive,
		Score:        85,
		Feedback:     "Excellent work on the assignment!",
		CreatedAt:    time.Now(),
	}

	err = enrollmentRepository.CreateStudentFeedback(feedback, createdEnrollment.ID.Hex())
	assert.NoError(t, err)

	// Verify feedback was added to enrollment
	updatedEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedEnrollment)
	assert.NotEmpty(t, updatedEnrollment.Feedback)
	assert.Equal(t, 1, len(updatedEnrollment.Feedback))

	createdFeedback := updatedEnrollment.Feedback[0]
	assert.Equal(t, "student-123", createdFeedback.StudentUUID)
	assert.Equal(t, "teacher-456", createdFeedback.TeacherUUID)
	assert.Equal(t, model.FeedbackTypePositive, createdFeedback.FeedbackType)
	assert.Equal(t, 85, createdFeedback.Score)
	assert.Equal(t, "Excellent work on the assignment!", createdFeedback.Feedback)
	assert.False(t, createdFeedback.CreatedAt.IsZero())
}

func TestCreateStudentFeedbackWithInvalidEnrollmentID(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	// Create student feedback with invalid enrollment ID
	feedback := model.StudentFeedback{
		StudentUUID:  "student-123",
		TeacherUUID:  "teacher-456",
		FeedbackType: model.FeedbackTypeNegative,
		Score:        60,
		Feedback:     "Needs improvement",
		CreatedAt:    time.Now(),
	}

	err := enrollmentRepository.CreateStudentFeedback(feedback, "invalid-enrollment-id")
	assert.Error(t, err)
}

func TestCreateMultipleStudentFeedbacks(t *testing.T) {
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
		Feedback:   []model.StudentFeedback{},
	}

	err = enrollmentRepository.CreateEnrollment(enrollment, createdCourse)
	assert.NoError(t, err)

	// Get the created enrollment to get its ID
	createdEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)

	// Create multiple feedbacks
	feedbacks := []model.StudentFeedback{
		{
			StudentUUID:  "student-123",
			TeacherUUID:  "teacher-456",
			FeedbackType: model.FeedbackTypePositive,
			Score:        85,
			Feedback:     "Great job on assignment 1!",
			CreatedAt:    time.Now(),
		},
		{
			StudentUUID:  "student-123",
			TeacherUUID:  "teacher-789",
			FeedbackType: model.FeedbackTypeNeutral,
			Score:        75,
			Feedback:     "Good effort on assignment 2",
			CreatedAt:    time.Now(),
		},
		{
			StudentUUID:  "student-123",
			TeacherUUID:  "teacher-456",
			FeedbackType: model.FeedbackTypePositive,
			Score:        90,
			Feedback:     "Excellent work on final project!",
			CreatedAt:    time.Now(),
		},
	}

	// Add all feedbacks
	for _, feedback := range feedbacks {
		err = enrollmentRepository.CreateStudentFeedback(feedback, createdEnrollment.ID.Hex())
		assert.NoError(t, err)
	}

	// Verify all feedbacks were added
	updatedEnrollment, err := enrollmentRepository.GetEnrollmentByStudentIdAndCourseId("student-123", createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedEnrollment)
	assert.Equal(t, 3, len(updatedEnrollment.Feedback))

	// Verify each feedback
	assert.Equal(t, "Great job on assignment 1!", updatedEnrollment.Feedback[0].Feedback)
	assert.Equal(t, model.FeedbackTypePositive, updatedEnrollment.Feedback[0].FeedbackType)
	assert.Equal(t, 85, updatedEnrollment.Feedback[0].Score)

	assert.Equal(t, "Good effort on assignment 2", updatedEnrollment.Feedback[1].Feedback)
	assert.Equal(t, model.FeedbackTypeNeutral, updatedEnrollment.Feedback[1].FeedbackType)
	assert.Equal(t, 75, updatedEnrollment.Feedback[1].Score)

	assert.Equal(t, "Excellent work on final project!", updatedEnrollment.Feedback[2].Feedback)
	assert.Equal(t, model.FeedbackTypePositive, updatedEnrollment.Feedback[2].FeedbackType)
	assert.Equal(t, 90, updatedEnrollment.Feedback[2].Score)
}
