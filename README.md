# Study Project: Step-by-Step: Configuring a Go API with Docker, Nginx, and mTLS
[To access the content in PT-BR, click here](README_PT.md)

#### This guide details how to create a Go API, configure it with Docker and Nginx, and add mutual TLS (mTLS) authentication for security.

### Part 1: Basic Configuration with Docker and Nginx
#### The API project contained in this repository is just a mock with fixed records used for testing; it's not a functional API.

    1. Create the Go API
        Develop your Go API.
        Ensure that the API listens on port 8080.
        Place the main application file in the root directory of the project.

    2. Configure the Dockerfile for the Go API

        Create a Dockerfile in the root directory of your Go API with the following content:
        Dockerfile

        # Uses the official Go image as base
        FROM golang:1.23.1-alpine AS builder

        # Sets the working directory inside the container
        WORKDIR /app

        # Copies the go.mod and go.sum files and downloads the dependencies
        COPY go.mod go.sum ./
        RUN go mod download

        # Copies the rest of the files from your project
        COPY . .

        # Compiles the Go code
        RUN go build -o main .

        # Uses a minimal alpine image for the final container
        FROM alpine:latest

        # Sets the working directory
        WORKDIR /app

        # Copies the compiled executable from the previous stage
        COPY --from=builder /app/main .

        # Exposes the port that the Go application uses
        EXPOSE 8080

        # Sets the command to run the application
        CMD ["./main"]

    3. Configure docker-compose.yml

        Create a docker-compose.yml file to orchestrate Docker services.

        Include the settings for your Go API and Nginx.
        YAML

        version: '3.8'
        services:
        api:
            build: ./sua-api # path to the Go API Dockerfile.
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

    4. Create the nginx/nginx.conf file

        Create an nginx directory and inside it, the nginx.conf file.

        This file configures Nginx as a reverse proxy for your Go API.
        Nginx

        events {
            worker_connections 1024;
        }

        http {
            server {
                listen 80;

                location / {
                    proxy_pass http://api:8080; # service name in docker-compose
                    proxy_set_header Host $host;
                    proxy_set_header X-Real-IP $remote_addr;
                    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                    proxy_set_header X-Forwarded-Proto $scheme;
                }
            }
        }

    5. Run with Docker Compose

        In the directory where docker-compose.yml is located, execute:
        Bash

        docker-compose up --build

    6. Test the Go API

        Access your Go API through the browser or with curl:
        Bash

        curl http://localhost/

### Part 2: Adding mTLS
    1. Generate Certificates (in the nginx/ssl directory)
        Navigate to the nginx/ssl directory. If it does not exist, create it.
        Execute the following OpenSSL commands to generate certificates:
**Certificate Authority (CA):**
```Bash
        openssl genrsa -out ca.key 2048
        openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt