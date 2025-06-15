package schemas

type EnrollStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type EnrollStudentResponse struct {
	Message string `json:"message"`
}

type UnenrollStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
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

type UnsetFavouriteCourseRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

type UnsetFavouriteCourseResponse struct {
	Message string `json:"message"`
}
