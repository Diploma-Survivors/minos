package dto

// PromptTemplateCreate represents the data structure for creating a new prompt template
// @Description Prompt template creation request body
type PromptTemplateCreate struct {
	// Name of the prompt template (unique per version)
	Name string `json:"name" example:"code-review-prompt" binding:"required"`

	// Version of the prompt template (semantic versioning recommended)
	Version string `json:"version" example:"v1.0.0" binding:"required"`

	// Description of what this prompt template does
	Description string `json:"description" example:"Prompt for AI code review functionality"`

	// The actual prompt content/template
	Content string `json:"content" example:"You are an expert code reviewer..." binding:"required"`

	// JSON string defining expected variables (optional)
	Variables string `json:"variables" example:"{\"language\": \"string\", \"code\": \"string\"}"`

	// Whether this template is active and can be used
	IsActive bool `json:"is_active" example:"true"`
}

// PromptTemplateUpdate represents the data structure for updating a prompt template
// @Description Prompt template update request body
type PromptTemplateUpdate struct {
	// Description of what this prompt template does
	Description *string `json:"description,omitempty" example:"Updated prompt description"`

	// The actual prompt content/template
	Content *string `json:"content,omitempty" example:"You are an expert code reviewer..."`

	// JSON string defining expected variables
	Variables *string `json:"variables,omitempty" example:"{\"language\": \"string\", \"code\": \"string\"}"`

	// Whether this template is active and can be used
	IsActive *bool `json:"is_active,omitempty" example:"false"`
}

// PromptTemplateQuery represents query parameters for listing prompt templates
// @Description Query parameters for filtering prompt templates
type PromptTemplateQuery struct {
	// Filter by template name
	Name string `form:"name" example:"code-review-prompt"`

	// Filter by version
	Version string `form:"version" example:"v1.0.0"`

	// Filter by active status
	IsActive *bool `form:"is_active" example:"true"`
}
