package repository

import (
	"minos/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateSubmission(submission *model.Submission) error
	FindSubmissionsByInterviewID(interviewID uuid.UUID) ([]model.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) CreateSubmission(submission *model.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) FindSubmissionsByInterviewID(interviewID uuid.UUID) ([]model.Submission, error) {
	var submissions []model.Submission
	err := r.db.Where("interview_id = ?", interviewID).Order("submitted_at DESC").Find(&submissions).Error
	return submissions, err
}

