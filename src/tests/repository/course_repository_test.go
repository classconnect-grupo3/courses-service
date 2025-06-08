package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"courses-service/src/tests/testutil"
	"fmt"
	"testing"
	"time"

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
		StudentID:  "student-123",
		CourseID:   resCourse.ID.Hex(),
		EnrolledAt: time.Now(),
		Status:     model.EnrollmentStatusActive,
		UpdatedAt:  time.Now(),
		Feedback:   []model.StudentFeedback{},
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

func TestCreateCourseFeedback(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course first
	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		Capacity:    10,
		TeacherUUID: "teacher-123",
		Feedback:    []model.CourseFeedback{},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create course feedback
	feedback := model.CourseFeedback{
		StudentUUID:  "student-456",
		FeedbackType: model.FeedbackTypePositive,
		Score:        5,
		Feedback:     "Excellent course! Very informative.",
	}

	createdFeedback, err := courseRepository.CreateCourseFeedback(createdCourse.ID.Hex(), feedback)
	assert.NoError(t, err)
	assert.NotNil(t, createdFeedback)
	assert.Equal(t, "student-456", createdFeedback.StudentUUID)
	assert.Equal(t, model.FeedbackTypePositive, createdFeedback.FeedbackType)
	assert.Equal(t, 5, createdFeedback.Score)
	assert.Equal(t, "Excellent course! Very informative.", createdFeedback.Feedback)
	assert.False(t, createdFeedback.ID.IsZero())

	// Verify feedback was added to course
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 1, len(updatedCourse.Feedback))
	assert.Equal(t, createdFeedback.ID, updatedCourse.Feedback[0].ID)
	assert.Equal(t, "student-456", updatedCourse.Feedback[0].StudentUUID)
}

func TestCreateCourseFeedbackWithInvalidCourseID(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Try to create feedback for non-existent course
	feedback := model.CourseFeedback{
		StudentUUID:  "student-456",
		FeedbackType: model.FeedbackTypeNegative,
		Score:        2,
		Feedback:     "Course was disappointing",
	}

	createdFeedback, err := courseRepository.CreateCourseFeedback("non-existent-course", feedback)
	assert.Error(t, err)
	assert.Nil(t, createdFeedback)
}

func TestCreateMultipleCourseFeedbacks(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course first
	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		Capacity:    10,
		TeacherUUID: "teacher-123",
		Feedback:    []model.CourseFeedback{},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Create multiple feedbacks
	feedbacks := []model.CourseFeedback{
		{
			StudentUUID:  "student-1",
			FeedbackType: model.FeedbackTypePositive,
			Score:        5,
			Feedback:     "Great course!",
		},
		{
			StudentUUID:  "student-2",
			FeedbackType: model.FeedbackTypeNeutral,
			Score:        3,
			Feedback:     "Average course",
		},
		{
			StudentUUID:  "student-3",
			FeedbackType: model.FeedbackTypeNegative,
			Score:        2,
			Feedback:     "Could be better",
		},
	}

	// Add all feedbacks
	for _, feedback := range feedbacks {
		createdFeedback, err := courseRepository.CreateCourseFeedback(createdCourse.ID.Hex(), feedback)
		assert.NoError(t, err)
		assert.NotNil(t, createdFeedback)
	}

	// Verify all feedbacks were added
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 3, len(updatedCourse.Feedback))

	// Verify feedback details
	feedbackMap := make(map[string]model.CourseFeedback)
	for _, feedback := range updatedCourse.Feedback {
		feedbackMap[feedback.StudentUUID] = feedback
	}

	assert.Equal(t, "Great course!", feedbackMap["student-1"].Feedback)
	assert.Equal(t, model.FeedbackTypePositive, feedbackMap["student-1"].FeedbackType)
	assert.Equal(t, 5, feedbackMap["student-1"].Score)

	assert.Equal(t, "Average course", feedbackMap["student-2"].Feedback)
	assert.Equal(t, model.FeedbackTypeNeutral, feedbackMap["student-2"].FeedbackType)
	assert.Equal(t, 3, feedbackMap["student-2"].Score)

	assert.Equal(t, "Could be better", feedbackMap["student-3"].Feedback)
	assert.Equal(t, model.FeedbackTypeNegative, feedbackMap["student-3"].FeedbackType)
	assert.Equal(t, 2, feedbackMap["student-3"].Score)
}

func TestCreateCourseFeedbackWithDifferentFeedbackTypes(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course first
	course := model.Course{
		Title:       "Test Course",
		Description: "Test Description",
		Capacity:    10,
		TeacherUUID: "teacher-123",
		Feedback:    []model.CourseFeedback{},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Test different feedback types
	testCases := []struct {
		feedbackType model.FeedbackType
		score        int
		feedback     string
	}{
		{model.FeedbackTypePositive, 5, "Excellent course"},
		{model.FeedbackTypeNeutral, 3, "Okay course"},
		{model.FeedbackTypeNegative, 1, "Poor course"},
	}

	for i, tc := range testCases {
		feedback := model.CourseFeedback{
			StudentUUID:  fmt.Sprintf("student-%d", i+1),
			FeedbackType: tc.feedbackType,
			Score:        tc.score,
			Feedback:     tc.feedback,
		}

		createdFeedback, err := courseRepository.CreateCourseFeedback(createdCourse.ID.Hex(), feedback)
		assert.NoError(t, err)
		assert.NotNil(t, createdFeedback)
		assert.Equal(t, tc.feedbackType, createdFeedback.FeedbackType)
		assert.Equal(t, tc.score, createdFeedback.Score)
		assert.Equal(t, tc.feedback, createdFeedback.Feedback)
	}

	// Verify all feedbacks were added with correct types
	updatedCourse, err := courseRepository.GetCourseById(createdCourse.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, updatedCourse)
	assert.Equal(t, 3, len(updatedCourse.Feedback))
}

func TestGetCourseFeedback(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course with feedback
	course := model.Course{
		Title:       "Feedback Test Course",
		Description: "Course for testing feedback retrieval",
		Capacity:    15,
		TeacherUUID: "teacher-456",
		Feedback: []model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-1",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Amazing course!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-2",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Decent course",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-3",
				FeedbackType: model.FeedbackTypeNegative,
				Score:        2,
				Feedback:     "Disappointing course",
				CreatedAt:    time.Now(),
			},
		},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Get all feedback
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}
	feedbacks, err := courseRepository.GetCourseFeedback(createdCourse.ID.Hex(), getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedbacks)
	assert.Equal(t, 3, len(feedbacks))

	// Verify feedback content
	feedbackMap := make(map[string]*model.CourseFeedback)
	for _, feedback := range feedbacks {
		feedbackMap[feedback.StudentUUID] = feedback
	}

	assert.Equal(t, "Amazing course!", feedbackMap["student-1"].Feedback)
	assert.Equal(t, model.FeedbackTypePositive, feedbackMap["student-1"].FeedbackType)
	assert.Equal(t, 5, feedbackMap["student-1"].Score)

	assert.Equal(t, "Decent course", feedbackMap["student-2"].Feedback)
	assert.Equal(t, model.FeedbackTypeNeutral, feedbackMap["student-2"].FeedbackType)
	assert.Equal(t, 3, feedbackMap["student-2"].Score)

	assert.Equal(t, "Disappointing course", feedbackMap["student-3"].Feedback)
	assert.Equal(t, model.FeedbackTypeNegative, feedbackMap["student-3"].FeedbackType)
	assert.Equal(t, 2, feedbackMap["student-3"].Score)
}

func TestGetCourseFeedbackWithFeedbackTypeFilter(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course with mixed feedback types
	course := model.Course{
		Title:       "Filter Test Course",
		Description: "Course for testing feedback filtering",
		Capacity:    10,
		TeacherUUID: "teacher-789",
		Feedback: []model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-positive",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Great!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-negative",
				FeedbackType: model.FeedbackTypeNegative,
				Score:        1,
				Feedback:     "Terrible!",
				CreatedAt:    time.Now(),
			},
		},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Filter for positive feedback only
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		FeedbackType: model.FeedbackTypePositive,
	}
	feedbacks, err := courseRepository.GetCourseFeedback(createdCourse.ID.Hex(), getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedbacks)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, model.FeedbackTypePositive, feedbacks[0].FeedbackType)
	assert.Equal(t, "student-positive", feedbacks[0].StudentUUID)
}

func TestGetCourseFeedbackWithScoreRangeFilter(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course with various scores
	course := model.Course{
		Title:       "Score Filter Test Course",
		Description: "Course for testing score filtering",
		Capacity:    10,
		TeacherUUID: "teacher-score",
		Feedback: []model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-high",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Perfect score!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-mid",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "Average score",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-low",
				FeedbackType: model.FeedbackTypeNegative,
				Score:        1,
				Feedback:     "Low score",
				CreatedAt:    time.Now(),
			},
		},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Filter for high scores (4-5)
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		StartScore: 4,
		EndScore:   5,
	}
	feedbacks, err := courseRepository.GetCourseFeedback(createdCourse.ID.Hex(), getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedbacks)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, 5, feedbacks[0].Score)
	assert.Equal(t, "student-high", feedbacks[0].StudentUUID)
}

func TestGetCourseFeedbackWithNoResults(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course without feedback
	course := model.Course{
		Title:       "Empty Course",
		Description: "Course with no feedback",
		Capacity:    10,
		TeacherUUID: "teacher-empty",
		Feedback:    []model.CourseFeedback{},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Try to get feedback
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}
	feedbacks, err := courseRepository.GetCourseFeedback(createdCourse.ID.Hex(), getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedbacks)
	assert.Equal(t, 0, len(feedbacks))
}

func TestGetCourseFeedbackWithInvalidCourseID(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Try to get feedback from non-existent course
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{}
	feedbacks, err := courseRepository.GetCourseFeedback("non-existent-course-id", getFeedbackRequest)
	assert.Error(t, err)
	assert.Nil(t, feedbacks)
}

func TestGetCourseFeedbackWithCombinedFilters(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	courseRepository := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a course with diverse feedback
	course := model.Course{
		Title:       "Combined Filter Test",
		Description: "Course for testing combined filters",
		Capacity:    20,
		TeacherUUID: "teacher-combined",
		Feedback: []model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-pos-high",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Excellent positive!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-pos-low",
				FeedbackType: model.FeedbackTypePositive,
				Score:        2,
				Feedback:     "Positive but low score",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-neg-high",
				FeedbackType: model.FeedbackTypeNegative,
				Score:        4,
				Feedback:     "Negative but high score",
				CreatedAt:    time.Now(),
			},
		},
	}
	createdCourse, err := courseRepository.CreateCourse(course)
	assert.NoError(t, err)

	// Filter for positive feedback with high scores (4-5)
	getFeedbackRequest := schemas.GetCourseFeedbackRequest{
		FeedbackType: model.FeedbackTypePositive,
		StartScore:   4,
		EndScore:     5,
	}
	feedbacks, err := courseRepository.GetCourseFeedback(createdCourse.ID.Hex(), getFeedbackRequest)
	assert.NoError(t, err)
	assert.NotNil(t, feedbacks)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, model.FeedbackTypePositive, feedbacks[0].FeedbackType)
	assert.Equal(t, 5, feedbacks[0].Score)
	assert.Equal(t, "student-pos-high", feedbacks[0].StudentUUID)
}
