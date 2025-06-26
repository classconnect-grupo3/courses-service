package controller_test

import (
	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"courses-service/src/schemas"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mockService      = &MockCourseService{}
	mockErrorService = &MockCourseServiceWithError{}
	normalController = controller.NewCourseController(mockService, nil, mockActivityService, mockNotificationsQueue)
	errorController  = controller.NewCourseController(mockErrorService, nil, mockActivityService, mockNotificationsQueue)
	normalRouter     = gin.Default()
	errorRouter      = gin.Default()
)

func init() {
	router.InitializeCoursesRoutes(normalRouter, normalController)
	router.InitializeCoursesRoutes(errorRouter, errorController)
}

type MockCourseService struct{}

// GetFavouriteCourses implements service.CourseServiceInterface.
func (m *MockCourseService) GetFavouriteCourses(studentId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

// RemoveAuxTeacherFromCourse implements service.CourseServiceInterface.
func (m *MockCourseService) RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

// AddAuxTeacherToCourse implements controller.CourseService.
func (m *MockCourseService) AddAuxTeacherToCourse(id string, teacherId string, auxTeacherId string) (*model.Course, error) {
	return &model.Course{}, nil
}

// GetCoursesByStudentId implements controller.CourseService.
func (m *MockCourseService) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

// GetCoursesByUserId implements controller.CourseService.
func (m *MockCourseService) GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error) {
	return &schemas.GetCoursesByUserIdResponse{
		Student:    []*model.Course{},
		Teacher:    []*model.Course{},
		AuxTeacher: []*model.Course{},
	}, nil
}

func (m *MockCourseService) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
	course := &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       c.Title,
		Description: c.Description,
		TeacherUUID: c.TeacherID,
		Capacity:    c.Capacity,
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
	}

	return course, nil
}

func (m *MockCourseService) GetCourses() ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockCourseService) GetCourseById(id string) (*model.Course, error) {
	return &model.Course{
		ID:          primitive.NewObjectID(),
		Title:       "Test Course",
		Description: "Test Description",
	}, nil
}

func (m *MockCourseService) DeleteCourse(id string, teacherId string) error {
	return nil
}

func (m *MockCourseService) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockCourseService) GetCourseByTitle(title string) ([]*model.Course, error) {
	return []*model.Course{}, nil
}

func (m *MockCourseService) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
	return &model.Course{}, nil
}

// CreateCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseService) CreateCourseFeedback(courseId string, feedbackRequest schemas.CreateCourseFeedbackRequest) (*model.CourseFeedback, error) {
	return &model.CourseFeedback{
		ID:           primitive.NewObjectID(),
		StudentUUID:  feedbackRequest.StudentUUID,
		FeedbackType: feedbackRequest.FeedbackType,
		Score:        feedbackRequest.Score,
		Feedback:     feedbackRequest.Feedback,
		CreatedAt:    time.Now(),
	}, nil
}

// GetCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseService) GetCourseFeedback(courseId string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	if courseId == "course-with-feedback" {
		feedbacks := []*model.CourseFeedback{
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-1",
				FeedbackType: model.FeedbackTypePositive,
				Score:        5,
				Feedback:     "Great course!",
				CreatedAt:    time.Now(),
			},
			{
				ID:           primitive.NewObjectID(),
				StudentUUID:  "student-2",
				FeedbackType: model.FeedbackTypeNeutral,
				Score:        3,
				Feedback:     "OK course",
				CreatedAt:    time.Now(),
			},
		}

		// Apply filters
		filteredFeedbacks := []*model.CourseFeedback{}
		for _, feedback := range feedbacks {
			if getCourseFeedbackRequest.FeedbackType != "" && feedback.FeedbackType != getCourseFeedbackRequest.FeedbackType {
				continue
			}
			if getCourseFeedbackRequest.StartScore > 0 && feedback.Score < getCourseFeedbackRequest.StartScore {
				continue
			}
			if getCourseFeedbackRequest.EndScore > 0 && feedback.Score > getCourseFeedbackRequest.EndScore {
				continue
			}
			filteredFeedbacks = append(filteredFeedbacks, feedback)
		}
		return filteredFeedbacks, nil
	}
	if courseId == "course-no-feedback" {
		return []*model.CourseFeedback{}, nil
	}
	return []*model.CourseFeedback{}, nil
}

func (m *MockCourseService) GetCourseMembers(courseId string) (*schemas.CourseMembersResponse, error) {
	if courseId == "course-123" {
		return &schemas.CourseMembersResponse{
			TeacherID:      "teacher-123",
			AuxTeachersIDs: []string{"aux-teacher-1", "aux-teacher-2"},
			StudentsIDs:    []string{"student-1", "student-2", "student-3"},
		}, nil
	}
	if courseId == "empty-course" {
		return &schemas.CourseMembersResponse{
			TeacherID:      "teacher-456",
			AuxTeachersIDs: []string{},
			StudentsIDs:    []string{},
		}, nil
	}
	return &schemas.CourseMembersResponse{}, nil
}

type MockCourseServiceWithError struct{}

// GetFavouriteCourses implements service.CourseServiceInterface.
func (m *MockCourseServiceWithError) GetFavouriteCourses(studentId string) ([]*model.Course, error) {
	return nil, errors.New("Error getting favourite courses")
}

// RemoveAuxTeacherFromCourse implements service.CourseServiceInterface.
func (m *MockCourseServiceWithError) RemoveAuxTeacherFromCourse(id string, titularTeacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, errors.New("Error removing aux teacher from course")
}

// AddAuxTeacherToCourse implements controller.CourseService.
func (m *MockCourseServiceWithError) AddAuxTeacherToCourse(id string, teacherId string, auxTeacherId string) (*model.Course, error) {
	return nil, errors.New("Error adding aux teacher to course")
}

// GetCoursesByStudentId implements controller.CourseService.
func (m *MockCourseServiceWithError) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	return nil, errors.New("Error getting courses by student ID")
}

// GetCoursesByUserId implements controller.CourseService.
func (m *MockCourseServiceWithError) GetCoursesByUserId(userId string) (*schemas.GetCoursesByUserIdResponse, error) {
	return nil, errors.New("Error getting courses by user ID")
}

func (m *MockCourseServiceWithError) CreateCourse(c schemas.CreateCourseRequest) (*model.Course, error) {
	return nil, errors.New("Error creating course")
}

func (m *MockCourseServiceWithError) GetCourses() ([]*model.Course, error) {
	return nil, errors.New("Error retrieving courses")
}

func (m *MockCourseServiceWithError) GetCourseById(id string) (*model.Course, error) {
	return nil, errors.New("Error getting course by ID")
}

func (m *MockCourseServiceWithError) DeleteCourse(id string, teacherId string) error {
	return errors.New("Error deleting course")
}

func (m *MockCourseServiceWithError) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	return nil, errors.New("Error getting course by teacher ID")
}

func (m *MockCourseServiceWithError) GetCourseByTitle(title string) ([]*model.Course, error) {
	return nil, errors.New("Error getting course by title")
}

func (m *MockCourseServiceWithError) UpdateCourse(id string, updateCourseRequest schemas.UpdateCourseRequest) (*model.Course, error) {
	return nil, errors.New("Error updating course")
}

// CreateCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseServiceWithError) CreateCourseFeedback(courseId string, feedbackRequest schemas.CreateCourseFeedbackRequest) (*model.CourseFeedback, error) {
	return nil, errors.New("Error creating course feedback")
}

// GetCourseFeedback implements service.CourseServiceInterface.
func (m *MockCourseServiceWithError) GetCourseFeedback(courseId string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	return nil, errors.New("Error getting course feedback")
}

func (m *MockCourseServiceWithError) GetCourseMembers(courseId string) (*schemas.CourseMembersResponse, error) {
	return nil, errors.New("Error getting course members")
}

func TestGetCourses(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestGetCoursesWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error retrieving courses\"}", w.Body.String())
}

func TestCreateCourse(t *testing.T) {
	w := httptest.NewRecorder()
	startTime := time.Now()
	endTime := startTime.Add(time.Second * 10)
	body := `{"title": "Test Course", "description": "Test Description", "teacher_id": "123", "capacity": 10, "start_date": "` + startTime.Format(time.RFC3339) + `", "end_date": "` + endTime.Format(time.RFC3339) + `"}`

	req, _ := http.NewRequest("POST", "/courses", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("POST", "/courses", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestCreateCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	startTime := time.Now()
	endTime := startTime.Add(time.Second * 10)
	body := `{"title": "Test Course", "description": "Test Description", "teacher_id": "123", "capacity": 10, "start_date": "` + startTime.Format(time.RFC3339) + `", "end_date": "` + endTime.Format(time.RFC3339) + `"}`

	req, _ := http.NewRequest("POST", "/courses", strings.NewReader(body))
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error creating course\"}", w.Body.String())
}

func TestGetCourseById(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourseByIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error getting course by ID\"}", w.Body.String())
}

func TestDeleteCourse(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/courses/123?teacherId=teacher123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/courses/123?teacherId=teacher123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteCourseWithoutTeacherId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/courses/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Teacher ID is required")
}

func TestGetCourseByTeacherId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/teacher/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourseByTeacherIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/teacher/123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCoursesByStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/student/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCoursesByStudentIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/student/123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCoursesByUserId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/user/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCoursesByUserIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/user/123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCourseByTitle(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/title/Test Course", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Course not found")
}

func TestUpdateCourse(t *testing.T) {
	w := httptest.NewRecorder()
	startTime := time.Now()
	endTime := startTime.Add(time.Second * 10)
	body := `{"title": "Test Course", "description": "Test Description", "teacher_id": "123", "capacity": 10, "start_date": "` + startTime.Format(time.RFC3339) + `", "end_date": "` + endTime.Format(time.RFC3339) + `"}`

	req, _ := http.NewRequest("PUT", "/courses/123", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid body`

	req, _ := http.NewRequest("PUT", "/courses/123", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	startTime := time.Now()
	endTime := startTime.Add(time.Second * 10)
	body := `{"title": "Test Course", "description": "Test Description", "teacher_id": "123", "capacity": 10, "start_date": "` + startTime.Format(time.RFC3339) + `", "end_date": "` + endTime.Format(time.RFC3339) + `"}`

	req, _ := http.NewRequest("PUT", "/courses/123", strings.NewReader(body))
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddAuxTeacherToCourse(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("POST", "/courses/123/aux-teacher/add", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddAuxTeacherToCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("POST", "/courses/123/aux-teacher/add", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestAddAuxTeacherToCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("POST", "/courses/123/aux-teacher/add", strings.NewReader(body))
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error adding aux teacher to course\"}", w.Body.String())
}

func TestAddAuxTeacherToCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("POST", "/courses//aux-teacher/add", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveAuxTeacherFromCourse(t *testing.T) {
	w := httptest.NewRecorder()
	teacherId := "123"
	auxTeacherId := "456"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/123/aux-teacher/remove?teacherId=%s&auxTeacherId=%s", teacherId, auxTeacherId), nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveAuxTeacherFromCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	teacherId := ""
	auxTeacherId := ""

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/123/aux-teacher/remove?teacherId=%s&auxTeacherId=%s", teacherId, auxTeacherId), nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Teacher ID and aux teacher ID are required")
}

func TestRemoveAuxTeacherFromCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	teacherId := "123"
	auxTeacherId := "456"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses/123/aux-teacher/remove?teacherId=%s&auxTeacherId=%s", teacherId, auxTeacherId), nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error removing aux teacher from course\"}", w.Body.String())
}

func TestRemoveAuxTeacherFromCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	teacherId := "123"
	auxTeacherId := "456"

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/courses//aux-teacher/remove?teacherId=%s&auxTeacherId=%s", teacherId, auxTeacherId), nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFavouriteCourses(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/student/123/favourite", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestGetFavouriteCoursesWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/student/123/favourite", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting favourite courses")
}

func TestGetFavouriteCoursesWithEmptyStudentId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/student//favourite", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Student ID is required")
}

func TestCreateCourseFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"score": 5,
		"feedback_type": "POSITIVO",
		"feedback": "Excellent course! Very informative and well structured."
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-123")
	assert.Contains(t, responseBody, "POSITIVO")
	assert.Contains(t, responseBody, "Excellent course!")
}

func TestCreateCourseFeedbackWithEmptyCourseID(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"score": 4,
		"feedback_type": "POSITIVO",
		"feedback": "Good course"
	}`

	req, _ := http.NewRequest("POST", "/courses//feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestCreateCourseFeedbackWithInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"invalid_field": "invalid"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestCreateCourseFeedbackWithMissingRequiredFields(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{
			name: "Missing student_uuid",
			body: `{
				"score": 5,
				"feedback_type": "POSITIVO",
				"feedback": "Great course"
			}`,
		},
		{
			name: "Missing score",
			body: `{
				"student_uuid": "student-123",
				"feedback_type": "POSITIVO",
				"feedback": "Great course"
			}`,
		},
		{
			name: "Missing feedback_type",
			body: `{
				"student_uuid": "student-123",
				"score": 5,
				"feedback": "Great course"
			}`,
		},
		{
			name: "Missing feedback",
			body: `{
				"student_uuid": "student-123",
				"score": 5,
				"feedback_type": "POSITIVO"
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			normalRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Error:Field validation")
		})
	}
}

func TestCreateCourseFeedbackWithDifferentFeedbackTypes(t *testing.T) {
	testCases := []struct {
		feedbackType string
		score        int
		feedback     string
	}{
		{"POSITIVO", 5, "Excellent course!"},
		{"NEUTRO", 3, "Average course"},
		{"NEGATIVO", 1, "Poor course"},
	}

	for _, tc := range testCases {
		t.Run("FeedbackType_"+tc.feedbackType, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := `{
				"student_uuid": "student-123",
				"score": ` + string(rune(tc.score+'0')) + `,
				"feedback_type": "` + tc.feedbackType + `",
				"feedback": "` + tc.feedback + `"
			}`

			req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			normalRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			responseBody := w.Body.String()
			assert.Contains(t, responseBody, tc.feedbackType)
			assert.Contains(t, responseBody, tc.feedback)
		})
	}
}

func TestCreateCourseFeedbackWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{
		"student_uuid": "student-123",
		"score": 4,
		"feedback_type": "POSITIVO",
		"feedback": "Good course"
	}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error creating course feedback")
}

func TestCreateCourseFeedbackWithScoreBoundaries(t *testing.T) {
	testCases := []struct {
		name  string
		score int
	}{
		{"Score_1", 1},
		{"Score_2", 2},
		{"Score_3", 3},
		{"Score_4", 4},
		{"Score_5", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := `{
				"student_uuid": "student-123",
				"score": ` + string(rune(tc.score+'0')) + `,
				"feedback_type": "POSITIVO",
				"feedback": "Course feedback"
			}`

			req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			normalRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			responseBody := w.Body.String()
			assert.Contains(t, responseBody, "student-123")
		})
	}
}

func TestCreateCourseFeedbackWithInvalidFeedbackType(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"student_uuid": "student-123", "score": 4, "feedback_type": "INVALID_TYPE", "feedback": "Good course"}`

	req, _ := http.NewRequest("POST", "/courses/course-123/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid feedback type")
}

func TestGetCourseFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Verify that response contains feedback data
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-1")
	assert.Contains(t, responseBody, "student-2")
	assert.Contains(t, responseBody, "Great course!")
	assert.Contains(t, responseBody, "OK course")
}

func TestGetCourseFeedbackWithEmptyCourseID(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{}`

	req, _ := http.NewRequest("PUT", "/courses//feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Course ID is required")
}

func TestGetCourseFeedbackWithNoFeedback(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{}`

	req, _ := http.NewRequest("PUT", "/courses/course-no-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Equal(t, "[]", strings.TrimSpace(responseBody))
}

func TestGetCourseFeedbackWithFeedbackTypeFilter(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"feedback_type": "POSITIVO"}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-1")
	assert.Contains(t, responseBody, "Great course!")
	// Should not contain neutral feedback
	assert.NotContains(t, responseBody, "student-2")
}

func TestGetCourseFeedbackWithScoreRangeFilter(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"start_score": 4, "end_score": 5}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-1")
	assert.Contains(t, responseBody, "Great course!")
	// Should not contain lower score feedback
	assert.NotContains(t, responseBody, "student-2")
}

func TestGetCourseFeedbackWithCombinedFilters(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"feedback_type": "POSITIVO", "start_score": 4, "end_score": 5}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-1")
	assert.Contains(t, responseBody, "Great course!")
}

func TestGetCourseFeedbackWithInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"feedback_type": "POSITIVO", "start_score": invalid}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestGetCourseFeedbackWithServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{}`

	req, _ := http.NewRequest("PUT", "/courses/error-course/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting course feedback")
}

func TestGetCourseFeedbackWithDateRangeFilter(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"start_date": "2024-01-01T00:00:00Z", "end_date": "2024-12-31T23:59:59Z"}`

	req, _ := http.NewRequest("PUT", "/courses/course-with-feedback/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "student-1")
	assert.Contains(t, responseBody, "student-2")
}

func TestGetCourseMembers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/courses/course-123/members", nil)
	w := httptest.NewRecorder()
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.CourseMembersResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "teacher-123", response.TeacherID)
	assert.Len(t, response.AuxTeachersIDs, 2)
	assert.Equal(t, "aux-teacher-1", response.AuxTeachersIDs[0])
	assert.Equal(t, "aux-teacher-2", response.AuxTeachersIDs[1])
	assert.Len(t, response.StudentsIDs, 3)
	assert.Equal(t, "student-1", response.StudentsIDs[0])
	assert.Equal(t, "student-2", response.StudentsIDs[1])
	assert.Equal(t, "student-3", response.StudentsIDs[2])
}

func TestGetCourseMembersEmptyCourse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/courses/empty-course/members", nil)
	w := httptest.NewRecorder()
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.CourseMembersResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "teacher-456", response.TeacherID)
	assert.Len(t, response.AuxTeachersIDs, 0)
	assert.Len(t, response.StudentsIDs, 0)
}

func TestGetCourseMembersWithError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/courses/course-123/members", nil)
	w := httptest.NewRecorder()
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Error getting course members", response.Error)
}
