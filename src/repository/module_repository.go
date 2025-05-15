package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
	filter := bson.M{"_id": courseID}
	err := r.moduleCollection.FindOne(context.TODO(), filter).Decode(&course)
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
	filter := bson.M{"_id": courseID}
	update := bson.M{"$push": bson.M{"modules": module}}

	var course model.Course
	err := r.moduleCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %v", err)
	}

	return &module, nil
}


