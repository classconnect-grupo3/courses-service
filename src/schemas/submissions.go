package schemas

// GradeSubmissionRequest represents the request to grade a submission
type GradeSubmissionRequest struct {
	Score    *float64 `json:"score" bson:"score"`
	Feedback string   `json:"feedback" bson:"feedback"`
}
