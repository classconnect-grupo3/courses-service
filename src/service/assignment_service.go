package service

import (
	"context"
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
	"errors"
	"time"
)

type AssignmentService struct {
	assignmentRepository repository.AssignmentRepository
	courseService        CourseService
}

func NewAssignmentService(assignmentRepository repository.AssignmentRepository, courseService CourseService) *AssignmentService {
	return &AssignmentService{assignmentRepository: assignmentRepository, courseService: courseService}
}

func (s *AssignmentService) GetAssignments() ([]*model.Assignment, error) {
	return s.assignmentRepository.GetAssignments()
}

func (s *AssignmentService) GetAssignmentById(id string) (*model.Assignment, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.assignmentRepository.GetByID(context.TODO(), id)
}

func (s *AssignmentService) CreateAssignment(c schemas.CreateAssignmentRequest) (*model.Assignment, error) {
	// Validate course exists
	course, err := s.courseService.GetCourseById(c.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	assignment := model.Assignment{
		Title:        c.Title,
		Description:  c.Description,
		Instructions: c.Instructions,
		Type:         c.Type,
		CourseID:     c.CourseID,
		DueDate:      c.DueDate,
		GracePeriod:  c.GracePeriod,
		Status:       c.Status,
		Questions:    c.Questions,
		TotalPoints:  c.TotalPoints,
		PassingScore: c.PassingScore,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.assignmentRepository.CreateAssignment(assignment)
}

func (s *AssignmentService) UpdateAssignment(id string, updateAssignmentRequest schemas.UpdateAssignmentRequest) (*model.Assignment, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Check if assignment exists
	existingAssignment, err := s.assignmentRepository.GetByID(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	if existingAssignment == nil {
		return nil, errors.New("assignment not found")
	}

	assignment := model.Assignment{
		Title:        updateAssignmentRequest.Title,
		Description:  updateAssignmentRequest.Description,
		Instructions: updateAssignmentRequest.Instructions,
		Type:         updateAssignmentRequest.Type,
		CourseID:     existingAssignment.CourseID,
		DueDate:      updateAssignmentRequest.DueDate,
		GracePeriod:  updateAssignmentRequest.GracePeriod,
		Status:       updateAssignmentRequest.Status,
		Questions:    updateAssignmentRequest.Questions,
		TotalPoints:  updateAssignmentRequest.TotalPoints,
		PassingScore: updateAssignmentRequest.PassingScore,
		UpdatedAt:    time.Now(),
	}

	return s.assignmentRepository.UpdateAssignment(id, assignment)
}

func (s *AssignmentService) DeleteAssignment(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.assignmentRepository.DeleteAssignment(id)
}

func (s *AssignmentService) GetAssignmentsByCourseId(courseId string) ([]*model.Assignment, error) {
	if courseId == "" {
		return nil, errors.New("course id is required")
	}
	return s.assignmentRepository.GetAssignmentsByCourseId(courseId)
}
