package service

import (
	"courses-service/src/model"
	"courses-service/src/repository"
	"errors"
	"slices"
)

type ForumService struct {
	forumRepository  repository.ForumRepositoryInterface
	courseRepository repository.CourseRepositoryInterface
}

func NewForumService(forumRepository repository.ForumRepositoryInterface, courseRepository repository.CourseRepositoryInterface) *ForumService {
	return &ForumService{
		forumRepository:  forumRepository,
		courseRepository: courseRepository,
	}
}

// Question operations

func (s *ForumService) CreateQuestion(courseID, authorID, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error) {
	// Validate required fields
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}
	if authorID == "" {
		return nil, errors.New("author ID is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}

	// Validate course exists
	_, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// Validate tags
	if err := s.validateTags(tags); err != nil {
		return nil, err
	}

	question := model.ForumQuestion{
		CourseID:    courseID,
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Tags:        tags,
	}

	return s.forumRepository.CreateQuestion(question)
}

func (s *ForumService) GetQuestionById(id string) (*model.ForumQuestion, error) {
	if id == "" {
		return nil, errors.New("question ID is required")
	}

	return s.forumRepository.GetQuestionById(id)
}

func (s *ForumService) GetQuestionsByCourseId(courseID string) ([]model.ForumQuestion, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	// Validate course exists
	if _, err := s.courseRepository.GetCourseById(courseID); err != nil {
		return nil, errors.New("course not found")
	}

	return s.forumRepository.GetQuestionsByCourseId(courseID)
}

func (s *ForumService) UpdateQuestion(id, title, description string, tags []model.QuestionTag) (*model.ForumQuestion, error) {
	if id == "" {
		return nil, errors.New("question ID is required")
	}

	// Get existing question to validate ownership later if needed
	existingQuestion, err := s.forumRepository.GetQuestionById(id)
	if err != nil {
		return nil, err
	}

	// Validate fields if provided
	if title == "" && description == "" && len(tags) == 0 {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Validate tags if provided
	if len(tags) > 0 {
		if err := s.validateTags(tags); err != nil {
			return nil, err
		}
	}

	updateQuestion := model.ForumQuestion{
		Title:       title,
		Description: description,
		Tags:        tags,
	}

	// If only tags are being updated, preserve existing title and description
	if title == "" {
		updateQuestion.Title = existingQuestion.Title
	}
	if description == "" {
		updateQuestion.Description = existingQuestion.Description
	}
	if len(tags) == 0 {
		updateQuestion.Tags = existingQuestion.Tags
	}

	return s.forumRepository.UpdateQuestion(id, updateQuestion)
}

func (s *ForumService) DeleteQuestion(id, authorID string) error {
	if id == "" {
		return errors.New("question ID is required")
	}
	if authorID == "" {
		return errors.New("author ID is required")
	}

	// Validate question exists and check ownership
	question, err := s.forumRepository.GetQuestionById(id)
	if err != nil {
		return err
	}

	if question.AuthorID != authorID {
		return errors.New("you can only delete your own questions")
	}

	return s.forumRepository.DeleteQuestion(id)
}

// Answer operations

func (s *ForumService) AddAnswer(questionID, authorID, content string) (*model.ForumAnswer, error) {
	if questionID == "" {
		return nil, errors.New("question ID is required")
	}
	if authorID == "" {
		return nil, errors.New("author ID is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}

	// Validate question exists
	_, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return nil, err
	}

	answer := model.ForumAnswer{
		AuthorID: authorID,
		Content:  content,
	}

	return s.forumRepository.AddAnswer(questionID, answer)
}

func (s *ForumService) UpdateAnswer(questionID, answerID, authorID, content string) (*model.ForumAnswer, error) {
	if questionID == "" {
		return nil, errors.New("question ID is required")
	}
	if answerID == "" {
		return nil, errors.New("answer ID is required")
	}
	if authorID == "" {
		return nil, errors.New("author ID is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}

	// Validate question exists and check answer ownership
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return nil, err
	}

	// Find the answer and check ownership
	var answerFound bool
	for _, answer := range question.Answers {
		if answer.ID == answerID {
			if answer.AuthorID != authorID {
				return nil, errors.New("you can only update your own answers")
			}
			answerFound = true
			break
		}
	}

	if !answerFound {
		return nil, errors.New("answer not found")
	}

	return s.forumRepository.UpdateAnswer(questionID, answerID, content)
}

func (s *ForumService) DeleteAnswer(questionID, answerID, authorID string) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if answerID == "" {
		return errors.New("answer ID is required")
	}
	if authorID == "" {
		return errors.New("author ID is required")
	}

	// Validate question exists and check answer ownership
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	// Find the answer and check ownership
	var answerFound bool
	for _, answer := range question.Answers {
		if answer.ID == answerID {
			if answer.AuthorID != authorID {
				return errors.New("you can only delete your own answers")
			}
			answerFound = true
			break
		}
	}

	if !answerFound {
		return errors.New("answer not found")
	}

	return s.forumRepository.DeleteAnswer(questionID, answerID)
}

func (s *ForumService) AcceptAnswer(questionID, answerID, authorID string) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if answerID == "" {
		return errors.New("answer ID is required")
	}
	if authorID == "" {
		return errors.New("author ID is required")
	}

	// Validate question exists and check question ownership
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	if question.AuthorID != authorID {
		return errors.New("only the question author can accept answers")
	}

	// Validate answer exists
	var answerFound bool
	for _, answer := range question.Answers {
		if answer.ID == answerID {
			answerFound = true
			break
		}
	}

	if !answerFound {
		return errors.New("answer not found")
	}

	return s.forumRepository.AcceptAnswer(questionID, answerID)
}

// Vote operations

func (s *ForumService) VoteQuestion(questionID, userID string, voteType int) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}
	if voteType != model.VoteTypeUp && voteType != model.VoteTypeDown {
		return errors.New("invalid vote type")
	}

	// Validate question exists
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	// Check if user is voting on their own question
	if question.AuthorID == userID {
		return errors.New("you cannot vote on your own question")
	}

	return s.forumRepository.AddVoteToQuestion(questionID, userID, voteType)
}

func (s *ForumService) VoteAnswer(questionID, answerID, userID string, voteType int) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if answerID == "" {
		return errors.New("answer ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}
	if voteType != model.VoteTypeUp && voteType != model.VoteTypeDown {
		return errors.New("invalid vote type")
	}

	// Validate question and answer exist
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	// Find the answer and check if user is voting on their own answer
	var answerFound bool
	for _, answer := range question.Answers {
		if answer.ID == answerID {
			if answer.AuthorID == userID {
				return errors.New("you cannot vote on your own answer")
			}
			answerFound = true
			break
		}
	}

	if !answerFound {
		return errors.New("answer not found")
	}

	return s.forumRepository.AddVoteToAnswer(questionID, answerID, userID, voteType)
}

func (s *ForumService) RemoveVoteFromQuestion(questionID, userID string) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Validate question exists
	_, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	return s.forumRepository.RemoveVoteFromQuestion(questionID, userID)
}

func (s *ForumService) RemoveVoteFromAnswer(questionID, answerID, userID string) error {
	if questionID == "" {
		return errors.New("question ID is required")
	}
	if answerID == "" {
		return errors.New("answer ID is required")
	}
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Validate question and answer exist
	question, err := s.forumRepository.GetQuestionById(questionID)
	if err != nil {
		return err
	}

	// Find the answer
	var answerFound bool
	for _, answer := range question.Answers {
		if answer.ID == answerID {
			answerFound = true
			break
		}
	}

	if !answerFound {
		return errors.New("answer not found")
	}

	return s.forumRepository.RemoveVoteFromAnswer(questionID, answerID, userID)
}

// Search and filter operations

func (s *ForumService) SearchQuestions(courseID, query string, tags []model.QuestionTag, status model.QuestionStatus) ([]model.ForumQuestion, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	// Validate course exists
	_, err := s.courseRepository.GetCourseById(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// Validate tags if provided
	if len(tags) > 0 {
		if err := s.validateTags(tags); err != nil {
			return nil, err
		}
	}

	// Validate status if provided
	if status != "" && !s.isValidStatus(status) {
		return nil, errors.New("invalid question status")
	}

	return s.forumRepository.SearchQuestions(courseID, query, tags, status)
}

// Helper methods

func (s *ForumService) validateTags(tags []model.QuestionTag) error {
	if len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		if !slices.Contains(model.QuestionTagValues, tag) {
			return errors.New("invalid tag: " + string(tag))
		}
	}
	return nil
}

func (s *ForumService) isValidStatus(status model.QuestionStatus) bool {
	return slices.Contains(model.QuestionStatusValues, status)
}
