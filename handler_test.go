package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 1. Nosso "Mock" do Serviço
// mockGeminiService simula o comportamento do GeminiService
type mockGeminiService struct {
	// Definimos o que ele deve retornar no teste
	mockAnswer string
	mockError  error
}

// Implementamos a interface GenerationService para o nosso mock
func (m *mockGeminiService) GetCompletion(ctx context.Context, prompt string) (string, error) {
	return m.mockAnswer, m.mockError
}

// 2. Testes

func TestChatHandler_HappyPath(t *testing.T) {
	// Configuração do Mock
	mockSvc := &mockGeminiService{
		mockAnswer: "Resposta simulada",
		mockError:  nil,
	}

	// Configuração do Handler
	handler := NewChatHandler(mockSvc)

	// Criar requisição
	reqBody := `{"prompt":"teste"}`
	req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(reqBody))
	rr := httptest.NewRecorder() // Grava a resposta

	// Executar
	handler.HandleChat(rr, req)

	// Verificações
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado status 200 OK, obteve %d", rr.Code)
	}

	var resp chatResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Não foi possível decodificar a resposta JSON: %v", err)
	}

	if resp.Answer != "Resposta simulada" {
		t.Errorf("Esperado 'Resposta simulada', obteve '%s'", resp.Answer)
	}
}

func TestChatHandler_ServiceError(t *testing.T) {
	// Configuração do Mock (retornando um erro)
	mockSvc := &mockGeminiService{
		mockError: errors.New("falha na IA"),
	}
	handler := NewChatHandler(mockSvc)

	reqBody := `{"prompt":"teste"}`
	req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.HandleChat(rr, req)

	// Verificação (esperamos Erro Interno do Servidor)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Esperado status 500, obteve %d", rr.Code)
	}
}

func TestChatHandler_BadMethod(t *testing.T) {
	handler := NewChatHandler(nil)                           // Nem precisamos do mock
	req := httptest.NewRequest(http.MethodGet, "/chat", nil) // Usando GET
	rr := httptest.NewRecorder()

	handler.HandleChat(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Esperado status 405, obteve %d", rr.Code)
	}
}

func TestChatHandler_InvalidJSON(t *testing.T) {
	handler := NewChatHandler(nil)
	reqBody := `{"prompt":` // JSON inválido
	req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.HandleChat(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status 400, obteve %d", rr.Code)
	}
}

func TestChatHandler_MissingPrompt(t *testing.T) {
	handler := NewChatHandler(nil)
	reqBody := `{"foo":"bar"}` // JSON válido, mas sem o campo 'prompt'
	req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	handler.HandleChat(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status 400, obteve %d", rr.Code)
	}
}
