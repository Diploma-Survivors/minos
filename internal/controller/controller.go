package controller

import "github.com/gin-gonic/gin"

type Controller struct {
	PromptTemplate *PromptTemplateController
	Interview      *InterviewController
}

func NewController(pt *PromptTemplateController, interview *InterviewController) *Controller {
	return &Controller{
		PromptTemplate: pt,
		Interview:      interview,
	}
}

func (c *Controller) RegisterRoutes(router *gin.Engine, apiPrefix string) {
	c.PromptTemplate.RegisterRoutes(router, apiPrefix)
	c.Interview.RegisterRoutes(router, apiPrefix)
}

