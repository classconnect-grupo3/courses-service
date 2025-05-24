package repository

import (
	"context"
	"courses-service/src/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentRepository struct {
	db                   *mongo.Client
	dbName               string
	enrollmentCollection *mongo.Collection
}

func NewEnrollmentRepository(db *mongo.Client, dbName string) *EnrollmentRepository {
	return &EnrollmentRepository{db: db, dbName: dbName, enrollmentCollection: db.Database(dbName).Collection("enrollments")}
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment model.Enrollment) error {
	res, err := r.enrollmentCollection.InsertOne(context.TODO(), enrollment)
	if err != nil {
		return err
	}
	enrollment.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *EnrollmentRepository) IsEnrolled(studentID, courseID string) (bool, error) {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
	}

	var enrollment model.Enrollment
	err := r.enrollmentCollection.FindOne(context.TODO(), filter).Decode(&enrollment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
