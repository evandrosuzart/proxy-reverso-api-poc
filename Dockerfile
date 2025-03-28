# Usa a imagem oficial do Go como base
FROM golang:1.23.1-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante dos arquivos do seu projeto
COPY . .

# Compila o código Go
RUN go build -o main .

# Usa uma imagem alpine mínima para o container final
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /app

# Copia o executável compilado do estágio anterior
COPY --from=builder /app/main .

# Exponha a porta que a aplicação Go utiliza
EXPOSE 8080

# Define o comando para executar a aplicação
CMD ["./main"]
