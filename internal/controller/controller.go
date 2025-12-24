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

func (c *Controller) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", c.HealthCheck)
	v1 := router.Group("/api/v1")
	{
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

// ===== Prompt Template Handlers =====

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
func (c *Controller) GetPromptTemplateByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid ID format", nil))
		return
	}

	template, err := c.service.GetPromptTemplateByID(uint(id))
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to fetch prompt template")
		ctx.JSON(http.StatusNotFound, model.NewResponse(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, model.NewResponse("Prompt template fetched successfully", template))
}

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
func (c *Controller) GetPromptTemplateByNameVersion(ctx *gin.Context) {
	name := ctx.Query("name")
	version := ctx.Query("version")

	if name == "" || version == "" {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Both name and version parameters are required", nil))
		return
	}

	template, err := c.service.GetPromptTemplateByNameVersion(name, version)
	if err != nil {
		log.Error().Err(err).Str("name", name).Str("version", version).Msg("Failed to fetch prompt template")
		ctx.JSON(http.StatusNotFound, model.NewResponse(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, model.NewResponse("Prompt template fetched successfully", template))
}

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
func (c *Controller) CreatePromptTemplate(ctx *gin.Context) {
	var input dto.PromptTemplateCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid input", nil))
		return
	}

	template, err := c.service.CreatePromptTemplate(&input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("Failed to create prompt template")
		// Check if it's a duplicate error
		if err.Error() == fmt.Sprintf("prompt template with name '%s' and version '%s' already exists", input.Name, input.Version) {
			ctx.JSON(http.StatusConflict, model.NewResponse(err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewResponse("Failed to create prompt template", nil))
		return
	}

	ctx.JSON(http.StatusCreated, model.NewResponse("Prompt template created successfully", template))
}

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
func (c *Controller) UpdatePromptTemplate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid ID format", nil))
		return
	}

	var input dto.PromptTemplateUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewResponse("Invalid input", nil))
		return
	}

	template, err := c.service.UpdatePromptTemplate(uint(id), &input)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to update prompt template")
		ctx.JSON(http.StatusInternalServerError, model.NewResponse(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, model.NewResponse("Prompt template updated successfully", template))
}

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
