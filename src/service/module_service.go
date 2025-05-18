package service

import (
	"courses-service/src/model"
	"errors"
	"fmt"
	"log/slog"
)

type ModuleService struct {
	moduleRepository ModuleRepository
}

type ModuleRepository interface {
	GetNextModuleOrder(courseID string) (int, error)
	CreateModule(courseID string, module model.Module) (*model.Module, error)
	GetModuleById(id string) (*model.Module, error)
	UpdateModule(id string, module model.Module) (*model.Module, error)
	DeleteModule(id string) error
	GetModulesByCourseId(courseId string) ([]model.Module, error)
	GetModuleByName(courseID string, moduleName string) (*model.Module, error)
}

func NewModuleService(moduleRepository ModuleRepository) *ModuleService {
	return &ModuleService{moduleRepository: moduleRepository}
}

func (s *ModuleService) CreateModule(module model.Module) (*model.Module, error) {
	slog.Debug("Creating module", "module", module)
	if module.Order == 0 {
		order, err := s.moduleRepository.GetNextModuleOrder(module.CourseID)
		if err != nil {
			return nil, err
		}
		module.Order = order
	}

	if _, err := s.moduleRepository.GetModuleByName(module.CourseID, module.Title); err == nil {
		return nil, fmt.Errorf("module with title %s already exists in course %s", module.Title, module.CourseID)
	}

	return s.moduleRepository.CreateModule(module.CourseID, module)
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

func (s *ModuleService) UpdateModule(id string, module model.Module) (*model.Module, error) {
	slog.Debug("Updating module", "id", id, "module", module)
	if id == "" {
		return nil, errors.New("module id is required")
	}

	existingModule, err := s.moduleRepository.GetModuleByName(module.CourseID, module.Title)
	if err != nil {
		return nil, err
	}

	// Check if the module we are updating is the same as the existing module
	if existingModule.ID != module.ID {
		return nil, fmt.Errorf("module with title %s already exists in course %s", module.Title, module.CourseID)
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
