package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// GenerationService define o contrato que nosso handler espera.
type GenerationService interface {
	GetCompletion(ctx context.Context, prompt string) (string, error)
}

// --- DTOs (Data Transfer Objects) ---

type chatRequest struct {
	Prompt string `json:"prompt"`
}

type chatResponse struct {
	Answer string `json:"answer"`
}

// --- Handler ---

// ChatHandler agora depende da interface
type ChatHandler struct {
	service GenerationService
}

// NewChatHandler agora aceita a interface
func NewChatHandler(s GenerationService) *ChatHandler {
	return &ChatHandler{
		service: s,
	}
}

// HandleChat é a função que responde às requisições POST.
// (Esta era a parte que estava faltando)
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	// 1. Validar o método
	if r.Method != http.MethodPost {
		httpError(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decodificar o JSON de entrada
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validação simples
	if req.Prompt == "" {
		httpError(w, "O campo 'prompt' é obrigatório", http.StatusBadRequest)
		return
	}

	// 4. Chamar a lógica de negócio (o serviço)
	ctx := r.Context()
	answer, err := h.service.GetCompletion(ctx, req.Prompt)
	if err != nil {
		log.Printf("Erro ao chamar o serviço Gemini: %v", err)
		httpError(w, "Erro interno ao processar a solicitação", http.StatusInternalServerError)
		return
	}

	// 5. Codificar a resposta JSON
	resp := chatResponse{
		Answer: answer,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		// Se falhar aqui, apenas logamos, pois o header já foi enviado
		log.Printf("Erro ao escrever a resposta JSON: %v", err)
	}
}

// Função helper para padronizar respostas de erro
func httpError(w http.ResponseWriter, message string, code int) {
	// Garante que o content-type do erro também é JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
