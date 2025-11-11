package main

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestGetCompletion_Integration é um teste de integração.
func TestGetCompletion_Integration(t *testing.T) {
	// Pulamos este teste por padrão, pois ele é "caro" (faz chamada de rede)
	// e requer uma API key válida.
	// Para rodar, use: go test -v -run TestGetCompletion_Integration
	// t.Skip("Pulando teste de integração. Remova esta linha para rodar.")

	// --- Configuração ---
	ctx := context.Background()

	// Carregar .env (necessário para GOOGLE_API_KEY)
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Erro ao carregar .env: %v", err)
	}

	if os.Getenv("GOOGLE_API_KEY") == "" {
		t.Fatal("GOOGLE_API_KEY não encontrada no .env. Teste não pode continuar.")
	}

	// Criar o serviço real
	svc, err := NewGeminiService(ctx)
	if err != nil {
		t.Fatalf("Falha ao criar NewGeminiService: %v", err)
	}

	// --- Execução ---
	prompt := "Responda apenas com a palavra 'teste'."
	answer, err := svc.GetCompletion(ctx, prompt)

	// --- Verificações ---
	if err != nil {
		t.Errorf("GetCompletion falhou: %v", err)
	}

	if answer == "" {
		t.Error("GetCompletion retornou uma resposta vazia")
	}

	t.Logf("Resposta da IA: %s", answer)
}
