package controller

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"courses-service/src/service"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModuleController struct {
	service service.ModuleServiceInterface
}

func NewModuleController(service service.ModuleServiceInterface) *ModuleController {
	return &ModuleController{
		service: service,
	}
}

// @Summary Module creation
// @Description Create a new module
// @Tags modules
// @Accept json
// @Produce json
// @Param module body schemas.CreateModuleRequest true "Module to create"
// @Success 201 {object} model.Module
// @Router /modules [post]
func (c *ModuleController) CreateModule(ctx *gin.Context) {
	slog.Debug("Creating module")

	var module schemas.CreateModuleRequest
	if err := ctx.ShouldBindJSON(&module); err != nil {
		slog.Error("Error binding JSON", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("module: %v\n", module)

	createdModule, err := c.service.CreateModule(module)
	if err != nil {
		slog.Error("Error creating module", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Debug("Module created", "module", createdModule)
	ctx.JSON(http.StatusCreated, createdModule)
}

// @Summary Get modules by course ID
// @Description Get modules by course ID
// @Tags modules
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Success 200 {array} model.Module
// @Router /modules/course/{courseId} [get]
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

// @Summary Get a module by ID
// @Description Get a module by ID
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} model.Module
// @Router /modules/{id} [get]
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

// @Summary Update a module
// @Description Update a module by ID
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param module body schemas.UpdateModuleRequest true "Module to update"
// @Success 200 {object} model.Module
// @Router /modules/{id} [put]
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

// @Summary Delete a module
// @Description Delete a module by ID
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 204 {string} string "Module deleted successfully"
// @Router /modules/{id} [delete]
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
