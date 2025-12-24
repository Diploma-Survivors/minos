package model

import (
	"time"

	"gorm.io/gorm"
)

// PromptTemplate represents a versioned prompt template for AI reasoning
// @Description Prompt template entity with versioning support
type PromptTemplate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;type:text" example:"code-review-prompt"`
	Version     string         `json:"version" gorm:"not null;type:text" example:"v1.0.0"`
	Description string         `json:"description" gorm:"type:text" example:"Prompt for AI code review functionality"`
	Content     string         `json:"content" gorm:"not null;type:text" example:"You are an expert code reviewer..."`
	Variables   string         `json:"variables,omitempty" gorm:"type:text" example:"{\"language\": \"string\"}"`
	IsActive    bool           `json:"is_active" gorm:"not null;default:true" example:"true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the PromptTemplate model
func (PromptTemplate) TableName() string {
	return "prompt_templates"
}