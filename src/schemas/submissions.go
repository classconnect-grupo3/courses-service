package schemas

// GradeSubmissionRequest represents the request to grade a submission
type GradeSubmissionRequest struct {
	Score    *float64 `json:"score" bson:"score"`
	Feedback string   `json:"feedback" bson:"feedback"`
}

// AiCorrectionResponse represents the response from AI correction
type AiCorrectionResponse struct {
	AIScore           float64 `json:"ai_score"`
	AIFeedback        string  `json:"ai_feedback"`
	NeedsManualReview bool    `json:"needs_manual_review"`
}
