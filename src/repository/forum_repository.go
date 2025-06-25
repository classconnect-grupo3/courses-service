package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ForumRepository struct {
	db                 *mongo.Client
	dbName             string
	questionCollection *mongo.Collection
}

func NewForumRepository(db *mongo.Client, dbName string) *ForumRepository {
	return &ForumRepository{
		db:                 db,
		dbName:             dbName,
		questionCollection: db.Database(dbName).Collection("forum_questions"),
	}
}

// Question operations

func (r *ForumRepository) CreateQuestion(question model.ForumQuestion) (*model.ForumQuestion, error) {
	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()
	question.Status = model.QuestionStatusOpen
	question.Votes = []model.Vote{}
	question.Answers = []model.ForumAnswer{}

	_, err := r.questionCollection.InsertOne(context.TODO(), question)
	if err != nil {
		return nil, fmt.Errorf("failed to create question: %v", err)
	}

	return &question, nil
}

func (r *ForumRepository) GetQuestionById(id string) (*model.ForumQuestion, error) {
	questionUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{"_id": questionUUID}
	var question model.ForumQuestion
	err = r.questionCollection.FindOne(context.TODO(), filter).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("question with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to find question: %v", err)
	}

	return &question, nil
}

func (r *ForumRepository) GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error) {
	filter := bson.M{"course_id": courseID}

	// Sort by created_at descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.questionCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find questions: %v", err)
	}
	defer cursor.Close(context.TODO())

	var questions []model.ForumQuestion
	if err = cursor.All(context.TODO(), &questions); err != nil {
		return nil, fmt.Errorf("failed to decode questions: %v", err)
	}

	return questions, nil
}

func (r *ForumRepository) UpdateQuestion(id string, question model.ForumQuestion) (*model.ForumQuestion, error) {
	questionUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{"_id": questionUUID}
	updateFields := bson.M{}

	if question.Title != "" {
		updateFields["title"] = question.Title
	}
	if question.Description != "" {
		updateFields["description"] = question.Description
	}
	if len(question.Tags) > 0 {
		updateFields["tags"] = question.Tags
	}
	if question.Status != "" {
		updateFields["status"] = question.Status
	}

	updateFields["updated_at"] = time.Now()

	update := bson.M{"$set": updateFields}

	var updatedQuestion model.ForumQuestion
	err = r.questionCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedQuestion)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("question not found")
		}
		return nil, fmt.Errorf("failed to update question: %v", err)
	}

	return &updatedQuestion, nil
}

func (r *ForumRepository) DeleteQuestion(id string) error {
	questionUUID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{"_id": questionUUID}
	result, err := r.questionCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete question: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("question not found")
	}

	return nil
}

// Answer operations

func (r *ForumRepository) AddAnswer(questionID string, answer model.ForumAnswer) (*model.ForumAnswer, error) {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return nil, fmt.Errorf("invalid question ID: %v", err)
	}

	// Generate a unique answer ID
	answer.ID = primitive.NewObjectID().Hex()
	answer.CreatedAt = time.Now()
	answer.UpdatedAt = time.Now()
	answer.Votes = []model.Vote{}
	answer.IsAccepted = false

	filter := bson.M{"_id": questionUUID}
	update := bson.M{
		"$push": bson.M{"answers": answer},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to add answer: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("question not found")
	}

	return &answer, nil
}

func (r *ForumRepository) UpdateAnswer(questionID string, answerID string, content string) (*model.ForumAnswer, error) {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return nil, fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{
		"_id":        questionUUID,
		"answers.id": answerID,
	}
	update := bson.M{
		"$set": bson.M{
			"answers.$.content":    content,
			"answers.$.updated_at": time.Now(),
			"updated_at":           time.Now(),
		},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update answer: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("question or answer not found")
	}

	// Get the updated question to return the answer
	question, err := r.GetQuestionById(questionID)
	if err != nil {
		return nil, err
	}

	// Find and return the updated answer
	for _, ans := range question.Answers {
		if ans.ID == answerID {
			return &ans, nil
		}
	}

	return nil, fmt.Errorf("answer not found after update")
}

func (r *ForumRepository) DeleteAnswer(questionID string, answerID string) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{"_id": questionUUID}
	update := bson.M{
		"$pull": bson.M{"answers": bson.M{"id": answerID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete answer: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("question not found")
	}

	return nil
}

func (r *ForumRepository) AcceptAnswer(questionID string, answerID string) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	// First, unmark any previously accepted answer
	filter := bson.M{"_id": questionUUID}
	update := bson.M{
		"$set": bson.M{
			"answers.$[].is_accepted": false,
			"updated_at":              time.Now(),
		},
	}

	_, err = r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to unmark previous accepted answers: %v", err)
	}

	// Now mark the specific answer as accepted and update question status
	filter = bson.M{
		"_id":        questionUUID,
		"answers.id": answerID,
	}
	update = bson.M{
		"$set": bson.M{
			"answers.$.is_accepted": true,
			"accepted_answer_id":    answerID,
			"status":                model.QuestionStatusResolved,
			"updated_at":            time.Now(),
		},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to accept answer: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("question or answer not found")
	}

	return nil
}

// Vote operations

func (r *ForumRepository) AddVoteToQuestion(questionID string, userID string, voteType int) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	// First, remove any existing vote from this user
	err = r.RemoveVoteFromQuestion(questionID, userID)
	if err != nil {
		// Ignore error if no vote exists
	}

	// Add the new vote
	vote := model.Vote{
		UserID:    userID,
		VoteType:  voteType,
		CreatedAt: time.Now(),
	}

	filter := bson.M{"_id": questionUUID}
	update := bson.M{
		"$push": bson.M{"votes": vote},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add vote: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("question not found")
	}

	return nil
}

func (r *ForumRepository) AddVoteToAnswer(questionID string, answerID string, userID string, voteType int) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	// First, remove any existing vote from this user on this answer
	err = r.RemoveVoteFromAnswer(questionID, answerID, userID)
	if err != nil {
		// Ignore error if no vote exists
	}

	// Add the new vote
	vote := model.Vote{
		UserID:    userID,
		VoteType:  voteType,
		CreatedAt: time.Now(),
	}

	filter := bson.M{
		"_id":        questionUUID,
		"answers.id": answerID,
	}
	update := bson.M{
		"$push": bson.M{"answers.$.votes": vote},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to add vote to answer: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("question or answer not found")
	}

	return nil
}

func (r *ForumRepository) RemoveVoteFromQuestion(questionID string, userID string) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{"_id": questionUUID}
	update := bson.M{
		"$pull": bson.M{"votes": bson.M{"user_id": userID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %v", err)
	}

	return nil
}

func (r *ForumRepository) RemoveVoteFromAnswer(questionID string, answerID string, userID string) error {
	questionUUID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %v", err)
	}

	filter := bson.M{
		"_id":        questionUUID,
		"answers.id": answerID,
	}
	update := bson.M{
		"$pull": bson.M{"answers.$.votes": bson.M{"user_id": userID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = r.questionCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove vote from answer: %v", err)
	}

	return nil
}

// Search and filter operations

func (r *ForumRepository) SearchQuestions(courseID string, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error) {
	filter := bson.M{"course_id": courseID}

	// Add text search if query is provided
	if query != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	// Add tags filter if provided
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}

	// Add status filter if provided
	if status != "" {
		filter["status"] = status
	}

	// Set up options - sort by newest first
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.questionCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search questions: %v", err)
	}
	defer cursor.Close(context.TODO())

	var questions []model.ForumQuestion
	if err = cursor.All(context.TODO(), &questions); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %v", err)
	}

	return questions, nil
}

// CountQuestions returns the total number of questions
func (r *ForumRepository) CountQuestions() (int64, error) {
	count, err := r.questionCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count questions: %v", err)
	}
	return count, nil
}

// CountQuestionsByStatus returns the number of questions by status
func (r *ForumRepository) CountQuestionsByStatus(status model.QuestionStatus) (int64, error) {
	filter := bson.M{"status": status}
	count, err := r.questionCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count questions by status: %v", err)
	}
	return count, nil
}

// CountAnswers returns the total number of answers across all questions
func (r *ForumRepository) CountAnswers() (int64, error) {
	pipeline := []bson.M{
		{"$unwind": "$answers"},
		{"$count": "total_answers"},
	}
	
	cursor, err := r.questionCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count answers: %v", err)
	}
	defer cursor.Close(context.TODO())
	
	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return 0, fmt.Errorf("failed to decode answers count: %v", err)
	}
	
	if len(result) == 0 {
		return 0, nil
	}
	
	count, ok := result[0]["total_answers"].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected result format for answers count")
	}
	
	return int64(count), nil
}
