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
