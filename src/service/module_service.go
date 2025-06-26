package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type ModuleService struct {
	moduleRepository repository.ModuleRepositoryInterface
}

func NewModuleService(moduleRepository repository.ModuleRepositoryInterface) *ModuleService {
	return &ModuleService{moduleRepository: moduleRepository}
}

func (s *ModuleService) CreateModule(module schemas.CreateModuleRequest) (*model.Module, error) {
	fmt.Printf("Creating module: %v\n", module)
	if _, err := s.moduleRepository.GetModuleByName(module.CourseID, module.Title); err == nil {
		return nil, fmt.Errorf("module with title %s already exists in course %s", module.Title, module.CourseID)
	}

	moduleModel := model.Module{
		Title:       module.Title,
		Description: module.Description,
		CourseID:    module.CourseID,
		Resources:   []model.ModuleResource{},
	}

	order, err := s.moduleRepository.GetNextModuleOrder(module.CourseID)
	if err != nil {
		return nil, err
	}
	moduleModel.Order = order

	return s.moduleRepository.CreateModule(module.CourseID, moduleModel)
}

func (s *ModuleService) GetModulesByCourseId(courseId string) ([]model.Module, error) {
	slog.Debug("Getting modules by course id", "courseId", courseId)
	if courseId == "" {
		return nil, errors.New("courseId is required")
	}
	return s.moduleRepository.GetModulesByCourseId(courseId)
}

func (s *ModuleService) GetModuleById(id string) (*model.Module, error) {
	slog.Debug("Getting module by id", "id", id)
	if id == "" {
		return nil, errors.New("module id is required")
	}
	return s.moduleRepository.GetModuleById(id)
}

func (s *ModuleService) GetModuleByOrder(courseID string, order int) (*model.Module, error) {
	slog.Debug("Getting module by order", "courseID", courseID, "order", order)
	if courseID == "" {
		return nil, errors.New("courseId is required")
	}
	return s.moduleRepository.GetModuleByOrder(courseID, order)
}

func (s *ModuleService) UpdateModule(id string, module model.Module) (*model.Module, error) {
	slog.Debug("Updating module", "id", id, "module", module)
	if id == "" {
		return nil, errors.New("module id is required")
	}

	// First, get the current module to check if title is changing
	currentModule, err := s.moduleRepository.GetModuleById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current module: %v", err)
	}

	// Only check for duplicate titles if the title is actually changing
	if module.Title != "" && module.Title != currentModule.Title {
		existingModule, err := s.moduleRepository.GetModuleByName(module.CourseID, module.Title)
		if err == nil {
			// Module with this title exists, check if it's a different module
			if existingModule.ID != module.ID {
				return nil, fmt.Errorf("module with title %s already exists in course %s", module.Title, module.CourseID)
			}
		} else {
			// Check if it's a "not found" error or a real error
			if !strings.Contains(err.Error(), "module not found") {
				// It's a real error (DB connection, etc.), propagate it
				return nil, fmt.Errorf("error checking module title: %v", err)
			}
			// If it's "module not found", that's fine - the title is unique
		}
	}

	return s.moduleRepository.UpdateModule(id, module)
}

func (s *ModuleService) DeleteModule(id string) error {
	slog.Debug("Deleting module", "id", id)
	if id == "" {
		return errors.New("module id is required")
	}
	return s.moduleRepository.DeleteModule(id)
}
