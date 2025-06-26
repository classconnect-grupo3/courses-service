package repository

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/schemas"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssignmentRepository struct {
	db                   *mongo.Client
	dbName               string
	assignmentCollection *mongo.Collection
}

func NewAssignmentRepository(db *mongo.Client, dbName string) AssignmentRepositoryInterface {
	return &AssignmentRepository{
		db:                   db,
		dbName:               dbName,
		assignmentCollection: db.Database(dbName).Collection("assignments"),
	}
}

func (r *AssignmentRepository) CreateAssignment(assignment model.Assignment) (*model.Assignment, error) {
	result, err := r.assignmentCollection.InsertOne(context.TODO(), assignment)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignment: %v", err)
	}

	assignment.ID = result.InsertedID.(primitive.ObjectID)
	return &assignment, nil
}

func (r *AssignmentRepository) GetAssignments() ([]*model.Assignment, error) {
	cursor, err := r.assignmentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %v", err)
	}

	var assignments []*model.Assignment
	if err := cursor.All(context.TODO(), &assignments); err != nil {
		return nil, fmt.Errorf("failed to get assignments: %v", err)
	}

	// Ensure we return an empty slice instead of nil when no documents are found
	if assignments == nil {
		assignments = []*model.Assignment{}
	}

	return assignments, nil
}

func (r *AssignmentRepository) GetByID(ctx context.Context, id string) (*model.Assignment, error) {
	var assignment model.Assignment
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment by id: %v", err)
	}

	err = r.assignmentCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&assignment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get assignment by id: %v", err)
	}

	return &assignment, nil
}

func (r *AssignmentRepository) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	cursor, err := r.assignmentCollection.Find(context.TODO(), bson.M{"course_id": courseId})
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments by course id: %v", err)
	}

	var assignments []*model.Assignment
	if err := cursor.All(context.TODO(), &assignments); err != nil {
		return nil, fmt.Errorf("failed to get assignments by course id: %v", err)
	}

	// Ensure we return an empty slice instead of nil when no documents are found
	if assignments == nil {
		assignments = []*model.Assignment{}
	}

	return assignments, nil
}

func filterEmptyAssignmentFields(assignment model.Assignment) bson.M {
	update := bson.M{}

	if assignment.Title != "" {
		update["title"] = assignment.Title
	}
	if assignment.Description != "" {
		update["description"] = assignment.Description
	}
	if assignment.Type != "" {
		update["type"] = assignment.Type
	}
	if !assignment.DueDate.IsZero() {
		update["due_date"] = assignment.DueDate
	}
	if assignment.Status != "" {
		update["status"] = assignment.Status
	}
	if assignment.GracePeriod > 0 {
		update["grace_period"] = assignment.GracePeriod
	}
	if len(assignment.SubmissionRules) > 0 {
		update["submission_rules"] = assignment.SubmissionRules
	}
	if assignment.Instructions != "" {
		update["instructions"] = assignment.Instructions
	}
	if len(assignment.Questions) > 0 {
		update["questions"] = assignment.Questions
	}
	if assignment.TotalPoints > 0 {
		update["total_points"] = assignment.TotalPoints
	}
	if assignment.PassingScore > 0 {
		update["passing_score"] = assignment.PassingScore
	}
	update["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return update
}

func (r *AssignmentRepository) UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignment: %v", err)
	}

	update := filterEmptyAssignmentFields(updateAssignment)

	_, err = r.assignmentCollection.UpdateOne(context.TODO(), bson.M{"_id": objectId}, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update assignment: %v", err)
	}

	updatedAssignment, err := r.GetByID(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated assignment: %v", err)
	}

	return updatedAssignment, nil
}

func (r *AssignmentRepository) DeleteAssignment(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %v", err)
	}

	_, err = r.assignmentCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %v", err)
	}

	return nil
}

// CountAssignments returns the total number of assignments
func (r *AssignmentRepository) CountAssignments() (int64, error) {
	count, err := r.assignmentCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments: %v", err)
	}
	return count, nil
}

// CountAssignmentsByType returns the number of assignments by type
func (r *AssignmentRepository) CountAssignmentsByType(assignmentType string) (int64, error) {
	filter := bson.M{"type": assignmentType}
	count, err := r.assignmentCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments by type: %v", err)
	}
	return count, nil
}

// CountAssignmentsByStatus returns the number of assignments by status
func (r *AssignmentRepository) CountAssignmentsByStatus(status string) (int64, error) {
	filter := bson.M{"status": status}
	count, err := r.assignmentCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments by status: %v", err)
	}
	return count, nil
}

// CountAssignmentsCreatedThisMonth returns the number of assignments created this month
func (r *AssignmentRepository) CountAssignmentsCreatedThisMonth() (int64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	count, err := r.assignmentCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments created this month: %v", err)
	}
	return count, nil
}

// GetAssignmentDistribution returns assignment distribution by type and status
func (r *AssignmentRepository) GetAssignmentDistribution() ([]schemas.AssignmentDistribution, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id": bson.M{
				"type":   "$type",
				"status": "$status",
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"type":   "$_id.type",
			"status": "$_id.status",
			"count":  1,
			"_id":    0,
		}},
		{"$sort": bson.M{"type": 1, "status": 1}},
	}

	cursor, err := r.assignmentCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment distribution: %v", err)
	}
	defer cursor.Close(context.TODO())

	var distribution []schemas.AssignmentDistribution
	if err = cursor.All(context.TODO(), &distribution); err != nil {
		return nil, fmt.Errorf("failed to decode assignment distribution: %v", err)
	}

	return distribution, nil
}

// GetRecentAssignments returns recent assignments with basic information
func (r *AssignmentRepository) GetRecentAssignments(limit int) ([]schemas.AssignmentBasicInfo, error) {
	pipeline := []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{"$limit": limit},
		{"$project": bson.M{
			"id":         bson.M{"$toString": "$_id"},
			"title":      1,
			"type":       1,
			"status":     1,
			"course_id":  1,
			"created_at": 1,
			"due_date":   1,
		}},
	}

	cursor, err := r.assignmentCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent assignments: %v", err)
	}
	defer cursor.Close(context.TODO())

	var assignments []schemas.AssignmentBasicInfo
	if err = cursor.All(context.TODO(), &assignments); err != nil {
		return nil, fmt.Errorf("failed to decode recent assignments: %v", err)
	}

	return assignments, nil
}
