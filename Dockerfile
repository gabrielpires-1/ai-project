# ----- Estágio 1: Build -----
FROM golang:1.25.1-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de módulo primeiro (para otimizar o cache)
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia todo o resto do código-fonte
COPY . .

# Compila o aplicativo de forma estática (CGO_ENABLED=0)
# -ldflags "-w -s" remove informações de debug, reduzindo o tamanho
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o /app/server .

# ----- Estágio 2: Final -----
# Usamos a imagem 'scratch', que é a menor possível (vazia)
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expõe a porta que sua aplicação usa (definida no main.go como 8080)
EXPOSE 8080

# Copia APENAS o binário compilado do estágio 'builder'
COPY --from=builder /app/server .

# Comando para executar a API
# (Este é o ponto de entrada consistente)
CMD ["./server"]