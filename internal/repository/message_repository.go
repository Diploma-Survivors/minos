package repository

import (
	"minos/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *model.Message) error
	FindMessagesByInterviewID(interviewID uuid.UUID) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) FindMessagesByInterviewID(interviewID uuid.UUID) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("interview_id = ?", interviewID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

