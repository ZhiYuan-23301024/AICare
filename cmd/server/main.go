package main

import (
	"AICare/internal/handler"
	"AICare/internal/repository"
	"AICare/internal/service"
	"AICare/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化各层依赖（依赖注入）
	// 从最底层的Repo（数据源）开始
	aiRepo := repository.NewAIRepository(cfg.AIApiKey, cfg.AIBaseURL)
	// 然后初始化依赖Repo的Service（业务逻辑）
	aiService := service.NewAIService(aiRepo)
	// 最后初始化依赖Service的Handler（控制器）
	aiHandler := handler.NewAIHandler(aiService)

	// 3. 创建Gin路由并注册
	r := gin.Default()
	r.POST("/ask", aiHandler.Ask) // 路由变得非常简洁

	// 4. 启动服务器
	r.Run(cfg.ServerPort)
}
