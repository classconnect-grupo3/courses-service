package schemas

type CreateModuleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Content     string `json:"content"`
	CourseID    string `json:"course_id" binding:"required"`
}

type CreateModuleResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Content     string `json:"content"`
}

type UpdateModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Content     string `json:"content"`
}
