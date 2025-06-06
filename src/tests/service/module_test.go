package service_test

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockModuleRepository struct{}

// GetNextModuleOrder implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) GetNextModuleOrder(courseID string) (int, error) {
	if courseID == "course-with-modules" {
		return 3, nil // Course already has 2 modules
	}
	if courseID == "empty-course" {
		return 1, nil // Empty course
	}
	if courseID == "error-course" {
		return 0, errors.New("Error getting next module order")
	}
	if courseID == "invalid-course-id" {
		return 0, errors.New("invalid course ID")
	}
	return 1, nil
}

// CreateModule implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) CreateModule(courseID string, module model.Module) (*model.Module, error) {
	if courseID == "error-creating-course" {
		return nil, errors.New("Error creating module")
	}
	if courseID == "invalid-course-id" {
		return nil, errors.New("invalid course ID")
	}

	// Simulate successful creation
	module.ID = primitive.NewObjectID()
	module.CourseID = courseID
	module.CreatedAt = time.Now()
	module.UpdatedAt = time.Now()

	return &module, nil
}

// GetModuleById implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) GetModuleById(id string) (*model.Module, error) {
	if id == "valid-module-id" {
		return &model.Module{
			ID:          mustParseModuleObjectID(id),
			Title:       "Test Module",
			Description: "Test Description",
			Order:       1,
			Content:     "Test Content",
			CourseID:    "valid-course-id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if id == "error-module-id" {
		return nil, errors.New("Error getting module by ID")
	}
	if id == "invalid-module-id" {
		return nil, errors.New("invalid module ID")
	}
	return nil, errors.New("module not found")
}

// UpdateModule implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) UpdateModule(id string, module model.Module) (*model.Module, error) {
	if id == "valid-module-id" {
		module.ID = mustParseModuleObjectID(id)
		module.UpdatedAt = time.Now()
		return &module, nil
	}
	if id == "error-updating-module" {
		return nil, errors.New("Error updating module")
	}
	if id == "invalid-module-id" {
		return nil, errors.New("invalid module ID")
	}
	return nil, errors.New("module not found")
}

// DeleteModule implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) DeleteModule(id string) error {
	if id == "valid-module-id" {
		return nil
	}
	if id == "error-deleting-module" {
		return errors.New("Error deleting module")
	}
	if id == "invalid-module-id" {
		return errors.New("invalid module ID")
	}
	return errors.New("module not found")
}

// GetModulesByCourseId implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) GetModulesByCourseId(courseId string) ([]model.Module, error) {
	if courseId == "course-with-modules" {
		return []model.Module{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 1",
				Description: "First module",
				Order:       1,
				Content:     "Content 1",
				CourseID:    courseId,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Content:     "Content 2",
				CourseID:    courseId,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}, nil
	}
	if courseId == "empty-course" {
		return []model.Module{}, nil
	}
	if courseId == "error-course" {
		return nil, errors.New("Error getting modules by course ID")
	}
	if courseId == "invalid-course-id" {
		return nil, errors.New("invalid course ID")
	}
	return []model.Module{}, nil
}

// GetModuleByName implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) GetModuleByName(courseID string, moduleName string) (*model.Module, error) {
	if courseID == "valid-course-id" && moduleName == "Existing Module" {
		return &model.Module{
			ID:          primitive.NewObjectID(),
			Title:       moduleName,
			Description: "Existing module",
			Order:       1,
			Content:     "Existing content",
			CourseID:    courseID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if courseID == "valid-course-id" && moduleName == "Different Module" {
		return &model.Module{
			ID:          mustParseModuleObjectID("different-module-id"),
			Title:       moduleName,
			Description: "Different module",
			Order:       2,
			Content:     "Different content",
			CourseID:    courseID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if courseID == "valid-course-id" && moduleName == "Updated Module Title" {
		// Return a module with the same ID that we're trying to update, to simulate updating the same module
		return &model.Module{
			ID:          mustParseModuleObjectID("valid-module-id"),
			Title:       moduleName,
			Description: "Original description",
			Order:       1,
			Content:     "Original content",
			CourseID:    courseID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if courseID == "error-course" {
		return nil, errors.New("Error getting module by name")
	}
	return nil, errors.New("module not found")
}

// GetModuleByOrder implements repository.ModuleRepositoryInterface.
func (m *MockModuleRepository) GetModuleByOrder(courseID string, order int) (*model.Module, error) {
	if courseID == "valid-course-id" && order == 1 {
		return &model.Module{
			ID:          primitive.NewObjectID(),
			Title:       "First Module",
			Description: "First module description",
			Order:       order,
			Content:     "First module content",
			CourseID:    courseID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}
	if courseID == "error-course" {
		return nil, errors.New("Error getting module by order")
	}
	return nil, errors.New("module not found")
}

// Helper function to create consistent ObjectIDs for testing
func mustParseModuleObjectID(id string) primitive.ObjectID {
	switch id {
	case "valid-module-id":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
		return objectID
	case "different-module-id":
		objectID, _ := primitive.ObjectIDFromHex("123456789012345678901235")
		return objectID
	default:
		return primitive.NewObjectID()
	}
}

// Tests for CreateModule
func TestCreateModule(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	request := schemas.CreateModuleRequest{
		Title:       "New Module",
		Description: "New module description",
		Content:     "New module content",
		CourseID:    "empty-course",
	}

	module, err := moduleService.CreateModule(request)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, request.Title, module.Title)
	assert.Equal(t, request.Description, module.Description)
	assert.Equal(t, request.Content, module.Content)
	assert.Equal(t, request.CourseID, module.CourseID)
	assert.Equal(t, 1, module.Order) // First module in empty course
	assert.False(t, module.ID.IsZero())
}

func TestCreateModuleWithExistingTitle(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	request := schemas.CreateModuleRequest{
		Title:       "Existing Module",
		Description: "New module description",
		Content:     "New module content",
		CourseID:    "valid-course-id",
	}

	module, err := moduleService.CreateModule(request)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module with title Existing Module already exists")
}

func TestCreateModuleWithErrorGettingOrder(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	request := schemas.CreateModuleRequest{
		Title:       "New Module",
		Description: "New module description",
		Content:     "New module content",
		CourseID:    "error-course",
	}

	module, err := moduleService.CreateModule(request)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "Error getting next module order")
}

func TestCreateModuleWithErrorCreating(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	request := schemas.CreateModuleRequest{
		Title:       "New Module",
		Description: "New module description",
		Content:     "New module content",
		CourseID:    "error-creating-course",
	}

	module, err := moduleService.CreateModule(request)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "Error creating module")
}

// Tests for GetModulesByCourseId
func TestGetModulesByCourseId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	modules, err := moduleService.GetModulesByCourseId("course-with-modules")
	assert.NoError(t, err)
	assert.NotNil(t, modules)
	assert.Equal(t, 2, len(modules))
	assert.Equal(t, "Module 1", modules[0].Title)
	assert.Equal(t, "Module 2", modules[1].Title)
}

func TestGetModulesByCourseIdWithEmptyCourse(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	modules, err := moduleService.GetModulesByCourseId("empty-course")
	assert.NoError(t, err)
	assert.NotNil(t, modules)
	assert.Equal(t, 0, len(modules))
}

func TestGetModulesByCourseIdWithEmptyId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	modules, err := moduleService.GetModulesByCourseId("")
	assert.Error(t, err)
	assert.Nil(t, modules)
	assert.Contains(t, err.Error(), "courseId is required")
}

func TestGetModulesByCourseIdWithError(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	modules, err := moduleService.GetModulesByCourseId("error-course")
	assert.Error(t, err)
	assert.Nil(t, modules)
	assert.Contains(t, err.Error(), "Error getting modules by course ID")
}

// Tests for GetModuleById
func TestGetModuleById(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleById("valid-module-id")
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, "Test Module", module.Title)
	assert.Equal(t, "Test Description", module.Description)
	assert.Equal(t, "valid-course-id", module.CourseID)
}

func TestGetModuleByIdWithEmptyId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleById("")
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module id is required")
}

func TestGetModuleByIdWithError(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleById("error-module-id")
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "Error getting module by ID")
}

func TestGetModuleByIdWithInvalidId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleById("invalid-module-id")
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "invalid module ID")
}

// Tests for GetModuleByOrder
func TestGetModuleByOrder(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleByOrder("valid-course-id", 1)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, "First Module", module.Title)
	assert.Equal(t, 1, module.Order)
}

func TestGetModuleByOrderWithEmptyCourseId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleByOrder("", 1)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "courseId is required")
}

func TestGetModuleByOrderWithError(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	module, err := moduleService.GetModuleByOrder("error-course", 1)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "Error getting module by order")
}

// Tests for UpdateModule
func TestUpdateModule(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	updateModule := model.Module{
		ID:          mustParseModuleObjectID("valid-module-id"),
		Title:       "Updated Module Title",
		Description: "Updated description",
		Content:     "Updated content",
		Order:       2,
		CourseID:    "valid-course-id",
	}

	module, err := moduleService.UpdateModule("valid-module-id", updateModule)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, updateModule.Title, module.Title)
	assert.Equal(t, updateModule.Description, module.Description)
	assert.Equal(t, updateModule.Content, module.Content)
}

func TestUpdateModuleWithEmptyId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	updateModule := model.Module{
		Title:       "Updated Module",
		Description: "Updated description",
		Content:     "Updated content",
		Order:       1,
		CourseID:    "valid-course-id",
	}

	module, err := moduleService.UpdateModule("", updateModule)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module id is required")
}

func TestUpdateModuleWithExistingTitleConflict(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	updateModule := model.Module{
		ID:          mustParseModuleObjectID("valid-module-id"),
		Title:       "Different Module", // This title exists but belongs to different module
		Description: "Updated description",
		Content:     "Updated content",
		Order:       1,
		CourseID:    "valid-course-id",
	}

	module, err := moduleService.UpdateModule("valid-module-id", updateModule)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module with title Different Module already exists")
}

func TestUpdateModuleWithErrorGettingByName(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	updateModule := model.Module{
		Title:       "Some Module",
		Description: "Updated description",
		Content:     "Updated content",
		Order:       1,
		CourseID:    "error-course",
	}

	module, err := moduleService.UpdateModule("valid-module-id", updateModule)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "Error getting module by name")
}

func TestUpdateModuleWithErrorUpdating(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	updateModule := model.Module{
		Title:       "Non-existing Module",
		Description: "Updated description",
		Content:     "Updated content",
		Order:       1,
		CourseID:    "valid-course-id",
	}

	module, err := moduleService.UpdateModule("error-updating-module", updateModule)
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module not found")
}

// Tests for DeleteModule
func TestDeleteModule(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	err := moduleService.DeleteModule("valid-module-id")
	assert.NoError(t, err)
}

func TestDeleteModuleWithEmptyId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	err := moduleService.DeleteModule("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "module id is required")
}

func TestDeleteModuleWithError(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	err := moduleService.DeleteModule("error-deleting-module")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error deleting module")
}

func TestDeleteModuleWithInvalidId(t *testing.T) {
	moduleService := service.NewModuleService(&MockModuleRepository{})

	err := moduleService.DeleteModule("invalid-module-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid module ID")
}
