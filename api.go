package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const JWT_SECRET = "secret"

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
	router.HandleFunc("GET /account", withJwtOut(makeHttpHandler(s.handleGetAccounts)))
	router.HandleFunc("POST /account", makeHttpHandler(s.handleCreateAccount))
	router.HandleFunc("PUT /account/{id}", withJwtOut(makeHttpHandler(s.handleUpdateAccount)))
	router.HandleFunc("DELETE /account/{id}", withJwtOut(makeHttpHandler(s.handleDeleteAccount)))
	router.HandleFunc("GET /account/{id}", withJwtOut(makeHttpHandler(s.handleGetAccountById)))

	router.HandleFunc("POST /transfer", withJwtOut(makeHttpHandler(s.handleTransfer)))

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
	validId, err := validateId(r)
	if err != nil {
		return err
	}

	account, err := s.Store.GetAccountById(validId)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, account)

}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	createAcctReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAcctReq); err != nil {
		log.Default().Println("Error decoding request:", err)
		return err
	}
	defer r.Body.Close()

	account := NewAccount(createAcctReq.FirstName, createAcctReq.LastName, createAcctReq.Balance)

	if err := s.Store.CreateAccount(account); err != nil {
		log.Default().Println("Error creating account:", err)
		return err
	}

	tokenString, err := GenerateToken(account.Id)
	if err != nil {
		log.Default().Println("Error generating token:", err)
		return fmt.Errorf("Error generating token: %v", err)
	}

	return writeJSON(w, http.StatusOK, map[string]any{"status": "created", "token": tokenString, "account": account})
}

func (s *ApiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	// handle account
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	validId, err := validateId(r)
	if err != nil {
		return err
	}

	if err := s.Store.DeleteAccount(validId); err != nil {
		log.Default().Println("Error deleting account:", err)
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": strconv.Itoa(validId)})

}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		log.Default().Println("Error decoding request:", err)
		return err
	}
	defer r.Body.Close()

	return writeJSON(w, http.StatusOK, transferReq)
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

func validateId(r *http.Request) (int, error) {
	// validate id
	pathId := r.PathValue("id")
	validId, err := strconv.Atoi(pathId)
	if err != nil {
		log.Default().Printf("Error converting id to int: %v", err)
		return 0, fmt.Errorf("Invalid id given: %v", pathId)
	}
	return validId, nil

}

func withJwtOut(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate jwt
		tokenString := r.Header.Get("Authorization")
		if !strings.HasPrefix(tokenString, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, ApiError{Error: "token not found"})
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		_, err := ValidateToken(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     userId,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		errorMessage := fmt.Sprintf("Error generating token: %v", err)
		return "", errors.New(errorMessage)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	_, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("error getting claims")
	}
	return parsedToken, nil

}
