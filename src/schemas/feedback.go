package schemas

import (
	"courses-service/src/model"
)

type CreateStudentFeedbackRequest struct {
	StudentUUID  string             `json:"student_uuid" binding:"required"`
	TeacherUUID  string             `json:"teacher_uuid" binding:"required"`
	CourseID     string             `json:"course_id" binding:"required"`
	FeedbackType model.FeedbackType `json:"feedback_type" binding:"required"`
	Score        int                `json:"score" binding:"required"`
	Feedback     string             `json:"feedback" binding:"required"`
}
