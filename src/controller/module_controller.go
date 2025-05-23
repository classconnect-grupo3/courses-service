package controller

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModuleController struct {
	service ModuleService
}

type ModuleService interface {
	CreateModule(module schemas.CreateModuleRequest) (*model.Module, error)
	GetModuleById(id string) (*model.Module, error)
	GetModulesByCourseId(courseId string) ([]model.Module, error)
	UpdateModule(id string, module model.Module) (*model.Module, error)
	DeleteModule(id string) error
}

func NewModuleController(service ModuleService) *ModuleController {
	return &ModuleController{
		service: service,
	}
}

func (c *ModuleController) CreateModule(ctx *gin.Context) {
	slog.Debug("Creating module")

	var module schemas.CreateModuleRequest	
	if err := ctx.ShouldBindJSON(&module); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdModule, err := c.service.CreateModule(module)
	if err != nil {
		slog.Error("Error creating module", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Module created", "module", createdModule)
	ctx.JSON(http.StatusCreated, createdModule)
}

func (c *ModuleController) GetModulesByCourseId(ctx *gin.Context) {
	slog.Debug("Getting modules by course ID")
	courseId := ctx.Param("courseId")

	modules, err := c.service.GetModulesByCourseId(courseId)
	if err != nil {
		slog.Error("Error getting modules by course ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Modules retrieved", "modules", modules)
	ctx.JSON(http.StatusOK, modules)
}

func (c *ModuleController) GetModuleById(ctx *gin.Context) {
	slog.Debug("Getting module by ID")
	id := ctx.Param("id")

	module, err := c.service.GetModuleById(id)
	if err != nil {
		slog.Error("Error getting module by ID", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Module retrieved", "module", module)
	ctx.JSON(http.StatusOK, module)
}

func (c *ModuleController) UpdateModule(ctx *gin.Context) {
	slog.Debug("Updating module")
	id := ctx.Param("id")

	var module model.Module
	if err := ctx.ShouldBindJSON(&module); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedModule, err := c.service.UpdateModule(id, module)
	if err != nil {
		slog.Error("Error updating module", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Module updated", "module", updatedModule)
	ctx.JSON(http.StatusOK, updatedModule)
}

func (c *ModuleController) DeleteModule(ctx *gin.Context) {
	slog.Debug("Deleting module")
	id := ctx.Param("id")

	err := c.service.DeleteModule(id)
	if err != nil {
		slog.Error("Error deleting module", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Module deleted", "id", id)
	ctx.JSON(http.StatusNoContent, nil)
}
