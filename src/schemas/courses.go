package schemas

import (
	"courses-service/src/model"
	"time"
)

type CreateCourseRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	TeacherID   string    `json:"teacher_id" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	TeacherName string    `json:"teacher_name"` // TODO: this will later be consulted with users service to get the teacher name
}

type CreateCourseResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   string    `json:"teacher_id"`
	TeacherName string    `json:"teacher_name"`
	Capacity    int       `json:"capacity"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateCourseRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   string    `json:"teacher_id"`
	Capacity    int       `json:"capacity"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type UpdateCourseResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   string    `json:"teacher_id"`
	Capacity    int       `json:"capacity"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetCoursesByUserIdResponse struct {
	Teacher []*model.Course `json:"teacher"`
	Student []*model.Course `json:"student"`
}

type AddAuxTeacherToCourseRequest struct {
	TeacherID    string `json:"teacher_id" binding:"required"`
	AuxTeacherID string `json:"aux_teacher_id" binding:"required"`
}

type RemoveAuxTeacherFromCourseRequest struct {
	TeacherID    string `json:"teacher_id" binding:"required"`
	AuxTeacherID string `json:"aux_teacher_id" binding:"required"`
}

type DeleteCourseResponse struct {
	Message string `json:"message"`
}

type CourseMembersResponse struct {
	TeacherID      string   `json:"teacher_id"`
	AuxTeachersIDs []string `json:"aux_teachers_ids"`
	StudentsIDs    []string `json:"students_ids"`
}