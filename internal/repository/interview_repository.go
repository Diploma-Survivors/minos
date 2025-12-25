package repository

import (
	"minos/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	CreateInterview(interview *model.Interview) error
	FindInterviewByID(id uuid.UUID) (*model.Interview, error)
	FindInterviewsByUserID(userID uuid.UUID) ([]model.Interview, error)
	UpdateInterview(interview *model.Interview) error
}

type interviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

func (r *interviewRepository) CreateInterview(interview *model.Interview) error {
	return r.db.Create(interview).Error
}

func (r *interviewRepository) FindInterviewByID(id uuid.UUID) (*model.Interview, error) {
	var interview model.Interview
	// Preload related data
	err := r.db.Preload("Messages").Preload("Submissions").Preload("Evaluation").First(&interview, id).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}

func (r *interviewRepository) FindInterviewsByUserID(userID uuid.UUID) ([]model.Interview, error) {
	var interviews []model.Interview
	err := r.db.Where("user_id = ?", userID).Order("started_at DESC").Find(&interviews).Error
	return interviews, err
}

func (r *interviewRepository) UpdateInterview(interview *model.Interview) error {
	return r.db.Save(interview).Error
}

