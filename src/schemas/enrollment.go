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

