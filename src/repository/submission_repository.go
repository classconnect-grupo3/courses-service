package repository

import (
	"context"
	"time"

	"courses-service/src/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSubmissionRepository struct {
	collection *mongo.Collection
}

func NewMongoSubmissionRepository(db *mongo.Database) SubmissionRepositoryInterface {
	return &MongoSubmissionRepository{
		collection: db.Collection("submissions"),
	}
}

func (r *MongoSubmissionRepository) Create(ctx context.Context, submission *model.Submission) error {
	result, err := r.collection.InsertOne(ctx, submission)
	if err != nil {
		return err
	}
	submission.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MongoSubmissionRepository) Update(ctx context.Context, submission *model.Submission) error {
	filter := bson.M{"_id": submission.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, submission)
	return err
}

func (r *MongoSubmissionRepository) GetByID(ctx context.Context, id string) (*model.Submission, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var submission model.Submission
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&submission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &submission, nil
}

func (r *MongoSubmissionRepository) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error) {
	var submission model.Submission
	err := r.collection.FindOne(ctx, bson.M{
		"assignment_id": assignmentID,
		"student_uuid":  studentUUID,
	}).Decode(&submission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &submission, nil
}

func (r *MongoSubmissionRepository) GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"assignment_id": assignmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []model.Submission = make([]model.Submission, 0)
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *MongoSubmissionRepository) GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"student_uuid": studentUUID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []model.Submission = make([]model.Submission, 0)
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *MongoSubmissionRepository) DeleteByStudentAndCourse(ctx context.Context, studentUUID, courseID string) error {
	// First, get all assignments for the course
	assignmentsCursor, err := r.collection.Database().Collection("assignments").Find(ctx, bson.M{"course_id": courseID})
	if err != nil {
		return err
	}
	defer assignmentsCursor.Close(ctx)

	var assignmentIDs []string
	for assignmentsCursor.Next(ctx) {
		var assignment struct {
			ID string `bson:"_id"`
		}
		if err := assignmentsCursor.Decode(&assignment); err != nil {
			continue
		}
		assignmentIDs = append(assignmentIDs, assignment.ID)
	}

	// Delete all submissions for this student in any assignment of this course
	filter := bson.M{
		"student_uuid":  studentUUID,
		"assignment_id": bson.M{"$in": assignmentIDs},
	}

	_, err = r.collection.DeleteMany(ctx, filter)
	return err
}

// CountSubmissions returns the total number of submissions
func (r *MongoSubmissionRepository) CountSubmissions(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountSubmissionsByStatus returns the number of submissions by status
func (r *MongoSubmissionRepository) CountSubmissionsByStatus(ctx context.Context, status model.SubmissionStatus) (int64, error) {
	filter := bson.M{"status": status}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountSubmissionsThisMonth returns the number of submissions created this month
func (r *MongoSubmissionRepository) CountSubmissionsThisMonth(ctx context.Context) (int64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
