package ai

import (
	"context"
)

// Provider define la interfaz común para todos los proveedores de IA
type Provider interface {
	// Chat envía un mensaje y recibe una respuesta
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	
	// ChatStreaming envía un mensaje y recibe respuesta por streaming
	ChatStreaming(ctx context.Context, req *ChatRequest, handler StreamHandler) error
	
	// IsAvailable verifica si el proveedor está disponible
	IsAvailable(ctx context.Context) bool
	
	// GetType retorna el tipo de proveedor
	GetType() string
}

// ChatRequest representa una petición de chat
type ChatRequest struct {
	Prompt      string
	Model       string
	Temperature float64
	MaxTokens   int
	SystemPrompt string
	History     []Message
}

// ChatResponse representa una respuesta de chat
type ChatResponse struct {
	Content   string
	Model     string
	TokensUsed int
	LatencyMs int64
}

// Message representa un mensaje en el historial de chat
type Message struct {
	Role    string // user, assistant, system
	Content string
}

// StreamHandler maneja chunks de respuesta en streaming
type StreamHandler func(chunk string) error

// TODO: Implementar proveedores concretos en:
// - omniroute.go (OmniRouteProvider)
// - ollama.go (OllamaProvider)
// - openai.go (OpenAIProvider)
// - deepseek.go (DeepSeekProvider)
