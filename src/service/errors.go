package service

import "errors"

var (
	// ... existing code ...
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrAssignmentNotFound = errors.New("assignment not found")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrLateSubmission     = errors.New("submission is past due date")
)
