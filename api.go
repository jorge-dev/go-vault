package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiServer struct {
	listenAddr string
}

type ApiError struct {
	Error string
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func NewApiServer(listenAddr string) *ApiServer {
	return &ApiServer{listenAddr: listenAddr}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func makeHttpHandler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *ApiServer) Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /account", makeHttpHandler(s.handleGetAccount))

	log.Println("Starting server on", s.listenAddr)

	server := http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	log.Println("Server started on", s.listenAddr)
	server.ListenAndServe()
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	account := NewAccount(1, "John", "Doe", "1234567890", 100.0)
	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	return nil
}

func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// handle account
	return nil
}
