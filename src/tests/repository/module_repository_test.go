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
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Data: []model.ModuleData{
					{
						Id:          "data-1",
						ModuleId:    "module-2-id",
						Title:       "Test Data",
						Description: "Test data description",
						Resources: []model.ModuleDataResource{
							{
								Id:   1,
								Name: "Test Resource",
								Url:  "https://example.com/test",
							},
						},
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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

	// Test creating a module with empty data
	newModule := model.Module{
		Title:       "New Module",
		Description: "New module description",
		Order:       1,
		Data:        []model.ModuleData{},
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

	// Verify Data field is properly handled
	if createdModule.Data == nil {
		t.Error("Expected Data field to be initialized")
	}

	if len(createdModule.Data) != 0 {
		t.Errorf("Expected empty Data array, got %d elements", len(createdModule.Data))
	}

	// Test creating module with invalid course ID
	_, err = moduleRepo.CreateModule("invalid-id", newModule)
	if err == nil {
		t.Error("Expected error for invalid course ID, got nil")
	}
}

func TestCreateModuleWithData(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course
	course := createEmptyTestCourse(t, courseRepo)

	// Test creating a module with data
	testData := []model.ModuleData{
		{
			Id:          "test-data-1",
			ModuleId:    "will-be-set-after-creation",
			Title:       "Test Data Item",
			Description: "Test data description",
			Resources: []model.ModuleDataResource{
				{
					Id:   1,
					Name: "Test Resource",
					Url:  "https://example.com/resource",
				},
			},
		},
	}

	newModule := model.Module{
		Title:       "Module With Data",
		Description: "Module with test data",
		Order:       1,
		Data:        testData,
		CourseID:    course.ID.Hex(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdModule, err := moduleRepo.CreateModule(course.ID.Hex(), newModule)
	if err != nil {
		t.Fatalf("Failed to create module with data: %v", err)
	}

	// Verify data was properly stored
	if len(createdModule.Data) != 1 {
		t.Errorf("Expected 1 data item, got %d", len(createdModule.Data))
	}

	if createdModule.Data[0].Title != "Test Data Item" {
		t.Errorf("Expected data title 'Test Data Item', got %s", createdModule.Data[0].Title)
	}

	if len(createdModule.Data[0].Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(createdModule.Data[0].Resources))
	}

	if createdModule.Data[0].Resources[0].Name != "Test Resource" {
		t.Errorf("Expected resource name 'Test Resource', got %s", createdModule.Data[0].Resources[0].Name)
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

	// Verify Data field is properly loaded
	if foundModule.Data == nil {
		t.Error("Expected Data field to be loaded")
	}

	// Test getting module with data
	moduleWithDataID := course.Modules[1].ID.Hex()
	foundModuleWithData, err := moduleRepo.GetModuleById(moduleWithDataID)
	if err != nil {
		t.Fatalf("Failed to get module with data by ID: %v", err)
	}

	if len(foundModuleWithData.Data) != 1 {
		t.Errorf("Expected 1 data item, got %d", len(foundModuleWithData.Data))
	}

	if foundModuleWithData.Data[0].Title != "Test Data" {
		t.Errorf("Expected data title 'Test Data', got %s", foundModuleWithData.Data[0].Title)
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

	// Verify Data field is properly loaded
	if foundModule.Data == nil {
		t.Error("Expected Data field to be loaded")
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
		// Verify Data field is properly loaded for all modules
		if module.Data == nil {
			t.Errorf("Expected Data field to be loaded for module %s", module.Title)
		}
	}

	// Verify first module has empty data, second has data
	for _, module := range modules {
		if module.Title == "Module 1" && len(module.Data) != 0 {
			t.Errorf("Expected Module 1 to have empty data, got %d items", len(module.Data))
		}
		if module.Title == "Module 2" && len(module.Data) != 1 {
			t.Errorf("Expected Module 2 to have 1 data item, got %d items", len(module.Data))
		}
	}

	// Test with empty course
	emptyCourse := createEmptyTestCourse(t, courseRepo)
	emptyModules, err := moduleRepo.GetModulesByCourseId(emptyCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules for empty course: %v", err)
	}

	if len(emptyModules) != 0 {
		t.Errorf("Expected 0 modules for empty course, got %d", len(emptyModules))
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
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 3",
				Description: "Third module",
				Order:       3,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 4",
				Description: "Fourth module",
				Order:       4,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 5",
				Description: "Fifth module",
				Order:       5,
				Data:        []model.ModuleData{},
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

func TestUpdateModuleWithEmptyData(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)
	moduleToUpdate := course.Modules[1] // Module 2 has data

	// Update with empty data
	moduleToUpdate.Title = "Updated Module Title"
	moduleToUpdate.Description = "Updated description"
	moduleToUpdate.Data = []model.ModuleData{} // Explicitly set to empty

	updatedModule, err := moduleRepo.UpdateModule(moduleToUpdate.ID.Hex(), moduleToUpdate)
	if err != nil {
		t.Fatalf("Failed to update module with empty data: %v", err)
	}

	if updatedModule.Title != "Updated Module Title" {
		t.Errorf("Expected updated title 'Updated Module Title', got %s", updatedModule.Title)
	}

	// Verify data was cleared
	if len(updatedModule.Data) != 0 {
		t.Errorf("Expected empty data array, got %d items", len(updatedModule.Data))
	}
}

func TestUpdateModuleWithData(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with modules
	course := createTestCourseWithModules(t, courseRepo)
	moduleToUpdate := course.Modules[0] // Module 1 has empty data

	// Update with new data
	newData := []model.ModuleData{
		{
			Id:          "new-data-1",
			ModuleId:    moduleToUpdate.ID.Hex(),
			Title:       "New Data Item",
			Description: "New data description",
			Resources: []model.ModuleDataResource{
				{
					Id:   2,
					Name: "New Resource",
					Url:  "https://example.com/newresource",
				},
			},
		},
		{
			Id:          "new-data-2",
			ModuleId:    moduleToUpdate.ID.Hex(),
			Title:       "Second Data Item",
			Description: "Second data description",
			Resources:   []model.ModuleDataResource{},
		},
	}

	moduleToUpdate.Title = "Updated Module With Data"
	moduleToUpdate.Data = newData

	updatedModule, err := moduleRepo.UpdateModule(moduleToUpdate.ID.Hex(), moduleToUpdate)
	if err != nil {
		t.Fatalf("Failed to update module with data: %v", err)
	}

	if updatedModule.Title != "Updated Module With Data" {
		t.Errorf("Expected updated title 'Updated Module With Data', got %s", updatedModule.Title)
	}

	// Verify data was updated
	if len(updatedModule.Data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(updatedModule.Data))
	}

	if updatedModule.Data[0].Title != "New Data Item" {
		t.Errorf("Expected first data title 'New Data Item', got %s", updatedModule.Data[0].Title)
	}

	if updatedModule.Data[1].Title != "Second Data Item" {
		t.Errorf("Expected second data title 'Second Data Item', got %s", updatedModule.Data[1].Title)
	}

	// Test first item has resources, second doesn't
	if len(updatedModule.Data[0].Resources) != 1 {
		t.Errorf("Expected 1 resource in first data item, got %d", len(updatedModule.Data[0].Resources))
	}

	if len(updatedModule.Data[1].Resources) != 0 {
		t.Errorf("Expected 0 resources in second data item, got %d", len(updatedModule.Data[1].Resources))
	}
}

func TestDeleteModuleWithReordering(t *testing.T) {
	t.Cleanup(func() {
		dbSetup.CleanupCollection("courses")
	})

	moduleRepo := repository.NewModuleRepository(dbSetup.Client, dbSetup.DBName)
	courseRepo := repository.NewCourseRepository(dbSetup.Client, dbSetup.DBName)

	// Create a test course with 5 modules in order 1, 2, 3, 4, 5
	course := model.Course{
		Title:          "Test Course for Reordering",
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
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 2",
				Description: "Second module",
				Order:       2,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 3",
				Description: "Third module",
				Order:       3,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 4",
				Description: "Fourth module",
				Order:       4,
				Data:        []model.ModuleData{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				Title:       "Module 5",
				Description: "Fifth module",
				Order:       5,
				Data:        []model.ModuleData{},
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

	// Delete Module 3 (middle module)
	moduleToDelete := createdCourse.Modules[2] // Module 3 (index 2)
	err = moduleRepo.DeleteModule(moduleToDelete.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete module: %v", err)
	}

	// Get remaining modules
	remainingModules, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get remaining modules: %v", err)
	}

	// Should have 4 modules now
	if len(remainingModules) != 4 {
		t.Errorf("Expected 4 remaining modules, got %d", len(remainingModules))
	}

	// Verify correct reordering: should be 1, 2, 3, 4 (no gaps)
	expectedOrders := map[string]int{
		"Module 1": 1, // Should remain unchanged
		"Module 2": 2, // Should remain unchanged
		"Module 4": 3, // Should shift down from 4 to 3
		"Module 5": 4, // Should shift down from 5 to 4
	}

	moduleOrders := make(map[string]int)
	for _, module := range remainingModules {
		moduleOrders[module.Title] = module.Order
	}

	for title, expectedOrder := range expectedOrders {
		if actualOrder, exists := moduleOrders[title]; !exists {
			t.Errorf("Module '%s' not found after deletion", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}

	// Verify Module 3 was actually deleted
	_, err = moduleRepo.GetModuleById(moduleToDelete.ID.Hex())
	if err == nil {
		t.Error("Expected error when getting deleted module, got nil")
	}

	// Test case 2: Delete first module
	// Delete Module 1
	moduleToDeleteFirst := remainingModules[0] // Should be Module 1
	if moduleToDeleteFirst.Title != "Module 1" {
		t.Fatalf("Expected first module to be 'Module 1', got '%s'", moduleToDeleteFirst.Title)
	}

	err = moduleRepo.DeleteModule(moduleToDeleteFirst.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete first module: %v", err)
	}

	// Get modules after deleting first
	modulesAfterFirstDeletion, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get modules after first deletion: %v", err)
	}

	// Should have 3 modules now
	if len(modulesAfterFirstDeletion) != 3 {
		t.Errorf("Expected 3 modules after first deletion, got %d", len(modulesAfterFirstDeletion))
	}

	// Verify reordering after deleting first module: should be 1, 2, 3
	expectedOrdersAfterFirst := map[string]int{
		"Module 2": 1, // Should shift down from 2 to 1
		"Module 4": 2, // Should shift down from 3 to 2 (was already 3 after previous deletion)
		"Module 5": 3, // Should shift down from 4 to 3 (was already 4 after previous deletion)
	}

	moduleOrdersAfterFirst := make(map[string]int)
	for _, module := range modulesAfterFirstDeletion {
		moduleOrdersAfterFirst[module.Title] = module.Order
	}

	for title, expectedOrder := range expectedOrdersAfterFirst {
		if actualOrder, exists := moduleOrdersAfterFirst[title]; !exists {
			t.Errorf("Module '%s' not found after first deletion", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("After first deletion - Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}

	// Test case 3: Delete last module
	// Delete Module 5 (which should now be at order 3)
	var moduleToDeleteLast model.Module
	for _, module := range modulesAfterFirstDeletion {
		if module.Title == "Module 5" {
			moduleToDeleteLast = module
			break
		}
	}

	if moduleToDeleteLast.Title != "Module 5" {
		t.Fatalf("Could not find Module 5 to delete")
	}

	err = moduleRepo.DeleteModule(moduleToDeleteLast.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete last module: %v", err)
	}

	// Get final modules
	finalModules, err := moduleRepo.GetModulesByCourseId(createdCourse.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get final modules: %v", err)
	}

	// Should have 2 modules now
	if len(finalModules) != 2 {
		t.Errorf("Expected 2 final modules, got %d", len(finalModules))
	}

	// Verify final ordering: should be 1, 2
	expectedFinalOrders := map[string]int{
		"Module 2": 1, // Should remain unchanged
		"Module 4": 2, // Should remain unchanged (was already 2)
	}

	finalModuleOrders := make(map[string]int)
	for _, module := range finalModules {
		finalModuleOrders[module.Title] = module.Order
	}

	for title, expectedOrder := range expectedFinalOrders {
		if actualOrder, exists := finalModuleOrders[title]; !exists {
			t.Errorf("Module '%s' not found in final modules", title)
		} else if actualOrder != expectedOrder {
			t.Errorf("Final modules - Module '%s' expected order %d, got %d", title, expectedOrder, actualOrder)
		}
	}
}
