package controller_test

import (
	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"courses-service/src/schemas"
	"errors"
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
	normalController = controller.NewCourseController(mockService)
	errorController  = controller.NewCourseController(mockErrorService)
	normalRouter     = gin.Default()
	errorRouter      = gin.Default()
)

func init() {
	router.InitializeCoursesRoutes(normalRouter, normalController)
	router.InitializeCoursesRoutes(errorRouter, errorController)
}

type MockCourseService struct{}

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
		Student: []*model.Course{},
		Teacher: []*model.Course{},
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

func (m *MockCourseService) DeleteCourse(id string) error {
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

type MockCourseServiceWithError struct{}

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

func (m *MockCourseServiceWithError) DeleteCourse(id string) error {
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
	req, _ := http.NewRequest("DELETE", "/courses/123", nil)
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/courses/123", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourseByTitleWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses/title/Test Course", nil)
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("DELETE", "/courses/123/aux-teacher/remove", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveAuxTeacherFromCourseWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`

	req, _ := http.NewRequest("DELETE", "/courses/123/aux-teacher/remove", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestRemoveAuxTeacherFromCourseWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("DELETE", "/courses/123/aux-teacher/remove", strings.NewReader(body))
	errorRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"Error removing aux teacher from course\"}", w.Body.String())
}

func TestRemoveAuxTeacherFromCourseWithEmptyCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"teacher_id": "123", "aux_teacher_id": "456"}`

	req, _ := http.NewRequest("DELETE", "/courses//aux-teacher/remove", strings.NewReader(body))
	normalRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
