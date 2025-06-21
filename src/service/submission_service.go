package service

import (
	"context"
	"errors"
	"time"

	"courses-service/src/ai"
	"courses-service/src/model"
	"courses-service/src/repository"
	"courses-service/src/schemas"
)

type SubmissionService struct {
	submissionRepo repository.SubmissionRepositoryInterface
	assignmentRepo repository.AssignmentRepositoryInterface
	courseService  CourseServiceInterface
	aiClient       *ai.AiClient
}

func NewSubmissionService(submissionRepo repository.SubmissionRepositoryInterface, assignmentRepo repository.AssignmentRepositoryInterface, courseService CourseServiceInterface, aiClient *ai.AiClient) *SubmissionService {
	return &SubmissionService{
		submissionRepo: submissionRepo,
		assignmentRepo: assignmentRepo,
		courseService:  courseService,
		aiClient:       aiClient,
	}
}

func (s *SubmissionService) CreateSubmission(ctx context.Context, submission *model.Submission) error {
	// Get assignment to validate submission
	assignment, err := s.assignmentRepo.GetByID(ctx, submission.AssignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return ErrAssignmentNotFound
	}

	// Initialize submission
	submission.CreatedAt = time.Now()
	submission.UpdatedAt = submission.CreatedAt
	submission.Status = model.SubmissionStatusDraft

	return s.submissionRepo.Create(ctx, submission)
}

func (s *SubmissionService) UpdateSubmission(ctx context.Context, submission *model.Submission) error {
	existing, err := s.submissionRepo.GetByID(ctx, submission.ID.Hex())
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSubmissionNotFound
	}

	submission.UpdatedAt = time.Now()
	return s.submissionRepo.Update(ctx, submission)
}

func (s *SubmissionService) SubmitSubmission(ctx context.Context, submissionID string) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return err
	}
	if submission == nil {
		return ErrSubmissionNotFound
	}

	assignment, err := s.assignmentRepo.GetByID(ctx, submission.AssignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return ErrAssignmentNotFound
	}

	now := time.Now()
	submission.SubmittedAt = &now
	submission.UpdatedAt = now

	// Check if submission is late
	if now.After(assignment.DueDate) {
		gracePeriodEnd := assignment.DueDate.Add(time.Duration(assignment.GracePeriod) * time.Minute)
		if now.After(gracePeriodEnd) {
			submission.Status = model.SubmissionStatusLate
		}
	} else {
		submission.Status = model.SubmissionStatusSubmitted
	}

	return s.submissionRepo.Update(ctx, submission)
}

func (s *SubmissionService) GetSubmission(ctx context.Context, id string) (*model.Submission, error) {
	return s.submissionRepo.GetByID(ctx, id)
}

func (s *SubmissionService) GetSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return s.submissionRepo.GetByAssignment(ctx, assignmentID)
}

func (s *SubmissionService) GetSubmissionsByStudent(ctx context.Context, studentUUID string) ([]model.Submission, error) {
	return s.submissionRepo.GetByStudent(ctx, studentUUID)
}

func (s *SubmissionService) GetOrCreateSubmission(ctx context.Context, assignmentID, studentUUID, studentName string) (*model.Submission, error) {
	submission, err := s.submissionRepo.GetByAssignmentAndStudent(ctx, assignmentID, studentUUID)
	if err != nil {
		return nil, err
	}

	if submission != nil {
		return submission, nil
	}

	// Create new submission
	newSubmission := &model.Submission{
		AssignmentID: assignmentID,
		StudentUUID:  studentUUID,
		StudentName:  studentName,
		Status:       model.SubmissionStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.submissionRepo.Create(ctx, newSubmission)
	if err != nil {
		return nil, err
	}

	return newSubmission, nil
}

// GradeSubmission updates the score and feedback of a submission
func (s *SubmissionService) GradeSubmission(ctx context.Context, submissionID string, score *float64, feedback string) (*model.Submission, error) {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if submission == nil {
		return nil, ErrSubmissionNotFound
	}

	// Update submission with grading information
	submission.Score = score
	submission.Feedback = feedback
	submission.UpdatedAt = time.Now()

	err = s.submissionRepo.Update(ctx, submission)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

// ValidateTeacherPermissions validates if a teacher can grade submissions for a given assignment
func (s *SubmissionService) ValidateTeacherPermissions(ctx context.Context, assignmentID, teacherUUID string) error {
	// Get assignment
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return ErrAssignmentNotFound
	}

	// Get course
	course, err := s.courseService.GetCourseById(assignment.CourseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}

	// Check if teacher is the main teacher
	if course.TeacherUUID == teacherUUID {
		return nil
	}

	// Check if teacher is an auxiliary teacher
	for _, auxTeacher := range course.AuxTeachers {
		if auxTeacher == teacherUUID {
			return nil
		}
	}

	return errors.New("teacher not authorized to grade this assignment")
}

// GenerateFeedbackSummary generates an AI summary of the feedback for a submission
func (s *SubmissionService) GenerateFeedbackSummary(ctx context.Context, submissionID string) (*schemas.AiSummaryResponse, error) {
	// Get submission
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if submission == nil {
		return nil, ErrSubmissionNotFound
	}

	// Check if submission has feedback
	if submission.Feedback == "" {
		return nil, errors.New("submission has no feedback to summarize")
	}

	// Generate summary using AI
	summary, err := s.aiClient.SummarizeSubmissionFeedback(submission.Score, submission.Feedback)
	if err != nil {
		return nil, err
	}

	return &schemas.AiSummaryResponse{
		Summary: summary,
	}, nil
}
