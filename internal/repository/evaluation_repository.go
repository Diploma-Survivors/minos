package repository

import (
	"minos/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EvaluationRepository interface {
	CreateEvaluation(evaluation *model.Evaluation) error
	FindEvaluationByInterviewID(interviewID uuid.UUID) (*model.Evaluation, error)
}

type evaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) EvaluationRepository {
	return &evaluationRepository{db: db}
}

func (r *evaluationRepository) CreateEvaluation(evaluation *model.Evaluation) error {
	return r.db.Create(evaluation).Error
}

func (r *evaluationRepository) FindEvaluationByInterviewID(interviewID uuid.UUID) (*model.Evaluation, error) {
	var evaluation model.Evaluation
	err := r.db.Where("interview_id = ?", interviewID).First(&evaluation).Error
	if err != nil {
		return nil, err
	}
	return &evaluation, nil
}

