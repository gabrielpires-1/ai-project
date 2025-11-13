package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// 1. Carregar o .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível encontrar ou carregar o arquivo .env")
	}

	// 2. Configuração
	// Apenas verificamos se a key existe, mas não a lemos para uma var
	if os.Getenv("GOOGLE_API_KEY") == "" {
		log.Fatal("A variável de ambiente GOOGLE_API_KEY é obrigatória (defina no .env).")
	}

	log.Println("Valor da GOOGLE_API_KEY lido:", os.Getenv("GOOGLE_API_KEY"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}

	// 3. Injeção de Dependência (Simplificada)
	// NewGeminiService agora sabe como se configurar sozinho
	geminiSvc, err := NewGeminiService(ctx)
	if err != nil {
		log.Fatalf("Falha ao inicializar o GeminiService: %v", err)
	}

	// O Handler não muda
	chatHandler := NewChatHandler(geminiSvc)

	// 4. Configuração do Servidor HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHandler.HandleChat)

	log.Printf("Servidor Go rodando na porta %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Servidor falhou: %v", err)
	}
}
