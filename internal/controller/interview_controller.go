package controller

import (
	"minos/internal/dto"
	"minos/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InterviewController struct {
	interviewService service.InterviewService
	chatService      service.ChatService
}

func NewInterviewController(
	interviewService service.InterviewService,
	chatService service.ChatService,
) *InterviewController {
	return &InterviewController{
		interviewService: interviewService,
		chatService:      chatService,
	}
}

func (c *InterviewController) StartInterview(ctx *gin.Context) {
	var req dto.StartInterviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.interviewService.StartInterview(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *InterviewController) GetInterview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}

	res, err := c.interviewService.GetInterview(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *InterviewController) SendMessage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}

	var req dto.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.chatService.SendMessage(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *InterviewController) GetHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}

	res, err := c.chatService.GetHistory(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *InterviewController) EndInterview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}

	res, err := c.interviewService.EndInterview(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *InterviewController) RegisterRoutes(router *gin.Engine, apiPrefix string) {
	v1 := router.Group(apiPrefix)
	{
		interviews := v1.Group("/interviews")
		{
			interviews.POST("", c.StartInterview)
			interviews.GET("/:id", c.GetInterview)
			interviews.POST("/:id/messages", c.SendMessage)
			interviews.GET("/:id/messages", c.GetHistory)
			interviews.POST("/:id/end", c.EndInterview)
		}
	}
}
