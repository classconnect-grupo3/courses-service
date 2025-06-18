package schemas

import "courses-service/src/model"

type CreateModuleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	CourseID    string `json:"course_id" binding:"required"`
}

type CreateModuleResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type UpdateModuleRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Order       int                    `json:"order"`
	Resources   []model.ModuleResource `json:"resources"`
}
