package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InterviewStatus string

const (
	InterviewStatusActive    InterviewStatus = "active"
	InterviewStatusCompleted InterviewStatus = "completed"
	InterviewStatusAbandoned InterviewStatus = "abandoned"
)

type Interview struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	ProblemID       uuid.UUID       `json:"problem_id" gorm:"type:uuid;not null"`
	ProblemSnapshot datatypes.JSON  `json:"problem_snapshot" gorm:"type:jsonb;not null"`
	Status          InterviewStatus `json:"status" gorm:"type:varchar(20);default:'active';index"`
	GeminiSessionID string          `json:"gemini_session_id" gorm:"type:varchar(255)"`
	StartedAt       time.Time       `json:"started_at" gorm:"autoCreateTime"`
	EndedAt         *time.Time      `json:"ended_at"`

	Messages    []Message    `json:"messages" gorm:"foreignKey:InterviewID"`
	Submissions []Submission `json:"submissions" gorm:"foreignKey:InterviewID"`
	Evaluation  *Evaluation  `json:"evaluation" gorm:"foreignKey:InterviewID"`
}

func (Interview) TableName() string {
	return "interviews"
}
