package repository

import (
	"context"
	"courses-service/src/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepository struct {
	db     *mongo.Client
	dbName string	
}

func NewCourseRepository(db *mongo.Client, dbName string) *CourseRepository {
	return &CourseRepository{db: db, dbName: dbName}
}

func (r *CourseRepository) GetCourses() ([]*model.Course, error) {
	cursor, err := r.db.Database(r.dbName).Collection("courses").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var courses []*model.Course
	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, err
	}
	return courses, nil
}
