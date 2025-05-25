package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrollmentRepository struct {
	db                   *mongo.Client
	dbName               string
	enrollmentCollection *mongo.Collection
	courseRepository     *CourseRepository
}

func NewEnrollmentRepository(db *mongo.Client, dbName string, courseRepository *CourseRepository) *EnrollmentRepository {
	return &EnrollmentRepository{db: db, dbName: dbName, enrollmentCollection: db.Database(dbName).Collection("enrollments"), courseRepository: courseRepository}
}

func (r *EnrollmentRepository) createEnrollmentAndModifyCourseCapacity(enrollment model.Enrollment, course *model.Course, ctx context.Context) (interface{}, error) {
	fmt.Printf("enrollment: %v", enrollment)
	res, err := r.enrollmentCollection.InsertOne(ctx, enrollment)
	if err != nil {
		slog.Error("Error creating enrollment", "error", err)
		return nil, err
	}

	enrollment.ID = res.InsertedID.(primitive.ObjectID)

	courseToUpdate := model.Course{
		ID:       course.ID,
		Capacity: course.Capacity - 1,
	}
	_, err = r.courseRepository.UpdateCourse(course.ID.Hex(), courseToUpdate)
	if err != nil {
		slog.Error("Error updating course capacity", "error", err)
		return nil, err
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment model.Enrollment, course *model.Course) error {
	_, err := r.createEnrollmentAndModifyCourseCapacity(enrollment, course, context.TODO())
	if err != nil {
		return err
	}
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

func (r *EnrollmentRepository) deleteEnrollmentAndModifyCourseCapacity(studentID string, course *model.Course, ctx context.Context) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  course.ID.Hex(),
	}

	_, err := r.enrollmentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	courseToUpdate := model.Course{
		ID:       course.ID,
		Capacity: course.Capacity + 1,
	}
	_, err = r.courseRepository.UpdateCourse(course.ID.Hex(), courseToUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (r *EnrollmentRepository) DeleteEnrollment(studentID string, course *model.Course) error {
	err := r.deleteEnrollmentAndModifyCourseCapacity(studentID, course, context.TODO())
	if err != nil {
		return err
	}
	return nil
}
