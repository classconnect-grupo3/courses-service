package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"

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

func (r *ModuleRepository) UpdateModule(id string, module model.Module) (*model.Module, error) {
	moduleUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid module ID: %v", err)
	}

	filter := bson.M{"modules._id": moduleUUID}
	update := bson.M{"$set": module}

	var course model.Course
	err = r.moduleCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to update module: %v", err)
	}

	return &course.Modules[0], nil
}

func (r *ModuleRepository) DeleteModule(id string) error {
	moduleUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid module ID: %v", err)
	}

	filter := bson.M{"modules._id": moduleUUID}
	update := bson.M{"$pull": bson.M{"modules": bson.M{"_id": moduleUUID}}}

	var course model.Course
	err = r.moduleCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&course)
	if err != nil {
		return fmt.Errorf("failed to delete module: %v", err)
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

	return &course.Modules[0], nil
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
