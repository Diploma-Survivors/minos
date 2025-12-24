package repository

import (
	"minos/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	// PromptTemplate methods
	CreatePromptTemplate(template *model.PromptTemplate) error
	FindAllPromptTemplates(name, version string, isActive *bool) ([]model.PromptTemplate, error)
	FindPromptTemplateByID(id uint) (*model.PromptTemplate, error)
	FindPromptTemplateByNameVersion(name, version string) (*model.PromptTemplate, error)
	UpdatePromptTemplate(template *model.PromptTemplate) error
	DeletePromptTemplate(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// PromptTemplate methods
func (r *repository) CreatePromptTemplate(template *model.PromptTemplate) error {
	return r.db.Create(template).Error
}

func (r *repository) FindAllPromptTemplates(name, version string, isActive *bool) ([]model.PromptTemplate, error) {
	var templates []model.PromptTemplate
	query := r.db.Model(&model.PromptTemplate{})

	if name != "" {
		query = query.Where("name = ?", name)
	}
	if version != "" {
		query = query.Where("version = ?", version)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *repository) FindPromptTemplateByID(id uint) (*model.PromptTemplate, error) {
	var template model.PromptTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) FindPromptTemplateByNameVersion(name, version string) (*model.PromptTemplate, error) {
	var template model.PromptTemplate
	err := r.db.Where("name = ? AND version = ?", name, version).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *repository) UpdatePromptTemplate(template *model.PromptTemplate) error {
	return r.db.Save(template).Error
}

func (r *repository) DeletePromptTemplate(id uint) error {
	return r.db.Delete(&model.PromptTemplate{}, id).Error
}
