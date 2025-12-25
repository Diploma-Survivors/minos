package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Evaluation struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InterviewID         uuid.UUID      `json:"interview_id" gorm:"type:uuid;not null;uniqueIndex"`
	ProblemSolvingScore int            `json:"problem_solving_score"`
	CodeQualityScore    int            `json:"code_quality_score"`
	CommunicationScore  int            `json:"communication_score"`
	TechnicalScore      int            `json:"technical_score"`
	OverallScore        int            `json:"overall_score"`
	Strengths           datatypes.JSON `json:"strengths" gorm:"type:jsonb"`    // []string
	Improvements        datatypes.JSON `json:"improvements" gorm:"type:jsonb"` // []string
	DetailedFeedback    string         `json:"detailed_feedback"`
	CreatedAt           time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

