# ==============================================================================
# Stage 1: Build
# ==============================================================================
FROM golang:1.23.5-alpine AS builder

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# ==============================================================================
# Stage 2: Runtime
# ==============================================================================
FROM alpine:latest

# Instalar certificados SSL e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -s /bin/sh appuser

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário da aplicação
COPY --from=builder /app/main .

# Copiar timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Mudar para usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]
