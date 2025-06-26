package repository

import (
	"context"
	"courses-service/src/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type TeacherActivityLogRepository struct {
	logCollection *mongo.Collection
}

// Ensure it implements the interface
var _ TeacherActivityLogRepositoryInterface = (*TeacherActivityLogRepository)(nil)

func NewTeacherActivityLogRepository(client *mongo.Client, dbName string) *TeacherActivityLogRepository {
	return &TeacherActivityLogRepository{
		logCollection: client.Database(dbName).Collection("teacher_activity_logs"),
	}
}

func (r *TeacherActivityLogRepository) LogActivity(courseID, teacherUUID, activityType, description string) error {
	log := model.TeacherActivityLog{
		CourseID:     courseID,
		TeacherUUID:  teacherUUID,
		ActivityType: activityType,
		Description:  description,
		Timestamp:    time.Now(),
	}

	_, err := r.logCollection.InsertOne(context.TODO(), log)
	if err != nil {
		return fmt.Errorf("failed to log teacher activity: %v", err)
	}

	return nil
}

func (r *TeacherActivityLogRepository) GetLogsByCourse(courseID string) ([]*model.TeacherActivityLog, error) {
	filter := map[string]interface{}{"course_id": courseID}
	
	cursor, err := r.logCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs by course: %v", err)
	}
	
	var logs []*model.TeacherActivityLog
	if err := cursor.All(context.TODO(), &logs); err != nil {
		return nil, fmt.Errorf("failed to decode logs: %v", err)
	}
	
	return logs, nil
} 