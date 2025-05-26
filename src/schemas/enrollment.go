package schemas

type EnrollStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}
