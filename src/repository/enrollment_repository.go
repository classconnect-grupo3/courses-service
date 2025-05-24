package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"

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
		return nil, err
	}

	enrollment.ID = res.InsertedID.(primitive.ObjectID)

	course.Capacity--
	_, err = r.courseRepository.UpdateCourse(course.ID.Hex(), *course)
	if err != nil {
		return nil, err
	}

	fmt.Printf("enrollment: %v", enrollment)
	return enrollment, nil
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment model.Enrollment, course *model.Course) error {
	session, err := r.db.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	if _, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		return r.createEnrollmentAndModifyCourseCapacity(enrollment, course, ctx)
	}); err != nil {
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

	course.Capacity++
	_, err = r.courseRepository.UpdateCourse(course.ID.Hex(), *course)
	if err != nil {
		return err
	}
	return nil
}

func (r *EnrollmentRepository) DeleteEnrollment(studentID string, course *model.Course) error {
	session, err := r.db.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	if _, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, r.deleteEnrollmentAndModifyCourseCapacity(studentID, course, ctx)
	}); err != nil {
		return err
	}

	return nil
}
