package repository

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

type CourseRepository struct {
	db                   *mongo.Client
	dbName               string
	courseCollection     *mongo.Collection
	enrollmentCollection *mongo.Collection
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
	return &CourseRepository{
		db:                   db,
		dbName:               dbName,
		courseCollection:     db.Database(dbName).Collection("courses"),
		enrollmentCollection: db.Database(dbName).Collection("enrollments"),
	}
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

func (r *CourseRepository) GetCoursesByStudentId(studentId string) ([]*model.Course, error) {
	// First, get all enrollment records for this student
	cursor, err := r.enrollmentCollection.Find(context.TODO(), bson.M{"student_id": studentId})
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments by student id: %v", err)
	}

	// Parse enrollments to get course IDs
	var enrollments []model.Enrollment
	if err := cursor.All(context.TODO(), &enrollments); err != nil {
		return nil, fmt.Errorf("failed to parse enrollments: %v", err)
	}

	if len(enrollments) == 0 {
		return []*model.Course{}, nil
	}

	// Extract course IDs from enrollments
	var courseIds []primitive.ObjectID
	for _, enrollment := range enrollments {
		courseId, err := primitive.ObjectIDFromHex(enrollment.CourseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course id in enrollment: %v", err)
		}
		courseIds = append(courseIds, courseId)
	}

	// Find all courses with these IDs
	filter := bson.M{"_id": bson.M{"$in": courseIds}}
	courseCursor, err := r.courseCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses by ids: %v", err)
	}

	// Parse courses
	var courses []*model.Course
	if err := courseCursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to parse courses: %v", err)
	}

	return courses, nil
}

func (r *CourseRepository) GetCoursesByAuxTeacherId(auxTeacherId string) ([]*model.Course, error) {
	// Find all courses where the auxTeacherId is in the aux_teachers array
	// Using $in to be explicit about searching within an array
	filter := bson.M{"aux_teachers": bson.M{"$in": []string{auxTeacherId}}}
	cursor, err := r.courseCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses by aux teacher id: %v", err)
	}

	var courses []*model.Course
	if err := cursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to get courses by aux teacher id: %v", err)
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

func (r *CourseRepository) AddAuxTeacherToCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	course.AuxTeachers = append(course.AuxTeachers, auxTeacherId)
	course.UpdatedAt = time.Now()

	// Direct MongoDB update to ensure we can set the exact AuxTeachers array
	update := bson.M{
		"$set": bson.M{
			"aux_teachers": course.AuxTeachers,
			"updated_at":   course.UpdatedAt,
		},
	}

	_, err := r.courseCollection.UpdateOne(context.TODO(), bson.M{"_id": course.ID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to add aux teacher to course: %v", err)
	}

	return r.GetCourseById(course.ID.Hex())
}

func (r *CourseRepository) RemoveAuxTeacherFromCourse(course *model.Course, auxTeacherId string) (*model.Course, error) {
	// Buscar y eliminar el auxTeacherId del slice
	for i, teacher := range course.AuxTeachers {
		if teacher == auxTeacherId {
			// Eliminar el elemento del slice
			course.AuxTeachers = append(course.AuxTeachers[:i], course.AuxTeachers[i+1:]...)
			break
		}
	}

	course.UpdatedAt = time.Now()

	// Direct MongoDB update to ensure we can set empty arrays
	update := bson.M{
		"$set": bson.M{
			"aux_teachers": course.AuxTeachers,
			"updated_at":   course.UpdatedAt,
		},
	}

	_, err := r.courseCollection.UpdateOne(context.TODO(), bson.M{"_id": course.ID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to remove aux teacher from course: %v", err)
	}

	return r.GetCourseById(course.ID.Hex())
}

func (r *CourseRepository) UpdateStudentsAmount(courseID string, newStudentsAmount int) error {
	objectId, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return fmt.Errorf("failed to update students amount: %v", err)
	}

	// Direct MongoDB update to ensure we can set StudentsAmount to 0
	update := bson.M{
		"$set": bson.M{
			"students_amount": newStudentsAmount,
			"updated_at":      time.Now(),
		},
	}

	_, err = r.courseCollection.UpdateOne(context.TODO(), bson.M{"_id": objectId}, update)
	if err != nil {
		return fmt.Errorf("failed to update students amount: %v", err)
	}

	return nil
}

func (r *CourseRepository) CreateCourseFeedback(courseID string, feedback model.CourseFeedback) (*model.CourseFeedback, error) {
	feedback.ID = primitive.NewObjectID()

	course, err := r.GetCourseById(courseID)
	if err != nil {
		return nil, err
	}

	if course.Feedback == nil {
		course.Feedback = []model.CourseFeedback{}
	}

	course.Feedback = append(course.Feedback, feedback)

	_, err = r.UpdateCourse(course.ID.Hex(), *course)
	if err != nil {
		return nil, err
	}

	return &feedback, nil
}

func (r *CourseRepository) matchesFeedbackFilters(feedback *model.CourseFeedback, request schemas.GetCourseFeedbackRequest) bool {
	// Filter by feedback type
	if request.FeedbackType != "" && feedback.FeedbackType != request.FeedbackType {
		return false
	}

	// Filter by score range
	if request.StartScore != 0 && feedback.Score < request.StartScore {
		return false
	}
	if request.EndScore != 0 && feedback.Score > request.EndScore {
		return false
	}

	// Filter by date range
	if !request.StartDate.IsZero() && feedback.CreatedAt.Before(request.StartDate) {
		return false
	}
	if !request.EndDate.IsZero() && feedback.CreatedAt.After(request.EndDate) {
		return false
	}

	return true
}

func (r *CourseRepository) GetCourseFeedback(courseID string, getCourseFeedbackRequest schemas.GetCourseFeedbackRequest) ([]*model.CourseFeedback, error) {
	// Convert courseID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return nil, errors.New("invalid course ID format: " + err.Error())
	}

	// Get the course document
	var course model.Course
	err = r.courseCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("course not found")
		}
		return nil, errors.New("error getting course: " + err.Error())
	}

	// If course has no feedback, return empty slice
	if len(course.Feedback) == 0 {
		return []*model.CourseFeedback{}, nil
	}

	// Apply filters to the feedback
	filteredFeedbacks := []*model.CourseFeedback{}
	for _, feedback := range course.Feedback {
		if r.matchesFeedbackFilters(&feedback, getCourseFeedbackRequest) {
			// Create a copy to avoid pointer issues
			feedbackCopy := feedback
			filteredFeedbacks = append(filteredFeedbacks, &feedbackCopy)
		}
	}

	return filteredFeedbacks, nil
}

// CountCourses returns the total number of courses
func (r *CourseRepository) CountCourses() (int64, error) {
	count, err := r.courseCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count courses: %v", err)
	}
	return count, nil
}

// CountActiveCourses returns the number of active courses (courses that have not ended yet)
func (r *CourseRepository) CountActiveCourses() (int64, error) {
	now := time.Now()
	filter := bson.M{"end_date": bson.M{"$gt": now}}
	count, err := r.courseCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count active courses: %v", err)
	}
	return count, nil
}

// CountFinishedCourses returns the number of finished courses (courses that have ended)
func (r *CourseRepository) CountFinishedCourses() (int64, error) {
	now := time.Now()
	filter := bson.M{"end_date": bson.M{"$lte": now}}
	count, err := r.courseCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count finished courses: %v", err)
	}
	return count, nil
}

// CountCoursesCreatedThisMonth returns the number of courses created this month
func (r *CourseRepository) CountCoursesCreatedThisMonth() (int64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	count, err := r.courseCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count courses created this month: %v", err)
	}
	return count, nil
}

// CountUniqueTeachers returns the number of unique teachers
func (r *CourseRepository) CountUniqueTeachers() (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$teacher_uuid"}},
		{"$count": "unique_teachers"},
	}

	cursor, err := r.courseCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique teachers: %v", err)
	}
	defer cursor.Close(context.TODO())

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return 0, fmt.Errorf("failed to decode unique teachers count: %v", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	count, ok := result[0]["unique_teachers"].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected result format for unique teachers count")
	}

	return int64(count), nil
}

// CountUniqueAuxTeachers returns the number of unique auxiliary teachers
func (r *CourseRepository) CountUniqueAuxTeachers() (int64, error) {
	pipeline := []bson.M{
		{"$unwind": "$aux_teachers"},
		{"$group": bson.M{"_id": "$aux_teachers"}},
		{"$count": "unique_aux_teachers"},
	}

	cursor, err := r.courseCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique aux teachers: %v", err)
	}
	defer cursor.Close(context.TODO())

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return 0, fmt.Errorf("failed to decode unique aux teachers count: %v", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	count, ok := result[0]["unique_aux_teachers"].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected result format for unique aux teachers count")
	}

	return int64(count), nil
}

// GetTopTeachersByCourseCount returns top teachers by course count
func (r *CourseRepository) GetTopTeachersByCourseCount(limit int) ([]schemas.CourseDistributionByTeacher, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id": bson.M{
				"teacher_id":   "$teacher_uuid",
				"teacher_name": "$teacher_name",
			},
			"course_count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"course_count": -1}},
		{"$limit": limit},
		{"$project": bson.M{
			"teacher_id":   "$_id.teacher_id",
			"teacher_name": "$_id.teacher_name",
			"course_count": 1,
			"_id":          0,
		}},
	}

	cursor, err := r.courseCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get top teachers: %v", err)
	}
	defer cursor.Close(context.TODO())

	var teachers []schemas.CourseDistributionByTeacher
	if err = cursor.All(context.TODO(), &teachers); err != nil {
		return nil, fmt.Errorf("failed to decode top teachers: %v", err)
	}

	return teachers, nil
}

// GetRecentCourses returns recent courses with basic information
func (r *CourseRepository) GetRecentCourses(limit int) ([]schemas.CourseBasicInfo, error) {
	pipeline := []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{"$limit": limit},
		{"$project": bson.M{
			"id":              bson.M{"$toString": "$_id"},
			"title":           1,
			"teacher_name":    1,
			"students_amount": 1,
			"capacity":        1,
			"created_at":      1,
		}},
	}

	cursor, err := r.courseCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent courses: %v", err)
	}
	defer cursor.Close(context.TODO())

	var courses []schemas.CourseBasicInfo
	if err = cursor.All(context.TODO(), &courses); err != nil {
		return nil, fmt.Errorf("failed to decode recent courses: %v", err)
	}

	return courses, nil
}
