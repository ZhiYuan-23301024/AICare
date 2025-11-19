// internal/model/request.go
package model

type ChatRequest struct {
	Question string `json:"question" binding:"required"`
}
