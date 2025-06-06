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
	res, err := r.enrollmentCollection.InsertOne(ctx, enrollment)
	if err != nil {
		slog.Error("Error creating enrollment", "error", err)
		return nil, err
	}

	enrollment.ID = res.InsertedID.(primitive.ObjectID)

	err = r.courseRepository.UpdateStudentsAmount(course.ID.Hex(), course.StudentsAmount+1)
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

func (r *EnrollmentRepository) GetEnrollmentsByCourseId(courseID string) ([]*model.Enrollment, error) {
	filter := bson.M{
		"course_id": courseID,
	}

	cursor, err := r.enrollmentCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var enrollments []*model.Enrollment
	if err := cursor.All(context.TODO(), &enrollments); err != nil {
		return nil, err
	}

	// Ensure we always return a non-nil slice
	if enrollments == nil {
		enrollments = []*model.Enrollment{}
	}

	return enrollments, nil
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

	result, err := r.enrollmentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Only update course capacity if we actually deleted an enrollment
	if result.DeletedCount > 0 {
		err = r.courseRepository.UpdateStudentsAmount(course.ID.Hex(), course.StudentsAmount-1)
		if err != nil {
			return err
		}
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

func (r *EnrollmentRepository) SetFavouriteCourse(studentID, courseID string) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
	}

	update := bson.M{
		"$set": bson.M{
			"favourite": true,
		},
	}

	res, err := r.enrollmentCollection.UpdateOne(context.TODO(), filter, update)
	if res.MatchedCount == 0 {
		return fmt.Errorf("enrollment not found for student %s in course %s", studentID, courseID)
	}
	if err != nil {
		return fmt.Errorf("error setting favourite course for student %s in course %s", studentID, courseID)
	}
	return nil
}

func (r *EnrollmentRepository) UnsetFavouriteCourse(studentID, courseID string) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
	}

	update := bson.M{
		"$set": bson.M{
			"favourite": false,
		},
	}

	res, err := r.enrollmentCollection.UpdateOne(context.TODO(), filter, update)
	if res.MatchedCount == 0 {
		return fmt.Errorf("enrollment not found for student %s in course %s", studentID, courseID)
	}
	if err != nil {
		return fmt.Errorf("error unsetting favourite course for student %s in course %s", studentID, courseID)
	}
	return nil
}

func (r *EnrollmentRepository) GetEnrollmentsByStudentId(studentID string) ([]*model.Enrollment, error) {
	filter := bson.M{
		"student_id": studentID,
	}

	cursor, err := r.enrollmentCollection.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []*model.Enrollment{}, nil
		}
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var enrollments []*model.Enrollment
	if err := cursor.All(context.TODO(), &enrollments); err != nil {
		return nil, err
	}

	return enrollments, nil
}