package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type StartInterviewRequest struct {
	UserID          uuid.UUID      `json:"user_id" binding:"required"`
	ProblemID       uuid.UUID      `json:"problem_id" binding:"required"`
	ProblemSnapshot datatypes.JSON `json:"problem_snapshot" binding:"required"`
}

type StartInterviewResponse struct {
	InterviewID uuid.UUID `json:"interview_id"`
	Greeting    string    `json:"greeting"`
}

type SendMessageRequest struct {
	Content  string `json:"content" binding:"required"`
	Code     string `json:"code,omitempty"`     // Optional attached code
	Language string `json:"language,omitempty"` // Optional language (e.g., "python", "go")
}

type SendMessageResponse struct {
	MessageID  uuid.UUID `json:"message_id"`
	AIResponse string    `json:"ai_response"`
}

type EndInterviewResponse struct {
	EvaluationID uuid.UUID `json:"evaluation_id"`
	OverallScore int       `json:"overall_score"`
	Feedback     string    `json:"feedback"`
}
