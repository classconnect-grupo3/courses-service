package schemas

import (
	"courses-service/src/model"
	"time"
)

type CreateStudentFeedbackRequest struct {
	StudentUUID  string             `json:"student_uuid" binding:"required"`
	TeacherUUID  string             `json:"teacher_uuid" binding:"required"`
	CourseID     string             `json:"course_id" binding:"required"`
	FeedbackType model.FeedbackType `json:"feedback_type" binding:"required"`
	Score        int                `json:"score" binding:"required"`
	Feedback     string             `json:"feedback" binding:"required"`
}

type GetFeedbackByStudentIdRequest struct {
	CourseID     string             `json:"course_id"`
	FeedbackType model.FeedbackType `json:"feedback_type"`
	StartScore   int                `json:"start_score"`
	EndScore     int                `json:"end_score"`
	StartDate    time.Time          `json:"start_date"`
	EndDate      time.Time          `json:"end_date"`
}

type CreateCourseFeedbackRequest struct {
	StudentUUID  string             `json:"student_uuid" binding:"required"`
	Score        int                `json:"score" binding:"required"`
	FeedbackType model.FeedbackType `json:"feedback_type" binding:"required"`
	Feedback     string             `json:"feedback" binding:"required"`
}
