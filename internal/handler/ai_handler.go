package handler

import (
	"AICare/internal/model"
	"AICare/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService service.AIService // 依赖业务逻辑接口
}

// NewAIHandler 构造函数，依赖注入 AIService
func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

func (h *AIHandler) Ask(c *gin.Context) {

	fmt.Println("Received Ask request")

	var req model.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	answer, err := h.aiService.AskQuestion(req.Question)
	if err != nil {
		// 可根据错误类型细化HTTP状态码和消息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"question": req.Question,
		"answer":   answer,
	})
}
