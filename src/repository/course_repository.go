package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

type CourseRepository struct {
	db               *mongo.Client
	dbName           string
	courseCollection *mongo.Collection
}

func filterEmptyFields(course model.Course) any {
	updates := bson.D{}

	courseType := reflect.TypeOf(course)
	courseValue := reflect.ValueOf(course)

	for i := 0; i < courseType.NumField(); i++ {
		field := courseType.Field(i)
		fieldValue := courseValue.Field(i)
		tag := field.Tag.Get("json")
		if !isZeroType(fieldValue) {
			update := bson.E{Key: tag, Value: fieldValue.Interface()}
			updates = append(updates, update)
		}
	}

	return updates
}

func isZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	case reflect.String:
		return value.String() == ""
	case reflect.Int:
		return value.Int() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Float64:
		return value.Float() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
	}
}
func NewCourseRepository(db *mongo.Client, dbName string) *CourseRepository {
	return &CourseRepository{db: db, dbName: dbName, courseCollection: db.Database(dbName).Collection("courses")}
}

func (r *CourseRepository) CreateCourse(course model.Course) (*model.Course, error) {
	result, err := r.courseCollection.InsertOne(context.TODO(), course)
	if err != nil {
		return nil, fmt.Errorf("failed to create course: %v", err)
	}

	course.ID = result.InsertedID.(primitive.ObjectID)
	return &course, nil
}

func (r *CourseRepository) GetCourses() ([]*model.Course, error) {
	cursor, err := r.courseCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	var courses []*model.Course
	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	return courses, nil
}

func (r *CourseRepository) GetCourseById(id string) (*model.Course, error) {
	var course model.Course
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by id: %v", err)
	}
	err = r.courseCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by id: %v", err)
	}
	return &course, nil
}

func (r *CourseRepository) GetCourseByTeacherId(teacherId string) ([]*model.Course, error) {
	cursor, err := r.courseCollection.Find(context.TODO(), bson.M{"teacher_uuid": teacherId})
	if err != nil {
		return nil, fmt.Errorf("failed to get course by teacher id: %v", err)
	}

	var courses []*model.Course
	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to get course by teacher id: %v", err)
	}
	return courses, nil
}

func (r *CourseRepository) GetCourseByTitle(title string) ([]*model.Course, error) {
	filter := bson.M{
		"title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}

	var courses []*model.Course
	cursor, err := r.courseCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by title: %v", err)
	}

	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to get course by title: %v", err)
	}

	return courses, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %v", err)
	}
	_, err = r.courseCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		return fmt.Errorf("failed to delete course: %v", err)
	}
	return nil
}

func (r *CourseRepository) UpdateCourse(id string, updateCourseRequest model.Course) (*model.Course, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %v", err)
	}

	update := filterEmptyFields(updateCourseRequest)

	_, err = r.courseCollection.UpdateOne(context.TODO(), bson.M{"_id": objectId}, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %v", err)
	}

	updatedCourse, err := r.GetCourseById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %v", err)
	}

	return updatedCourse, nil
}
