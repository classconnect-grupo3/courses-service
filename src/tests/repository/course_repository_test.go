package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/tests/testutil"
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
