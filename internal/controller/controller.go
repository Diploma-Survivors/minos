package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"minos/internal/dto"
	"minos/internal/model"
	"minos/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) RegisterRoutes(router *gin.Engine, apiPrefix string) {
	v1 := router.Group(apiPrefix)
	{
		v1.GET("/health", c.HealthCheck)
		// Prompt template routes
		prompts := v1.Group("/prompts")
		{
			prompts.GET("", c.GetAllPromptTemplates)
			prompts.GET("/:id", c.GetPromptTemplateByID)
			prompts.GET("/by-name-version", c.GetPromptTemplateByNameVersion)
			prompts.POST("", c.CreatePromptTemplate)
			prompts.PUT("/:id", c.UpdatePromptTemplate)
			prompts.DELETE("/:id", c.DeletePromptTemplate)
		}
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} model.Response
// @Router /health [get]
func (x *Controller) HealthCheck(ctx *gin.Context) {
	log.Info().Msg("Health check")
	ctx.JSON(http.StatusOK, model.NewResponse("OK", nil))
}

// GetAllPromptTemplates godoc
// @Summary Get all prompt templates
// @Description Get all prompt templates with optional filters
// @Tags prompts
// @Accept json
// @Produce json
// @Param name query string false "Filter by template name"
// @Param version query string false "Filter by template version"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} model.Response{data=[]model.PromptTemplate}
// @Failure 500 {object} model.Response
// @Router /prompts [get]
func (c *Controller) GetAllPromptTemplates(ctx *gin.Context) {
	var query dto.PromptTemplateQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid query parameters", nil))
		return
	}

	templates, err := c.service.GetAllPromptTemplates(&query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch prompt templates")
		ctx.JSON(http.StatusInternalServerError, model.NewResponse("Failed to fetch prompt templates", nil))
		return
	}

	ctx.JSON(http.StatusOK, model.NewResponse("Prompt templates fetched successfully", templates))
}

// GetPromptTemplateByID godoc
// @Summary Get a prompt template by ID
// @Description Get prompt template by ID
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path int true "Prompt Template ID"
// @Success 200 {object} model.Response{data=model.PromptTemplate}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /prompts/{id} [get]

// ...

// GetPromptTemplateByNameVersion godoc
// @Summary Get a prompt template by name and version
// @Description Get prompt template by name and version combination
// @Tags prompts
// @Accept json
// @Produce json
// @Param name query string true "Template name"
// @Param version query string true "Template version"
// @Success 200 {object} model.Response{data=model.PromptTemplate}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /prompts/by-name-version [get]

// ...

// CreatePromptTemplate godoc
// @Summary Create a prompt template
// @Description Create a new prompt template with versioning
// @Tags prompts
// @Accept json
// @Produce json
// @Param prompt body dto.PromptTemplateCreate true "Create prompt template"
// @Success 201 {object} model.Response{data=model.PromptTemplate}
// @Failure 400 {object} model.Response
// @Failure 409 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /prompts [post]

// ...

// UpdatePromptTemplate godoc
// @Summary Update a prompt template
// @Description Update an existing prompt template (name and version cannot be changed)
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path int true "Prompt Template ID"
// @Param prompt body dto.PromptTemplateUpdate true "Update prompt template"
// @Success 200 {object} model.Response{data=model.PromptTemplate}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /prompts/{id} [put]

// ...

// DeletePromptTemplate godoc
// @Summary Delete a prompt template
// @Description Delete a prompt template by ID (soft delete)
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path int true "Prompt Template ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /prompts/{id} [delete]
func (c *Controller) DeletePromptTemplate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid ID format", nil))
		return
	}

	if err := c.service.DeletePromptTemplate(uint(id)); err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete prompt template")
		ctx.JSON(http.StatusInternalServerError, model.NewResponse(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, model.NewResponse("Prompt template deleted successfully", nil))
}
