package main

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strconv"
        "strings"
)

type Response struct {
        Msg string `json:"msg"`
}

type Client struct {
        Nome     string `json:"nome"`
        Telefone string `json:"telefone"`
        Email    string `json:"email"`
        Codigo   int    `json:"codigo"`
}

var clients = []Client{
        {Nome: "João da Silva", Telefone: "(11) 9999-9999", Email: "joao@example.com", Codigo: 1},
        {Nome: "Maria Oliveira", Telefone: "(21) 8888-8888", Email: "maria@example.com", Codigo: 2},
        {Nome: "Pedro Souza", Telefone: "(31) 7777-7777", Email: "pedro@example.com", Codigo: 3},
}

func clientsHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(clients)
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        // Obtem o id da URL.
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) != 3 { // /clients/1
                http.Error(w, "Requisição inválida", http.StatusBadRequest)
                return
        }
        id, err := strconv.Atoi(parts[2])
        if err != nil {
                http.Error(w, "ID inválido", http.StatusBadRequest)
                return
        }

        // Busca o cliente pelo ID.
        for _, client := range clients {
                if client.Codigo == id {
                        json.NewEncoder(w).Encode(client)
                        return
                }
        }

        // Cliente não encontrado.
        w.WriteHeader(http.StatusNotFound)
        response := Response{Msg: "Cliente não encontrado"}
        json.NewEncoder(w).Encode(response)
}

func main() {
        http.HandleFunc("/clients", clientsHandler)
        http.HandleFunc("/clients/", clientHandler) // Adicionado a nova rota.
        fmt.Println("Servidor rodando na porta 8080...")
        http.ListenAndServe(":8080", nil)
}