package main

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// GeminiService agora armazena o cliente genérico.
type GeminiService struct {
	client *genai.Client
}

// NewGeminiService usa 'nil' nas opções para carregar a key
// automaticamente da variável de ambiente GOOGLE_API_KEY.
func NewGeminiService(ctx context.Context) (*GeminiService, error) {
	// 'nil' faz o SDK procurar a variável GOOGLE_API_KEY
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente GenAI: %w", err)
	}

	return &GeminiService{
		client: client,
	}, nil
}

// GetCompletion agora usa o método client.Models.GenerateContent,
// como no seu exemplo.
func (s *GeminiService) GetCompletion(ctx context.Context, prompt string) (string, error) {

	// Modelo hardcoded, como no exemplo fornecido
	modelName := "gemini-2.5-flash"

	result, err := s.client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(prompt),
		nil, // GenerationConfig (usando o padrão)
	)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar conteúdo: %w", err)
	}

	// Usamos o helper .Text() para extrair a string,
	// que é mais simples que o parsing anterior.
	return result.Text(), nil
}
