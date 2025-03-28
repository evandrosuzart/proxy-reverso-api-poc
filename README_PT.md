# Projeto de estudo com o  **Passo-a-Passo**: Configurando uma API Go com Docker, Nginx e mTLS

[Para acessar o conteúdo em Ingles, acesse o clique aqui](README.md)
#### Este guia detalha como criar uma API Go, configurá-la com Docker e Nginx, e adicionar autenticação TLS mútua (mTLS) para segurança.

### Parte 1: Configuração Básica com Docker e Nginx
####   O projeto de api contido nesse repositório é apenas um mock com registros fixos utilizados para testes, não trata-se de uma api funcional

    1. Criar a API Go
        Desenvolva sua API Go.
        Certifique-se de que a API escuta na porta 8080.
        Coloque o arquivo principal da aplicação no diretório raiz do projeto.

    2. Configurar o Dockerfile para a API Go

        Crie um Dockerfile na raiz do diretório da sua API Go com o seguinte conteúdo:
        Dockerfile

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

    3. Configurar o docker-compose.yml

        Crie um arquivo docker-compose.yml para orquestrar os serviços Docker.

        Inclua as configurações para sua API Go e Nginx.
        YAML

        version: '3.8'
        services:
        api:
            build: ./sua-api # caminho para o Dockerfile da api go.
            ports:
            - "8080:8080"
        nginx:
            image: nginx:latest
            ports:
            - "80:80"
            volumes:
            - ./nginx/nginx.conf:/etc/nginx/nginx.conf
            depends_on:
            - api

    4. Criar o arquivo nginx/nginx.conf

        Crie um diretório nginx e dentro dele o arquivo nginx.conf.

        Este arquivo configura o Nginx como um proxy reverso para sua API Go.
        Nginx

        events {
            worker_connections 1024;
        }

        http {
            server {
                listen 80;

                location / {
                    proxy_pass http://api:8080; # nome do serviço no docker-compose
                    proxy_set_header Host $host;
                    proxy_set_header X-Real-IP $remote_addr;
                    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                    proxy_set_header X-Forwarded-Proto $scheme;
                }
            }
        }

    5. Executar com Docker Compose

        No diretório onde docker-compose.yml está localizado, execute:
        Bash

        docker-compose up --build

    6. Testar a API Go

        Acesse sua API Go através do navegador ou com curl:
        Bash

        curl http://localhost/

### Parte 2: Adicionando mTLS
    1. Gerar Certificados (no diretório nginx/ssl)
       Navegue até o diretório nginx/ssl. Se ele não existir, crie-o.
       Execute os seguintes comandos OpenSSL para gerar certificados:
**Autoridade Certificadora (CA):**
```Bash
    openssl genrsa -out ca.key 2048
    openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt
```
**Certificado do Servidor:**
```Bash
    openssl genrsa -out server.key 2048
    openssl req -new -key server.key -out server.csr
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
    cat server.crt ca.crt >> ca.pem
```
**Certificado do Cliente:**
```Bash
    openssl genrsa -out client.key 2048
    openssl req -new -key client.key -out client.csr
    openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256
    openssl pkcs12 -export -out client.p12 -inkey client.key -in client.crt -certfile ca.crt #Opcional para navegadores.
```  

    2. Atualizar nginx/nginx.conf

        Modifique o arquivo nginx.conf para habilitar mTLS:
        Nginx

```config
    events {
        worker_connections 1024;
    }

    http {
        server {
            listen 443 ssl;

            ssl_certificate /etc/nginx/ssl/server.crt;
            ssl_certificate_key /etc/nginx/ssl/server.key;
            ssl_client_certificate /etc/nginx/ssl/ca.crt;
            ssl_verify_client on;

            location / {
                proxy_pass http://api:8080; # nome do serviço no docker-compose
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }
        }
    }
```

    3. Atualizar docker-compose.yml
       Adicione o volume dos certificados ao nginx
```YAML
    volumes:
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf
    - ./nginx/ssl:/etc/nginx/ssl
```
    4. Reiniciar Docker Compose

        Encerre os contêineres e inicie-os novamente:
```Bash
    docker-compose down
    docker-compose up --build
```
    5. Testar mTLS
       No diretório raiz do seu projeto, o seguinte comando deve funcionar:
        
```Bash
    curl --cert nginx/ssl/client.crt --key nginx/ssl/client.key --cacert nginx/ssl/ca.pem https://localhost
```
