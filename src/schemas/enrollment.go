package schemas

type EnrollStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type EnrollStudentResponse struct {
	Message string `json:"message"`
}

type UnenrollStudentResponse struct {
	Message string `json:"message"`
}

type SetFavouriteCourseRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type SetFavouriteCourseResponse struct {
	Message string `json:"message"`
}

type UnsetFavouriteCourseResponse struct {
	Message string `json:"message"`
}

type ApproveStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type ApproveStudentResponse struct {
	Message   string `json:"message"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
}

type DisapproveStudentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type DisapproveStudentResponse struct {
	Message   string `json:"message"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	Reason    string `json:"reason"`
}
