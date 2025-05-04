package schemas

import "time"

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	TeacherID   string `json:"teacher_id" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required"`
}

type CourseResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   string    `json:"teacher_id"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TeacherID   string `json:"teacher_id"`
	Capacity    int    `json:"capacity"`
}

type UpdateCourseResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   string    `json:"teacher_id"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
