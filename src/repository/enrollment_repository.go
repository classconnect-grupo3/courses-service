package repository

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"fmt"
	"log/slog"
	"time"

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

func (r *EnrollmentRepository) GetEnrollmentByStudentIdAndCourseId(studentID, courseID string) (*model.Enrollment, error) {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
	}

	var enrollment model.Enrollment
	err := r.enrollmentCollection.FindOne(context.TODO(), filter).Decode(&enrollment)
	if err != nil {
		return nil, err
	}

	return &enrollment, nil
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

func (r *EnrollmentRepository) CreateStudentFeedback(feedbackRequest model.StudentFeedback, enrollmentID string) error {
	feedbackRequest.ID = primitive.NewObjectID()

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(enrollmentID)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %v", err)
	}

	_, err = r.enrollmentCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$push": bson.M{"feedback": feedbackRequest}})
	if err != nil {
		return err
	}

	return nil
}

func (r *EnrollmentRepository) matchesFeedbackFilters(feedback *model.StudentFeedback, request schemas.GetFeedbackByStudentIdRequest) bool {
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

func (r *EnrollmentRepository) GetFeedbackByStudentId(studentID string, getFeedbackByStudentIdRequest schemas.GetFeedbackByStudentIdRequest) ([]*model.StudentFeedback, error) {
	// Build base filter for enrollments of the student
	filter := bson.M{
		"student_id": studentID,
		"feedback":   bson.M{"$exists": true, "$ne": []interface{}{}}, // Must have non-empty feedback array
	}

	// Add course filter if specified
	if getFeedbackByStudentIdRequest.CourseID != "" {
		filter["course_id"] = getFeedbackByStudentIdRequest.CourseID
	}

	// Find all enrollments for this student
	cursor, err := r.enrollmentCollection.Find(context.TODO(), filter)
	if err != nil {
		return []*model.StudentFeedback{}, nil // Return empty slice on error instead of nil
	}
	defer cursor.Close(context.TODO())

	var enrollments []*model.Enrollment
	if err := cursor.All(context.TODO(), &enrollments); err != nil {
		return []*model.StudentFeedback{}, nil // Return empty slice on error instead of nil
	}

	// Extract and filter feedbacks from enrollments
	var allFeedbacks []*model.StudentFeedback
	for _, enrollment := range enrollments {
		for _, feedback := range enrollment.Feedback {
			// Apply feedback filters
			if r.matchesFeedbackFilters(&feedback, getFeedbackByStudentIdRequest) {
				// Create a copy to avoid pointer issues
				feedbackCopy := feedback
				allFeedbacks = append(allFeedbacks, &feedbackCopy)
			}
		}
	}

	// Ensure we always return a non-nil slice
	if allFeedbacks == nil {
		allFeedbacks = []*model.StudentFeedback{}
	}

	return allFeedbacks, nil
}

// ApproveStudent updates an enrollment status to completed and sets completion date
func (r *EnrollmentRepository) ApproveStudent(studentID, courseID string) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
		"status":     model.EnrollmentStatusActive, // Only approve active enrollments
	}

	update := bson.M{
		"$set": bson.M{
			"status":         model.EnrollmentStatusCompleted,
			"completed_date": time.Now(),
			"updated_at":     time.Now(),
		},
	}

	result, err := r.enrollmentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating enrollment: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("enrollment not found or student is not active in course %s", courseID)
	}

	return nil
}

// DisapproveStudent updates an enrollment status to dropped and sets the reason for unenrollment
func (r *EnrollmentRepository) DisapproveStudent(studentID, courseID, reason string) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
		"status":     model.EnrollmentStatusActive, // Only disapprove active enrollments
	}

	update := bson.M{
		"$set": bson.M{
			"status":                  model.EnrollmentStatusDropped,
			"reason_for_unenrollment": reason,
			"updated_at":              time.Now(),
		},
	}

	result, err := r.enrollmentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating enrollment: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("enrollment not found or student is not active in course %s", courseID)
	}

	return nil
}

// ReactivateDroppedEnrollment reactivates a dropped enrollment and clears the reason
func (r *EnrollmentRepository) ReactivateDroppedEnrollment(studentID, courseID string) error {
	filter := bson.M{
		"student_id": studentID,
		"course_id":  courseID,
		"status":     model.EnrollmentStatusDropped, // Only reactivate dropped enrollments
	}

	update := bson.M{
		"$set": bson.M{
			"status":     model.EnrollmentStatusActive,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{
			"reason_for_unenrollment": "", // Remove the field completely
		},
	}

	result, err := r.enrollmentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error reactivating enrollment: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("dropped enrollment not found for student %s in course %s", studentID, courseID)
	}

	return nil
}

// CountEnrollments returns the total number of enrollments
func (r *EnrollmentRepository) CountEnrollments() (int64, error) {
	count, err := r.enrollmentCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count enrollments: %v", err)
	}
	return count, nil
}

// CountEnrollmentsByStatus returns the number of enrollments by status
func (r *EnrollmentRepository) CountEnrollmentsByStatus(status model.EnrollmentStatus) (int64, error) {
	filter := bson.M{"status": status}
	count, err := r.enrollmentCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count enrollments by status: %v", err)
	}
	return count, nil
}

// CountEnrollmentsThisMonth returns the number of enrollments created this month
func (r *EnrollmentRepository) CountEnrollmentsThisMonth() (int64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := bson.M{
		"enrolled_at": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	count, err := r.enrollmentCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count enrollments this month: %v", err)
	}
	return count, nil
}

// CountUniqueStudents returns the number of unique students
func (r *EnrollmentRepository) CountUniqueStudents() (int64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$student_id"}},
		{"$count": "unique_students"},
	}

	cursor, err := r.enrollmentCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique students: %v", err)
	}
	defer cursor.Close(context.TODO())

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return 0, fmt.Errorf("failed to decode unique students count: %v", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	count, ok := result[0]["unique_students"].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected result format for unique students count")
	}

	return int64(count), nil
}
