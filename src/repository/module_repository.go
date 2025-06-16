package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModuleRepository struct {
	db               *mongo.Client
	dbName           string
	moduleCollection *mongo.Collection
}

func NewModuleRepository(db *mongo.Client, dbName string) *ModuleRepository {
	return &ModuleRepository{db: db, dbName: dbName, moduleCollection: db.Database(dbName).Collection("courses")}
}

func (r *ModuleRepository) GetNextModuleOrder(courseID string) (int, error) {
	var course model.Course
	courseUUID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return 0, fmt.Errorf("invalid course ID: %v", err)
	}
	filter := bson.M{"_id": courseUUID}
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return 0, fmt.Errorf("failed to find course: %v", err)
	}

	maxOrder := 0
	for _, module := range course.Modules {
		if module.Order > maxOrder {
			maxOrder = module.Order
		}
	}

	return maxOrder + 1, nil
}

func (r *ModuleRepository) CreateModule(courseID string, module model.Module) (*model.Module, error) {
	module.ID = primitive.NewObjectID()

	courseUUID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %v", err)
	}

	filter := bson.M{"_id": courseUUID}
	update := bson.M{"$push": bson.M{"modules": module}}

	var course model.Course
	err = r.moduleCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %v", err)
	}

	return &module, nil
}

// reorderModules adjusts the order of modules when a module's position changes
func (r *ModuleRepository) reorderModules(modules []model.Module, targetModuleID primitive.ObjectID, newOrder int, oldOrder int) {
	for i := range modules {
		if modules[i].ID == targetModuleID {
			// This is the target module, set its new order
			modules[i].Order = newOrder
		} else {
			// Adjust order for other modules
			currentOrder := modules[i].Order

			if oldOrder < newOrder {
				// Module moved down: shift modules between oldOrder+1 and newOrder up
				if currentOrder > oldOrder && currentOrder <= newOrder {
					modules[i].Order = currentOrder - 1
				}
			} else if oldOrder > newOrder {
				// Module moved up: shift modules between newOrder and oldOrder-1 down
				if currentOrder >= newOrder && currentOrder < oldOrder {
					modules[i].Order = currentOrder + 1
				}
			}
		}
	}
}

// updateModuleFields updates the non-order fields of a target module
func (r *ModuleRepository) updateModuleFields(modules []model.Module, targetModuleID primitive.ObjectID, updatedModule model.Module) {
	for i := range modules {
		if modules[i].ID == targetModuleID {
			if updatedModule.Title != "" {
				modules[i].Title = updatedModule.Title
			}
			if updatedModule.Description != "" {
				modules[i].Description = updatedModule.Description
			}
			// Update Data field - explicit handling for slice
			if updatedModule.Data != nil {
				modules[i].Data = updatedModule.Data
			}
			break
		}
	}
}

func (r *ModuleRepository) UpdateModule(id string, module model.Module) (*model.Module, error) {
	moduleUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid module ID: %v", err)
	}

	// First, get the current module and course to check if order has changed
	filter := bson.M{"modules._id": moduleUUID}
	var currentCourse model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&currentCourse)
	if err != nil {
		return nil, fmt.Errorf("failed to find module: %v", err)
	}

	// Find the current module
	var currentModule *model.Module
	for _, mod := range currentCourse.Modules {
		if mod.ID == moduleUUID {
			currentModule = &mod
			break
		}
	}

	if currentModule == nil {
		return nil, fmt.Errorf("module not found")
	}

	// Check if order has changed
	if module.Order != 0 && module.Order != currentModule.Order {
		// Update module fields first
		r.updateModuleFields(currentCourse.Modules, moduleUUID, module)

		// Reorder modules
		r.reorderModules(currentCourse.Modules, moduleUUID, module.Order, currentModule.Order)

		// Sort modules by order to maintain consistency
		sort.Slice(currentCourse.Modules, func(i, j int) bool {
			return currentCourse.Modules[i].Order < currentCourse.Modules[j].Order
		})

		// Update the entire course with reordered modules
		courseFilter := bson.M{"_id": currentCourse.ID}
		update := bson.M{"$set": bson.M{"modules": currentCourse.Modules}}

		_, err = r.moduleCollection.UpdateOne(context.TODO(), courseFilter, update)
		if err != nil {
			return nil, fmt.Errorf("failed to update module order: %v", err)
		}

		// Return the updated module
		for _, mod := range currentCourse.Modules {
			if mod.ID == moduleUUID {
				return &mod, nil
			}
		}
	} else {
		// No order change, just update the specific module fields
		updateFields := bson.M{}
		if module.Title != "" {
			updateFields["modules.$.title"] = module.Title
		}
		if module.Description != "" {
			updateFields["modules.$.description"] = module.Description
		}
		// Handle Data field update
		if module.Data != nil {
			updateFields["modules.$.data"] = module.Data
		}

		if len(updateFields) > 0 {
			update := bson.M{"$set": updateFields}
			_, err = r.moduleCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return nil, fmt.Errorf("failed to update module: %v", err)
			}
		}

		// Get the updated course to return the module
		err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&currentCourse)
		if err != nil {
			return nil, fmt.Errorf("failed to find updated module: %v", err)
		}

		// Return the updated module
		for _, mod := range currentCourse.Modules {
			if mod.ID == moduleUUID {
				return &mod, nil
			}
		}
	}

	return nil, fmt.Errorf("module not found after update")
}

func (r *ModuleRepository) DeleteModule(id string) error {
	moduleUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid module ID: %v", err)
	}

	// First, get the course and find the module to delete
	filter := bson.M{"modules._id": moduleUUID}
	var course model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return fmt.Errorf("failed to find module: %v", err)
	}

	// Find the module to delete and its order
	var deletedModuleOrder int
	moduleFound := false
	for _, module := range course.Modules {
		if module.ID == moduleUUID {
			deletedModuleOrder = module.Order
			moduleFound = true
			break
		}
	}

	if !moduleFound {
		return fmt.Errorf("module not found")
	}

	// Remove the module from the array
	filteredModules := make([]model.Module, 0)
	for _, module := range course.Modules {
		if module.ID != moduleUUID {
			filteredModules = append(filteredModules, module)
		}
	}

	// Reorder modules: any module with order > deletedModuleOrder should decrease by 1
	for i := range filteredModules {
		if filteredModules[i].Order > deletedModuleOrder {
			filteredModules[i].Order = filteredModules[i].Order - 1
		}
	}

	// Sort modules by order to maintain consistency
	sort.Slice(filteredModules, func(i, j int) bool {
		return filteredModules[i].Order < filteredModules[j].Order
	})

	// Update the course with the reordered modules
	courseFilter := bson.M{"_id": course.ID}
	update := bson.M{"$set": bson.M{"modules": filteredModules}}

	_, err = r.moduleCollection.UpdateOne(context.TODO(), courseFilter, update)
	if err != nil {
		return fmt.Errorf("failed to update modules after deletion: %v", err)
	}

	return nil
}

func (r *ModuleRepository) GetModuleByName(courseID string, moduleName string) (*model.Module, error) {
	courseUUID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %v", err)
	}

	filter := bson.M{"_id": courseUUID, "modules.title": moduleName}

	var course model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to find course or module: %v", err)
	}

	// Find the module with the specified name
	for _, module := range course.Modules {
		if module.Title == moduleName {
			return &module, nil
		}
	}

	return nil, fmt.Errorf("module with name %s not found in course %s", moduleName, courseID)
}

func (r *ModuleRepository) GetModuleById(id string) (*model.Module, error) {
	moduleUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid module ID: %v", err)
	}

	filter := bson.M{"modules._id": moduleUUID}

	var course model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to find course or module: %v", err)
	}

	// Find the module with the specified ID
	for _, module := range course.Modules {
		if module.ID == moduleUUID {
			return &module, nil
		}
	}

	return nil, fmt.Errorf("module with ID %s not found", id)
}

func (r *ModuleRepository) GetModulesByCourseId(courseID string) ([]model.Module, error) {
	courseUUID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %v", err)
	}

	filter := bson.M{"_id": courseUUID}

	var course model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to find course: %v", err)
	}

	// Sort modules by order before returning
	sort.Slice(course.Modules, func(i, j int) bool {
		return course.Modules[i].Order < course.Modules[j].Order
	})

	return course.Modules, nil
}

func (r *ModuleRepository) GetModuleByOrder(courseID string, order int) (*model.Module, error) {
	courseUUID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %v", err)
	}

	filter := bson.M{"_id": courseUUID, "modules.order": order}

	var course model.Course
	err = r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to find course: %v", err)
	}

	for _, module := range course.Modules {
		if module.Order == order {
			return &module, nil
		}
	}

	return nil, fmt.Errorf("module with order %d not found in course %s", order, courseID)
}
