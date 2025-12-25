package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

type Message struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InterviewID uuid.UUID   `json:"interview_id" gorm:"type:uuid;not null;index"`
	Role        MessageRole `json:"role" gorm:"type:varchar(20);not null"`
	Content     string      `json:"content" gorm:"not null"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
}

func (Message) TableName() string {
	return "messages"
}

