package service

import (
	"AICare/internal/model"
	"AICare/internal/repository"
)

// AIService 定义业务逻辑接口
type AIService interface {
	AskQuestion(question string) (string, error)
}

type aiService struct {
	aiRepo repository.AIRepository // 依赖接口，而非具体实现
}

// NewAIService 构造函数，依赖注入 AiRepository 接口的实现
func NewAIService(aiRepo repository.AIRepository) AIService {
	return &aiService{aiRepo: aiRepo}
}

func (s *aiService) AskQuestion(question string) (string, error) {
	// 此处可扩展业务逻辑，如对话历史管理、提示词工程等
	messages := []model.Message{
		{Role: "user", Content: question},
	}

	response, err := s.aiRepo.SendMessage(messages)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "The AI did not provide a response.", nil
	}

	return response.Choices[0].Message.Content, nil
}
