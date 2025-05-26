package repository

import (
	"context"

	"courses-service/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	Update(ctx context.Context, submission *model.Submission) error
	GetByID(ctx context.Context, id string) (*model.Submission, error)
	GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentUUID string) (*model.Submission, error)
	GetByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error)
	GetByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error)
}

type MongoSubmissionRepository struct {
	collection *mongo.Collection
}

func NewMongoSubmissionRepository(db *mongo.Database) SubmissionRepository {
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

	var submissions []model.Submission
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

	var submissions []model.Submission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}
	return submissions, nil
}
