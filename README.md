# Agente Educacional (Go + Gemini)

Sistema backend em Go que utiliza a IA do Google Gemini para atuar como assistente de estudos.

## 1. Configuração (Obrigatório)

Antes de rodar, crie um arquivo chamado `.env` na raiz do projeto e adicione sua chave:

```env
GOOGLE_API_KEY=sua_chave_do_google_aqui
PORT=8080
```
## 2. Como rodar em Docker (Recomendado)

Certifique-se de estar na pasta do projeto e execute:

```bash
# 1. Construir a imagem
docker build -t agente-edu .

# 2. Rodar o container (lendo as variáveis do arquivo .env)
docker run -d -p 8080:8080 --env-file .env --name meu-agente agente-edu
```

## 3. Como Rodar Localmente (Sem Docker)

Se tiver o Go instalado:

```bash
# Baixar dependências
go mod tidy

# Rodar o servidor
go run .
```

## 4. Como testar

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Como posso melhorar meus estudos em Go?"}'
```