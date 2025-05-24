package controller_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"courses-service/src/controller"
	"courses-service/src/model"
	"courses-service/src/router"
	"courses-service/src/schemas"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mockModuleService      = &MockModuleService{}
	mockModuleErrorService = &MockModuleServiceWithError{}
	normalModuleController = controller.NewModuleController(mockModuleService)
	errorModuleController  = controller.NewModuleController(mockModuleErrorService)
	normalModuleRouter     = gin.Default()
	errorModuleRouter      = gin.Default()
)

func init() {
	gin.SetMode(gin.TestMode)
	router.InitializeModulesRoutes(normalModuleRouter, normalModuleController)
	router.InitializeModulesRoutes(errorModuleRouter, errorModuleController)
}

type MockModuleService struct{}

// DeleteModule implements controller.ModuleService.
func (m *MockModuleService) DeleteModule(id string) error {
	return nil
}

// GetModulesByCourseId implements controller.ModuleService.
func (m *MockModuleService) GetModulesByCourseId(courseId string) ([]model.Module, error) {
	return []model.Module{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Test Module",
			Description: "Test Description",
			Order:       1,
			CourseID:    "123",
		},
	}, nil
}

// UpdateModule implements controller.ModuleService.
func (m *MockModuleService) UpdateModule(id string, module model.Module) (*model.Module, error) {
	return &model.Module{
		ID:          primitive.NewObjectID(),
		Title:       "Test Module",
		Description: "Test Description",
		Order:       1,
		CourseID:    "123",
	}, nil
}

func (m *MockModuleService) CreateModule(module schemas.CreateModuleRequest) (*model.Module, error) {
	return &model.Module{
		ID:          primitive.NewObjectID(),
		Title:       module.Title,
		Description: module.Description,
		Order:       1,
		CourseID:    module.CourseID,
	}, nil
}

func (m *MockModuleService) GetModuleById(id string) (*model.Module, error) {
	return &model.Module{
		ID:          primitive.NewObjectID(),
		Title:       "Test Module",
		Description: "Test Description",
		Order:       1,
		CourseID:    "123",
	}, nil
}

type MockModuleServiceWithError struct{}

// GetModuleById implements controller.ModuleService.
func (m *MockModuleServiceWithError) GetModuleById(id string) (*model.Module, error) {
	return nil, errors.New("Error getting module by id")
}

func (m *MockModuleServiceWithError) DeleteModule(id string) error {
	return errors.New("Error deleting module")
}

func (m *MockModuleServiceWithError) GetModulesByCourseId(courseId string) ([]model.Module, error) {
	return nil, errors.New("Error getting modules by course id")
}

func (m *MockModuleServiceWithError) UpdateModule(id string, module model.Module) (*model.Module, error) {
	return nil, errors.New("Error updating module")
}

func (m *MockModuleServiceWithError) CreateModule(module schemas.CreateModuleRequest) (*model.Module, error) {
	return nil, errors.New("Error creating module")
}

func TestCreateModule(t *testing.T) {
	w := httptest.NewRecorder()

	body := `{"title": "Test Module", "description": "Test Description", "order": 1, "course_id": "123"}`
	req, _ := http.NewRequest("POST", "/modules", strings.NewReader(body))
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateModuleWithInvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"invalid": "body"}`
	req, _ := http.NewRequest("POST", "/modules", strings.NewReader(body))
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")
}

func TestCreateModuleWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"title": "Test Module", "description": "Test Description", "course_id": "123"}`
	req, _ := http.NewRequest("POST", "/modules", strings.NewReader(body))
	errorModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetModulesByCourseId(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/modules/course/123", nil)
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetModulesByCourseIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/modules/course/123", nil)
	errorModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetModuleById(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/modules/123", nil)
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetModuleByIdWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/modules/123", nil)
	errorModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateModule(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"title": "Test Module", "description": "Test Description", "order": 1, "course_id": "123"}`
	req, _ := http.NewRequest("PUT", "/modules/123", strings.NewReader(body))
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateModuleWithInvalidModule(t *testing.T) {
	w := httptest.NewRecorder()
	body := `invalid body`
	req, _ := http.NewRequest("PUT", "/modules/123", strings.NewReader(body))
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateModuleWithError(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"title": "Test Module", "description": "Test Description", "order": 1, "course_id": "123"}`
	req, _ := http.NewRequest("PUT", "/modules/123", strings.NewReader(body))
	errorModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteModule(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/modules/123", nil)
	normalModuleRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteModuleWithError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/modules/123", nil)
	errorModuleRouter.ServeHTTP(w, req)
}