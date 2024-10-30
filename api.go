package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiServer struct {
	listenAddr string
	Store      Storage
}

type ApiError struct {
	Error string
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{listenAddr: listenAddr, Store: store}
}

func (s *ApiServer) Start() {
	router := http.NewServeMux()
	router.HandleFunc("GET /account", makeHttpHandler(s.handleGetAccounts))
	router.HandleFunc("POST /account", makeHttpHandler(s.handleCreateAccount))
	router.HandleFunc("GET /account/{id}", makeHttpHandler(s.handleGetAccountById))

	log.Println("Starting server on", s.listenAddr)

	server := http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	log.Println("Server started on", s.listenAddr)
	server.ListenAndServe()
}

func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	// handle account
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		log.Default().Println("Error getting accounts hander:", err)
		return err
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	log.Println("id:", id)
	// handle account
	return writeJSON(w, http.StatusOK, &Account{})
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	createAcctReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAcctReq); err != nil {
		log.Default().Println("Error decoding request:", err)
		return err
	}

	account := NewAccount(createAcctReq.FirstName, createAcctReq.LastName, createAcctReq.Balance)

	if err := s.Store.CreateAccount(account); err != nil {
		log.Default().Println("Error creating account:", err)
		return err
	}

	return writeJSON(w, http.StatusOK, account)
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

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func makeHttpHandler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
