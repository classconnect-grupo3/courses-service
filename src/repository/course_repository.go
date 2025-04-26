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

func (r *CourseRepository) CreateCourse(course model.Course) (*model.Course, error) {
	collection := r.db.Database(r.dbName).Collection("courses")

	_, err := collection.InsertOne(context.TODO(), course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) GetCourseById(id string) (*model.Course, error) {
	collection := r.db.Database(r.dbName).Collection("courses")

	var course model.Course
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	collection := r.db.Database(r.dbName).Collection("courses")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
