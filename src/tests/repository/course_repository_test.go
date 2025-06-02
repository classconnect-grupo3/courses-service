package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/tests/testutil"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbSetup *testutil.DBSetup

func init() {
	// Initialize database connection for repository tests
	dbSetup = testutil.SetupTestDB()
}

func TestCreateCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	// Test creating a course
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Verify the course was created
	assert.False(t, createdCourse.ID.IsZero())
	assert.Equal(t, course.Title, createdCourse.Title)
}

func TestGetCourseByTitle(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course for this test
	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	courseRepository.CreateCourse(course)

	gotCourse, err := courseRepository.GetCourseByTitle(course.Title)
	assert.NoError(t, err)

	assert.NotEmpty(t, gotCourse)
	assert.Equal(t, course.Title, gotCourse[0].Title)
}

func TestGetCourseByTitleNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseByTitle("Non-existent course")
	assert.NoError(t, err)

	assert.Empty(t, gotCourse)
}

func TestGetCourseByTeacherId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "123e4567-e89b-12d3-a456-426614174000",
	}

	courseRepository.CreateCourse(course)

	gotCourse, err := courseRepository.GetCourseByTeacherId(course.TeacherUUID)
	assert.NoError(t, err)

	assert.NotEmpty(t, gotCourse)
	assert.Equal(t, course.TeacherUUID, gotCourse[0].TeacherUUID)
}

func TestGetCourseByTeacherIdNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseByTeacherId("Non-existent teacher UUID")
	assert.NoError(t, err)

	assert.Empty(t, gotCourse)
}

func TestGetCourseById(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	objectId, err := primitive.ObjectIDFromHex("663463666666666666666666")
	if err != nil {
		t.Fatalf("Error creating object ID: %v", err)
	}

	course := model.Course{
		ID:          objectId,
		Title:       "Test Course",
		Description: "Test Description",
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	gotCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)

	assert.Equal(t, createdCourse.ID, gotCourse.ID)
	assert.Equal(t, course.Title, gotCourse.Title)
}

func TestGetCourseByIdNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseById("663463666666666666666666")
	assert.Error(t, err)

	assert.Nil(t, gotCourse)
}

func TestGetCourses(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course1 := model.Course{
		Title:       "Test Course 1",
		Description: "Test Description 1",
	}

	course2 := model.Course{
		Title:       "Test Course 2",
		Description: "Test Description 2",
	}

	courseRepository.CreateCourse(course1)
	courseRepository.CreateCourse(course2)

	gotCourses, err := courseRepository.GetCourses()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(gotCourses))

	assert.Equal(t, course1.Title, gotCourses[0].Title)
	assert.Equal(t, course2.Title, gotCourses[1].Title)
}

func TestGetCoursesEmpty(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourses, err := courseRepository.GetCourses()
	assert.NoError(t, err)

	assert.Equal(t, 0, len(gotCourses))
}

func TestDeleteCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	err = courseRepository.DeleteCourse(createdCourse.ID.Hex())
	assert.NoError(t, err)

	gotCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.Error(t, err)

	assert.Nil(t, gotCourse)
}

func TestDeleteCourseNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	err := courseRepository.DeleteCourse("663463666666666666666666")
	assert.NoError(t, err)
}

func TestUpdateCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	expectedUpdatedCourse := model.Course{
		Title:       "Updated Course",
		Description: "Updated Description",
	}

	updatedCourse, err := courseRepository.UpdateCourse(createdCourse.ID.Hex(), expectedUpdatedCourse)
	assert.NoError(t, err)

	assert.Equal(t, expectedUpdatedCourse.Title, updatedCourse.Title)
	assert.Equal(t, expectedUpdatedCourse.Description, updatedCourse.Description)
}

func TestUpdateCourseNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	_, err := courseRepository.UpdateCourse("663463666666666666666666", model.Course{
		Title:       "Updated Course",
		Description: "Updated Description",
	})
	assert.Error(t, err)
}

func TestUpdateCourseOnlyTitle(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.UpdateCourse(createdCourse.ID.Hex(), model.Course{
		Title: "Updated Course",
	})
	assert.NoError(t, err)

	assert.Equal(t, "Updated Course", updatedCourse.Title)
	assert.Equal(t, course.Description, updatedCourse.Description)
}

func TestUpdatedCourseOnlyCapacity(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.UpdateCourse(createdCourse.ID.Hex(), model.Course{
		Capacity: 10,
	})
	assert.NoError(t, err)

	assert.Equal(t, 10, updatedCourse.Capacity)
	assert.Equal(t, course.Title, updatedCourse.Title)
	assert.Equal(t, course.Description, updatedCourse.Description)
}

func TestGetCoursesByStudentId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)
	enrollmentRepository := repository.NewEnrollmentRepository(dbSetup.Client, dbSetup.DBName, courseRepository)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
	}
	resCourse, _ := courseRepository.CreateCourse(course)

	enrollment := model.Enrollment{
		StudentID: "123e4567-e89b-12d3-a456-426614174000",
		CourseID:  resCourse.ID.Hex(),
	}

	fmt.Printf("resCourseId: %v", resCourse.ID.Hex())
	enrollmentRepository.CreateEnrollment(enrollment, resCourse)

	gotCourses, err := courseRepository.GetCoursesByStudentId(enrollment.StudentID)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(gotCourses))
	assert.Equal(t, course.Title, gotCourses[0].Title)
}

func TestGetCoursesByStudentIdEmpty(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
		dbSetup.CleanupCollection("enrollments")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourses, err := courseRepository.GetCoursesByStudentId("non-existent-student")
	assert.NoError(t, err)

	assert.Equal(t, 0, len(gotCourses))
}

func TestAddAuxTeacherToCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "titular-teacher",
		AuxTeachers: []string{},
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.AddAuxTeacherToCourse(createdCourse, "aux-teacher-1")
	assert.NoError(t, err)

	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 1, len(updatedCourse.AuxTeachers))
	assert.Equal(t, "aux-teacher-1", updatedCourse.AuxTeachers[0])
}

func TestAddMultipleAuxTeachersToCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "titular-teacher",
		AuxTeachers: []string{"existing-aux-teacher"},
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.AddAuxTeacherToCourse(createdCourse, "aux-teacher-2")
	assert.NoError(t, err)

	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 2, len(updatedCourse.AuxTeachers))
	assert.Contains(t, updatedCourse.AuxTeachers, "existing-aux-teacher")
	assert.Contains(t, updatedCourse.AuxTeachers, "aux-teacher-2")
}

func TestRemoveAuxTeacherFromCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "titular-teacher",
		AuxTeachers: []string{"aux-teacher-1", "aux-teacher-2"},
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.RemoveAuxTeacherFromCourse(createdCourse, "aux-teacher-1")
	assert.NoError(t, err)

	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 1, len(updatedCourse.AuxTeachers))
	assert.Equal(t, "aux-teacher-2", updatedCourse.AuxTeachers[0])
	assert.NotContains(t, updatedCourse.AuxTeachers, "aux-teacher-1")
}

func TestRemoveLastAuxTeacherFromCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "titular-teacher",
		AuxTeachers: []string{"aux-teacher-1"},
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.RemoveAuxTeacherFromCourse(createdCourse, "aux-teacher-1")
	assert.NoError(t, err)

	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 0, len(updatedCourse.AuxTeachers))
	assert.Equal(t, course.TeacherUUID, updatedCourse.TeacherUUID)
	assert.Equal(t, course.Title, updatedCourse.Title)
	assert.Equal(t, course.Description, updatedCourse.Description)
}

func TestRemoveNonExistentAuxTeacherFromCourse(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		TeacherUUID: "titular-teacher",
		AuxTeachers: []string{"aux-teacher-1"},
	}

	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	updatedCourse, err := courseRepository.RemoveAuxTeacherFromCourse(createdCourse, "non-existent-aux")
	assert.NoError(t, err)

	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 1, len(updatedCourse.AuxTeachers))
	assert.Equal(t, "aux-teacher-1", updatedCourse.AuxTeachers[0])
}

func TestGetCourseByIdWithInvalidId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseById("invalid-object-id")
	assert.Error(t, err)
	assert.Nil(t, gotCourse)
	assert.Contains(t, err.Error(), "failed to get course by id")
}

func TestDeleteCourseWithInvalidId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	err := courseRepository.DeleteCourse("invalid-object-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete course")
}

func TestUpdateCourseWithInvalidId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	_, err := courseRepository.UpdateCourse("invalid-object-id", model.Course{
		Title: "Updated Course",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update course")
}
