package repository

import (
	"AICare/internal/model" // 替换为你的项目模块名
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// AIRepository 定义数据访问接口，使业务层不依赖具体实现
type AIRepository interface {
	SendMessage(messages []model.Message) (*model.AIResponse, error)
}

type aiRepository struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewAIRepository 是构造函数，通过依赖注入配置和HTTP客户端
func NewAIRepository(apiKey, baseURL string) AIRepository {
	return &aiRepository{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *aiRepository) SendMessage(messages []model.Message) (*model.AIResponse, error) {
	aiReq := model.AIRequest{
		Model:       "deepseek-ai/DeepSeek-V3",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1024,
		Stream:      false,
	}
	//	序列化json
	jsonData, err := json.Marshal(aiReq)
	if err != nil {
		return nil, err
	}
	// 新建HTTP请求
	req, err := http.NewRequest("POST", r.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	// 发送请求
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err // 可自定义错误类型，携带更多信息
	}

	var aiResp model.AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, err
	}

	return &aiResp, nil
}
