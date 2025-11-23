package main

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

const systemPrompt = `
Você é um Assistente Educacional Inteligente. Seu objetivo é ajudar estudantes a aprender.
Diretrizes:
1. Use o método socrático: guie o aluno à resposta em vez de dá-la pronta.
2. Explique conceitos complexos com analogias simples.
3. Use Markdown para formatar suas respostas (listas, negrito, blocos de código).
4. Seja paciente, encorajador e focado na disciplina do aluno.
5. Se o aluno pedir código, explique a lógica passo a passo.
`

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

func (s *GeminiService) GetCompletion(ctx context.Context, userMessage string) (string, error) {
	modelName := "gemini-2.5-flash"

	// Configuração onde definimos o comportamento do agente
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemPrompt},
			},
		},
		Temperature: genai.Ptr(float32(0.4)),
	}

	// Chamada à API
	result, err := s.client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(userMessage),
		config,
	)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar conteúdo: %w", err)
	}

	return result.Text(), nil
}
