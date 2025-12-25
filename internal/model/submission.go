package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Submission struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InterviewID uuid.UUID      `json:"interview_id" gorm:"type:uuid;not null;index"`
	Code        string         `json:"code" gorm:"not null"`
	Language    string         `json:"language" gorm:"type:varchar(20);not null"`
	AIFeedback  string         `json:"ai_feedback"`
	IsCorrect   *bool          `json:"is_correct"`
	TestResults datatypes.JSON `json:"test_results" gorm:"type:jsonb"` // Store simulated test results
	SubmittedAt time.Time      `json:"submitted_at" gorm:"autoCreateTime"`
}

func (Submission) TableName() string {
	return "submissions"
}

