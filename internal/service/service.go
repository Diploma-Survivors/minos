package service

import (
	"errors"
	"fmt"
	"minos/internal/dto"
	"minos/internal/model"
	"minos/internal/repository"

	"gorm.io/gorm"
)

type Service interface {
	// PromptTemplate methods
	CreatePromptTemplate(input *dto.PromptTemplateCreate) (*model.PromptTemplate, error)
	GetAllPromptTemplates(query *dto.PromptTemplateQuery) ([]model.PromptTemplate, error)
	GetPromptTemplateByID(id uint) (*model.PromptTemplate, error)
	GetPromptTemplateByNameVersion(name, version string) (*model.PromptTemplate, error)
	UpdatePromptTemplate(id uint, input *dto.PromptTemplateUpdate) (*model.PromptTemplate, error)
	DeletePromptTemplate(id uint) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// PromptTemplate methods
func (s *service) CreatePromptTemplate(input *dto.PromptTemplateCreate) (*model.PromptTemplate, error) {
	// Check if template with same name and version already exists
	existing, err := s.repo.FindPromptTemplateByNameVersion(input.Name, input.Version)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("prompt template with name '%s' and version '%s' already exists", input.Name, input.Version)
	}

	template := &model.PromptTemplate{
		Name:        input.Name,
		Version:     input.Version,
		Description: input.Description,
		Content:     input.Content,
		Variables:   input.Variables,
		IsActive:    input.IsActive,
	}

	// Default to active if not specified
	if !input.IsActive {
		template.IsActive = true
	}

	err = s.repo.CreatePromptTemplate(template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *service) GetAllPromptTemplates(query *dto.PromptTemplateQuery) ([]model.PromptTemplate, error) {
	if query == nil {
		query = &dto.PromptTemplateQuery{}
	}
	return s.repo.FindAllPromptTemplates(query.Name, query.Version, query.IsActive)
}

func (s *service) GetPromptTemplateByID(id uint) (*model.PromptTemplate, error) {
	template, err := s.repo.FindPromptTemplateByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("prompt template with id %d not found", id)
		}
		return nil, err
	}
	return template, nil
}

func (s *service) GetPromptTemplateByNameVersion(name, version string) (*model.PromptTemplate, error) {
	template, err := s.repo.FindPromptTemplateByNameVersion(name, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("prompt template with name '%s' and version '%s' not found", name, version)
		}
		return nil, err
	}
	return template, nil
}

func (s *service) UpdatePromptTemplate(id uint, input *dto.PromptTemplateUpdate) (*model.PromptTemplate, error) {
	template, err := s.repo.FindPromptTemplateByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("prompt template with id %d not found", id)
		}
		return nil, err
	}

	// Update only provided fields
	if input.Description != nil {
		template.Description = *input.Description
	}
	if input.Content != nil {
		template.Content = *input.Content
	}
	if input.Variables != nil {
		template.Variables = *input.Variables
	}
	if input.IsActive != nil {
		template.IsActive = *input.IsActive
	}

	err = s.repo.UpdatePromptTemplate(template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *service) DeletePromptTemplate(id uint) error {
	template, err := s.repo.FindPromptTemplateByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("prompt template with id %d not found", id)
		}
		return err
	}

	if template == nil {
		return fmt.Errorf("prompt template with id %d not found", id)
	}

	return s.repo.DeletePromptTemplate(id)
}
