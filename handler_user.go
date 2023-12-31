package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserHandler struct {
	http.Handler
	userStore UserStore
}

func NewUserHandler(store UserStore) *UserHandler {
	uh := new(UserHandler)
	uh.userStore = store

	router := http.NewServeMux()
	router.HandleFunc("/login", http.HandlerFunc(uh.LoginHandler))
	router.HandleFunc("/logout", http.HandlerFunc(uh.LogoutHandler))
	router.HandleFunc("/signup", http.HandlerFunc(uh.SignupHandler))
	router.HandleFunc("/reset", http.HandlerFunc(uh.ResetHandler))

	uh.Handler = router
	return uh
}

func (uh *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		http.Error(w, fmt.Sprintf("Unsupported Media Type %s", ct), http.StatusUnsupportedMediaType)
		return
	}

	// read credentials from request body
	var c Credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// login
	token, login_err := uh.userStore.Login(c)
	if login_err != nil {
		http.Error(w, login_err.Error(), http.StatusBadRequest)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (uh *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Authorization") == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := splitToken[1]

	// logout
	logout_err := uh.userStore.Logout(token)
	if logout_err != nil {
		http.Error(w, logout_err.Error(), http.StatusBadRequest)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	return
}

func (uh *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		http.Error(w, fmt.Sprintf("Unsupported Media Type %s", ct), http.StatusUnsupportedMediaType)
		return
	}

	// read credentials from request body
	var c Credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// signup
	signup_err := uh.userStore.Signup(c)
	if signup_err != nil {
		http.Error(w, signup_err.Error(), http.StatusBadRequest)
		return
	}

	// return success
	w.WriteHeader(http.StatusCreated)
	return
}

func (uh *UserHandler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get authorization token
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := splitToken[1]

	// get user from request context
	user_id, err := uh.userStore.UserIdFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// reset user
	reset_err := uh.userStore.ResetUser(user_id)
	if reset_err != nil {
		http.Error(w, reset_err.Error(), http.StatusBadRequest)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
}
