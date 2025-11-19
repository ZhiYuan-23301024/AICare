// internal/model/response.go
package model

type Message struct {
	Role    string
	Content string
}

type Choice struct {
	Message Message `json:"message"`
}

type AIResponse struct {
	Choices []Choice `json:"choices"`
}
