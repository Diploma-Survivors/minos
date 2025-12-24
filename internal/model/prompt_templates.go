package model

import (
	"time"

	"gorm.io/gorm"
)

// PromptTemplate represents a versioned prompt template for AI reasoning
// @Description Prompt template entity with versioning support
type PromptTemplate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;type:text"`
	Version     string         `json:"version" gorm:"not null;type:text"`
	Description string         `json:"description" gorm:"type:text"`
	Content     string         `json:"content" gorm:"not null;type:text"`
	Variables   string         `json:"variables,omitempty" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index" gorm:"default:null"`
}

// TableName specifies the table name for the PromptTemplate model
func (PromptTemplate) TableName() string {
	return "prompt_templates"
}