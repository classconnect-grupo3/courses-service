package service

import (
	"courses-service/src/model"
	"courses-service/src/schemas"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type AssignmentRepository interface {
	CreateAssignment(a model.Assignment) (*model.Assignment, error)
	GetAssignments() ([]*model.Assignment, error)
	GetAssignmentById(id string) (*model.Assignment, error)
	GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error)
	UpdateAssignment(id string, updateAssignment model.Assignment) (*model.Assignment, error)
	DeleteAssignment(id string) error
}

type CourseGetter interface {
	GetCourseById(id string) (*model.Course, error)
}

type AssignmentService struct {
	assignmentRepository AssignmentRepository
	courseService       CourseGetter
}

func NewAssignmentService(assignmentRepository AssignmentRepository, courseService CourseGetter) *AssignmentService {
	return &AssignmentService{
		assignmentRepository: assignmentRepository,
		courseService:       courseService,
	}
}

func (s *AssignmentService) CreateAssignment(req schemas.CreateAssignmentRequest) (*model.Assignment, error) {
	// Verificar que el curso existe
	course, err := s.courseService.GetCourseById(req.CourseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// Verificar que la fecha de entrega es posterior a la fecha actual
	now := time.Now()
	
	slog.Info("Comparing dates",
		"now", now.Format(time.RFC3339),
		"due_date", req.DueDate.Format(time.RFC3339),
		"now_location", now.Location().String(),
		"due_date_location", req.DueDate.Location().String())

	if req.DueDate.Before(now) {
		return nil, fmt.Errorf("due date must be in the future. Now: %v, DueDate: %v", 
			now.Format(time.RFC3339), 
			req.DueDate.Format(time.RFC3339))
	}

	// Solo validar la fecha de fin del curso si está establecida (no es la fecha por defecto)
	if !course.EndDate.IsZero() && req.DueDate.After(course.EndDate) {
		return nil, errors.New("due date must be before course end date")
	}

	assignment := model.Assignment{
		Title:           req.Title,
		Description:     req.Description,
		Instructions:    req.Instructions,
		Type:            req.Type,
		CourseID:        req.CourseID,
		DueDate:         req.DueDate,
		GracePeriod:     req.GracePeriod,
		Status:          req.Status,
		SubmissionRules: req.SubmissionRules,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return s.assignmentRepository.CreateAssignment(assignment)
}

func (s *AssignmentService) GetAssignments() ([]*model.Assignment, error) {
	return s.assignmentRepository.GetAssignments()
}

func (s *AssignmentService) GetAssignmentById(id string) (*model.Assignment, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.assignmentRepository.GetAssignmentById(id)
}

func (s *AssignmentService) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	if courseId == "" {
		return nil, errors.New("course id is required")
	}
	
	// Verificar que el curso existe
	_, err := s.courseService.GetCourseById(courseId)
	if err != nil {
		return nil, errors.New("course not found")
	}

	return s.assignmentRepository.GetAssignmentsByCourseId(courseId)
}

func (s *AssignmentService) UpdateAssignment(id string, req schemas.UpdateAssignmentRequest) (*model.Assignment, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Verificar que el assignment existe
	assignment, err := s.assignmentRepository.GetAssignmentById(id)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	// Si se actualiza la fecha de entrega, verificar que es válida
	var dueDate time.Time
	if !req.DueDate.IsZero() {
		course, err := s.courseService.GetCourseById(assignment.CourseID)
		if err != nil {
			return nil, errors.New("course not found")
		}

		now := time.Now()
		dueDate = req.DueDate
		
		if dueDate.Before(now) {
			return nil, errors.New("due date must be in the future")
		}
		
		// Solo validar la fecha de fin del curso si está establecida (no es la fecha por defecto)
		if !course.EndDate.IsZero() && dueDate.After(course.EndDate) {
			return nil, errors.New("due date must be before course end date")
		}
	}

	updateAssignment := model.Assignment{
		Title:           req.Title,
		Description:     req.Description,
		Instructions:    req.Instructions,
		Type:            req.Type,
		DueDate:         dueDate,
		GracePeriod:     req.GracePeriod,
		Status:          req.Status,
		SubmissionRules: req.SubmissionRules,
		UpdatedAt:       time.Now(),
	}

	return s.assignmentRepository.UpdateAssignment(id, updateAssignment)
}

func (s *AssignmentService) DeleteAssignment(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.assignmentRepository.DeleteAssignment(id)
} 