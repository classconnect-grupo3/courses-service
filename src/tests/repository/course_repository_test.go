package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/tests/testutil"
	"testing"

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
	if err != nil {
		t.Fatalf("Error creating course: %v", err)
	}

	// Verify the course was created
	if createdCourse.ID.IsZero() {
		t.Error("Expected course to have an ID but got zero ID")
	}
	if createdCourse.Title != course.Title {
		t.Errorf("Expected title %s but got %s", course.Title, createdCourse.Title)
	}
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
	if err != nil {
		t.Fatalf("Error getting course: %v", err)
	}

	if len(gotCourse) == 0 {
		t.Error("Expected course but got empty array")
	}

	if gotCourse[0].Title != course.Title {
		t.Errorf("Expected title %s but got %s", course.Title, gotCourse[0].Title)
	}
}

func TestGetCourseByTitleNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseByTitle("Non-existent course")
	if err != nil {
		t.Fatalf("Error getting course: %v", err)
	}

	if len(gotCourse) != 0 {
		t.Error("Expected empty array but got courses")
	}
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
	if err != nil {
		t.Fatalf("Error getting course: %v", err)
	}

	if len(gotCourse) == 0 {
		t.Error("Expected course but got empty array")
	}

	if gotCourse[0].TeacherUUID != course.TeacherUUID {
		t.Errorf("Expected teacher UUID %s but got %s", course.TeacherUUID, gotCourse[0].TeacherUUID)
	}
}

func TestGetCourseByTeacherIdNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseByTeacherId("Non-existent teacher UUID")
	if err != nil {
		t.Fatalf("Error getting course: %v", err)
	}

	if len(gotCourse) != 0 {
		t.Error("Expected empty array but got courses")
	}
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
	if err != nil {
		t.Fatalf("Error creating course: %v", err)
	}

	gotCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Error getting course: %v", err)
	}

	if gotCourse.ID != createdCourse.ID {
		t.Errorf("Expected course ID %s but got %s", createdCourse.ID.Hex(), gotCourse.ID.Hex())
	}

	if gotCourse.Title != course.Title {
		t.Errorf("Expected title %s but got %s", course.Title, gotCourse.Title)
	}
}

func TestGetCourseByIdNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourse, err := courseRepository.GetCourseById("663463666666666666666666")
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	if gotCourse != nil {
		t.Error("Expected nil but got course")
	}
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
	if err != nil {
		t.Fatalf("Error getting courses: %v", err)
	}

	if len(gotCourses) != 2 {
		t.Errorf("Expected 2 courses but got %d", len(gotCourses))
	}

	if gotCourses[0].Title != course1.Title {
		t.Errorf("Expected title %s but got %s", course1.Title, gotCourses[0].Title)
	}

	if gotCourses[1].Title != course2.Title {
		t.Errorf("Expected title %s but got %s", course2.Title, gotCourses[1].Title)
	}
}

func TestGetCoursesEmpty(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	gotCourses, err := courseRepository.GetCourses()
	if err != nil {
		t.Fatalf("Error getting courses: %v", err)
	}

	if len(gotCourses) != 0 {
		t.Errorf("Expected 0 courses but got %d", len(gotCourses))
	}
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
	if err != nil {
		t.Fatalf("Error creating course: %v", err)
	}

	err = courseRepository.DeleteCourse(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Error deleting course: %v", err)
	}

	gotCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	if gotCourse != nil {
		t.Error("Expected nil but got course")
	}
}

func TestDeleteCourseNotFound(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	err := courseRepository.DeleteCourse("663463666666666666666666")
	if err != nil {
		t.Fatalf("Expected nil but got error: %v", err)
	}
}
