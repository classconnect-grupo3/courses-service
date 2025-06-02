package repository_test

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModuleRepositoryMock struct{}

// Helper function to create a test course with modules
func createTestCourseWithModules(t *testing.T, courseRepo *repository.CourseRepository) *model.Course {
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		TeacherUUID:    "teacher-123",
		TeacherName:    "Test Teacher",
		Capacity:       30,
		StudentsAmount: 0,
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(24 * time.Hour * 30),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Modules: []model.Module{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 1",
				Description: "First module",
				Order:       1,
				Content:     "Content 1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Content:     "Content 2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	createdCourse, err := courseRepo.CreateCourse(course)
	if err != nil {
		t.Fatalf("Failed to create test course: %v", err)
	}
	return createdCourse
}

// Helper function to create an empty test course
func createEmptyTestCourse(t *testing.T, courseRepo *repository.CourseRepository) *model.Course {
	course := model.Course{
		Title:          "Empty Test Course",
		Description:    "Test Description",
		TeacherUUID:    "teacher-123",
		TeacherName:    "Test Teacher",
		Capacity:       30,
		StudentsAmount: 0,
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(24 * time.Hour * 30),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Modules:        []model.Module{},
	}

	createdCourse, err := courseRepo.CreateCourse(course)
	if err != nil {
		t.Fatalf("Failed to create empty test course: %v", err)
	}
	return createdCourse
}

func TestGetNextModuleOrder(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Test with empty course
	emptyCourse := createEmptyTestCourse(t, courseRepo)
	nextOrder, err := moduleRepo.GetNextModuleOrder(emptyCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get next module order for empty course: %v", err)
	}
	if nextOrder != 1 {
		t.Errorf("Expected next order to be 1 for empty course, got %d", nextOrder)
	}

	// Test with course containing modules
	courseWithModules := createTestCourseWithModules(t, courseRepo)
	nextOrder, err = moduleRepo.GetNextModuleOrder(courseWithModules.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get next module order for course with modules: %v", err)
	}
	if nextOrder != 3 {
		t.Errorf("Expected next order to be 3 for course with 2 modules, got %d", nextOrder)
	}

	// Test with invalid course ID
	_, err = moduleRepo.GetNextModuleOrder("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestCreateModule(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course
	course := createEmptyTestCourse(t, courseRepo)

	// Test creating a module
	newModule := model.Module{
		Title:       "New Module",
		Description: "New module description",
		Order:       1,
		Content:     "New module content",
		CourseID:    course.ID.Hex(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdModule, err := moduleRepo.CreateModule(course.ID.Hex(), newModule)
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	if createdModule.Title != newModule.Title {
		t.Errorf("Expected module title %s, got %s", newModule.Title, createdModule.Title)
	}

	if createdModule.Description != newModule.Description {
		t.Errorf("Expected module description %s, got %s", newModule.Description, createdModule.Description)
	}

	if createdModule.ID.IsZero() {
		t.Error("Expected module to have an ID after creation")
	}

	// Test creating module with invalid course ID
	_, err = moduleRepo.CreateModule("invalid-id", newModule)
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestGetModuleById(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)
	moduleID := course.Modules[0].ID.Hex()

	// Test getting module by valid ID
	foundModule, err := moduleRepo.GetModuleById(moduleID)
	if err != nil {
		t.Fatalf("Failed to get module by ID: %v", err)
	}

	if foundModule.Title != course.Modules[0].Title {
		t.Errorf("Expected module title %s, got %s", course.Modules[0].Title, foundModule.Title)
	}

	// Test getting module with invalid ID
	_, err = moduleRepo.GetModuleById("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid module ID, got nil")
	}

	// Test getting non-existent module
	nonExistentID := primitive.NewObjectID().Hex()
	_, err = moduleRepo.GetModuleById(nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent module, got nil")
	}
}

func TestGetModuleByName(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)

	// Test getting module by valid name
	foundModule, err := moduleRepo.GetModuleByName(course.ID.Hex(), "Module 1")
	if err != nil {
		t.Fatalf("Failed to get module by name: %v", err)
	}

	if foundModule.Title != "Module 1" {
		t.Errorf("Expected module title 'Module 1', got %s", foundModule.Title)
	}

	// Test getting module with non-existent name
	_, err = moduleRepo.GetModuleByName(course.ID.Hex(), "Non-existent Module")
	if err == nil {
		t.Error("Expected error for non-existent module name, got nil")
	}

	// Test with invalid course ID
	_, err = moduleRepo.GetModuleByName("invalid-id", "Module 1")
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestGetModulesByCourseId(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Test with course containing modules
	course := createTestCourseWithModules(t, courseRepo)
	modules, err := moduleRepo.GetModulesByCourseId(course.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules by course ID: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}

	// Verify module titles
	expectedTitles := map[string]bool{"Module 1": true, "Module 2": true}
	for _, module := range modules {
		if !expectedTitles[module.Title] {
			t.Errorf("Unexpected module title: %s", module.Title)
		}
	}

	// Test with empty course
	emptyCourse := createEmptyTestCourse(t, courseRepo)
	modules, err = moduleRepo.GetModulesByCourseId(emptyCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules for empty course: %v", err)
	}

	if len(modules) != 0 {
		t.Errorf("Expected 0 modules for empty course, got %d", len(modules))
	}

	// Test with invalid course ID
	_, err = moduleRepo.GetModulesByCourseId("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestGetModuleByOrder(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)

	// Test getting module by valid order
	foundModule, err := moduleRepo.GetModuleByOrder(course.ID.Hex(), 1)
	if err != nil {
		t.Fatalf("Failed to get module by order: %v", err)
	}

	if foundModule.Order != 1 {
		t.Errorf("Expected module order 1, got %d", foundModule.Order)
	}

	if foundModule.Title != "Module 1" {
		t.Errorf("Expected module title 'Module 1', got %s", foundModule.Title)
	}

	// Test getting module by order 2
	foundModule, err = moduleRepo.GetModuleByOrder(course.ID.Hex(), 2)
	if err != nil {
		t.Fatalf("Failed to get module by order 2: %v", err)
	}

	if foundModule.Order != 2 {
		t.Errorf("Expected module order 2, got %d", foundModule.Order)
	}

	// Test getting module with non-existent order
	_, err = moduleRepo.GetModuleByOrder(course.ID.Hex(), 999)
	if err == nil {
		t.Error("Expected error for non-existent module order, got nil")
	}

	// Test with invalid course ID
	_, err = moduleRepo.GetModuleByOrder("invalid-id", 1)
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestDeleteModule(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)
	moduleToDelete := course.Modules[0]

	// Verify module exists before deletion
	_, err := moduleRepo.GetModuleById(moduleToDelete.ID.Hex())
	if err != nil {
		t.Fatalf("Module should exist before deletion: %v", err)
	}

	// Test deleting the module
	err = moduleRepo.DeleteModule(moduleToDelete.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete module: %v", err)
	}

	// Verify module no longer exists
	_, err = moduleRepo.GetModuleById(moduleToDelete.ID.Hex())
	if err == nil {
		t.Error("Expected error when getting deleted module, got nil")
	}

	// Verify the other module still exists
	remainingModules, err := moduleRepo.GetModulesByCourseId(course.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get remaining modules: %v", err)
	}

	if len(remainingModules) != 1 {
		t.Errorf("Expected 1 remaining module, got %d", len(remainingModules))
	}

	if remainingModules[0].Title != "Module 2" {
		t.Errorf("Expected remaining module to be 'Module 2', got %s", remainingModules[0].Title)
	}

	// Test deleting with invalid ID
	err = moduleRepo.DeleteModule("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid module ID, got nil")
	}

	// Test deleting non-existent module
	nonExistentID := primitive.NewObjectID().Hex()
	err = moduleRepo.DeleteModule(nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent module, got nil")
	}
}

func TestUpdateModuleReorderingFunctionality(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	// Setup
	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course
	course := model.Course{
		Title:          "Test Course",
		Description:    "Test Description",
		TeacherUUID:    "teacher-123",
		TeacherName:    "Test Teacher",
		Capacity:       30,
		StudentsAmount: 0,
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(24 * time.Hour * 30),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Modules: []model.Module{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 1",
				Description: "First module",
				Order:       1,
				Content:     "Content 1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Content:     "Content 2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 3",
				Description: "Third module",
				Order:       3,
				Content:     "Content 3",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 4",
				Description: "Fourth module",
				Order:       4,
				Content:     "Content 4",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 5",
				Description: "Fifth module",
				Order:       5,
				Content:     "Content 5",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	// Insert the course
	createdCourse, err := courseRepo.CreateCourse(course)
	if err != nil {
		t.Fatalf("Failed to create test course: %v", err)
	}

	// Test case: Move module 5 (order 5) to position 2
	// Expected result:
	// - Module 5 should have order 2
	// - Original modules 2, 3, 4 should shift down to orders 3, 4, 5
	// - Module 1 should remain at order 1

	moduleToUpdate := createdCourse.Modules[4] // Module 5 (index 4)
	moduleToUpdate.Order = 2
	moduleToUpdate.Title = "Updated Module 5"

	// Update the module
	updatedModule, err := moduleRepo.UpdateModule(moduleToUpdate.ID.Hex(), moduleToUpdate)
	if err != nil {
		t.Fatalf("Failed to update module: %v", err)
	}

	// Verify the updated module has correct order
	if updatedModule.Order != 2 {
		t.Errorf("Expected updated module to have order 2, got %d", updatedModule.Order)
	}

	if updatedModule.Title != "Updated Module 5" {
		t.Errorf("Expected updated module title to be 'Updated Module 5', got %s", updatedModule.Title)
	}

	// Get all modules to verify the reordering
	allModules, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules: %v", err)
	}

	// Create a map of module titles to their orders for easier verification
	moduleOrders := make(map[string]int)
	for _, module := range allModules {
		moduleOrders[module.Title] = module.Order
	}

	// Verify expected orders
	expectedOrders := map[string]int{
		"Module 1":         1, // Should remain unchanged
		"Updated Module 5": 2, // Moved from 5 to 2
		"Module 2":         3, // Shifted from 2 to 3
		"Module 3":         4, // Shifted from 3 to 4
		"Module 4":         5, // Shifted from 4 to 5
	}

	for title, expectedOrder := range expectedOrders {
		if actualOrder, exists := moduleOrders[title]; !exists {
			t.Errorf("Module '%s' not found", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}

	// Test case 2: Move a module up (from position 5 to position 1)
	// This should test the other direction of reordering
	// Since modules are now sorted by order, we can find Module 4 (which should be at order 5)
	var moduleToMoveUp model.Module
	for _, module := range allModules {
		if module.Title == "Module 4" {
			moduleToMoveUp = module
			break
		}
	}
	moduleToMoveUp.Order = 1

	_, err = moduleRepo.UpdateModule(moduleToMoveUp.ID.Hex(), moduleToMoveUp)
	if err != nil {
		t.Fatalf("Failed to update module for second test: %v", err)
	}

	// Get updated modules
	allModulesAfterSecondUpdate, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules after second update: %v", err)
	}

	// Verify that Module 4 is now at position 1 and others shifted accordingly
	moduleOrdersAfterSecond := make(map[string]int)
	for _, module := range allModulesAfterSecondUpdate {
		moduleOrdersAfterSecond[module.Title] = module.Order
	}

	expectedOrdersAfterSecond := map[string]int{
		"Module 4":         1, // Moved from 5 to 1
		"Module 1":         2, // Shifted from 1 to 2
		"Updated Module 5": 3, // Shifted from 2 to 3
		"Module 2":         4, // Shifted from 3 to 4
		"Module 3":         5, // Shifted from 4 to 5
	}

	for title, expectedOrder := range expectedOrdersAfterSecond {
		if actualOrder, exists := moduleOrdersAfterSecond[title]; !exists {
			t.Errorf("Module '%s' not found in second test", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("Second test - Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}

	// Test case 3: Update module without changing order
	// This should not affect other modules' orders
	// Since array is sorted, Module 4 with order 1 should be at index 0
	moduleToUpdateWithoutOrderChange := allModulesAfterSecondUpdate[0]
	if moduleToUpdateWithoutOrderChange.Title != "Module 4" {
		t.Fatalf("Expected first module to be 'Module 4', got '%s'", moduleToUpdateWithoutOrderChange.Title)
	}

	originalOrder := moduleToUpdateWithoutOrderChange.Order
	moduleToUpdateWithoutOrderChange.Title = "Updated Module 4 Title"
	moduleToUpdateWithoutOrderChange.Description = "Updated description"
	// Don't change the order

	_, err = moduleRepo.UpdateModule(moduleToUpdateWithoutOrderChange.ID.Hex(), moduleToUpdateWithoutOrderChange)
	if err != nil {
		t.Fatalf("Failed to update module without order change: %v", err)
	}

	// Verify that only the content changed, not the orders
	finalModules, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get final modules: %v", err)
	}

	finalModuleOrders := make(map[string]int)
	for _, module := range finalModules {
		finalModuleOrders[module.Title] = module.Order
	}

	// Orders should be the same as in expectedOrdersAfterSecond, except title changed
	expectedFinalOrders := map[string]int{
		"Updated Module 4 Title": 1, // Title changed but order same
		"Module 1":               2,
		"Updated Module 5":       3,
		"Module 2":               4,
		"Module 3":               5,
	}

	for title, expectedOrder := range expectedFinalOrders {
		if actualOrder, exists := finalModuleOrders[title]; !exists {
			t.Errorf("Module '%s' not found in final test", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("Final test - Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}

	// Verify that the updated module has the correct title and description
	var updatedTitleModule *model.Module
	for _, module := range finalModules {
		if module.Title == "Updated Module 4 Title" {
			updatedTitleModule = &module
			break
		}
	}

	if updatedTitleModule == nil {
		t.Error("Could not find module with updated title")
	} else {
		if updatedTitleModule.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got '%s'", updatedTitleModule.Description)
		}
		if updatedTitleModule.Order != originalOrder {
			t.Errorf("Expected order to remain %d, got %d", originalOrder, updatedTitleModule.Order)
		}
	}
}
