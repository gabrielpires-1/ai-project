package main

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestGetCompletion_Integration é um teste de integração.
func TestGetCompletion_Integration(t *testing.T) {
	// 1. Tenta carregar o .env (para testes locais)
	// Usamos _ para ignorar o erro "no such file or directory" na CI
	_ = godotenv.Load()

	// 2. Verifica se a key existe (seja do .env local ou dos secrets da CI)
	if os.Getenv("GOOGLE_API_KEY") == "" {
		t.Fatal("GOOGLE_API_KEY não encontrada. Defina no .env local ou nos Secrets do repositório.")
	}

	// --- Configuração ---
	ctx := context.Background()

	svc, err := NewGeminiService(ctx)
	if err != nil {
		t.Fatalf("Falha ao criar NewGeminiService: %v", err)
	}

	prompt := "Responda apenas com a palavra 'teste'."
	answer, err := svc.GetCompletion(ctx, prompt)

	if err != nil {
		t.Errorf("GetCompletion falhou: %v", err)
	}

	if answer == "" {
		t.Error("GetCompletion retornou uma resposta vazia")
	}

	t.Logf("Resposta da IA: %s", answer)
}
