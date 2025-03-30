package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Client struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
	Codigo   int    `json:"codigo"`
	Links    []Link `json:"links,omitempty"`
}

type Response struct {
	Msg string `json:"msg"`
}

var clients = make(map[int]Client)
var proximoCodigo = 1

func addClientLinks(client *Client, r *http.Request) {
	codigoStr := strconv.Itoa(client.Codigo)
	client.Links = []Link{
		{Rel: "self", Href: fmt.Sprintf("%s://%s/clients/%s", getScheme(r), r.Host, codigoStr)},
		{Rel: "update", Href: fmt.Sprintf("%s://%s/clients/%s", getScheme(r), r.Host, codigoStr)},
		{Rel: "delete", Href: fmt.Sprintf("%s://%s/clients/%s", getScheme(r), r.Host, codigoStr)},
	}
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func createClient(w http.ResponseWriter, r *http.Request) {
	var novoclient Client
	err := json.NewDecoder(r.Body).Decode(&novoclient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	novoclient.Codigo = proximoCodigo
	clients[proximoCodigo] = novoclient
	proximoCodigo++

	addClientLinks(&novoclient, r) //adicionado links

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novoclient)
}

func findAllClients(w http.ResponseWriter, r *http.Request) {
	listaclients := []Client{}
	for _, client := range clients {
		tempClient := client
		addClientLinks(&tempClient, r)
		listaclients = append(listaclients, tempClient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaclients)
}

func getClientByCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	codigo, err := strconv.Atoi(params["codigo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, existe := clients[codigo]
	if !existe {
		http.NotFound(w, r)
		return
	}

	addClientLinks(&client, r) // adicionado links

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

func updateClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	codigo, err := strconv.Atoi(params["codigo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, existe := clients[codigo]
	if !existe {
		http.NotFound(w, r)
		return
	}

	var clientAtualizado Client
	err = json.NewDecoder(r.Body).Decode(&clientAtualizado)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientAtualizado.Codigo = codigo
	clients[codigo] = clientAtualizado

	addClientLinks(&clientAtualizado, r) //adicionado links

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientAtualizado)
}

func deleteClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	codigo, err := strconv.Atoi(params["codigo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, existe := clients[codigo]
	if !existe {
		http.NotFound(w, r)
		return
	}

	delete(clients, codigo)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/clients", createClient).Methods("POST")
	r.HandleFunc("/clients", findAllClients).Methods("GET")
	r.HandleFunc("/clients/{codigo}", getClientByCode).Methods("GET")
	r.HandleFunc("/clients/{codigo}", updateClient).Methods("PUT")
	r.HandleFunc("/clients/{codigo}", deleteClient).Methods("DELETE")

	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
