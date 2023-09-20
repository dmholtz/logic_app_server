package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type QuizHandler struct {
	http.Handler
	userStore UserStore
}

func NewQuizHandler(store UserStore) *QuizHandler {
	qh := new(QuizHandler)
	qh.userStore = store

	router := http.NewServeMux()
	router.HandleFunc("/competition", http.HandlerFunc(qh.CompetitionHandler))
	router.HandleFunc("/find", http.HandlerFunc(qh.FindHandler))
	router.HandleFunc("/solve", http.HandlerFunc(qh.SolveHandler))

	qh.Handler = router
	return qh
}

func (uh *QuizHandler) CompetitionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// query quiz database
	quiz, err := uh.userStore.GetCompetitionQuiz(user_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// return quiz
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(quiz)
}

func (uh *QuizHandler) FindHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// get query parameters
	qc, err := QuizPropertiesFromUrlQuery(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// query quiz database
	quiz, err := uh.userStore.FindQuiz(user_id, qc)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// return quiz
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(quiz)

}

func (uh *QuizHandler) SolveHandler(w http.ResponseWriter, r *http.Request) {
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
	userId, err := uh.userStore.UserIdFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// parse SolveSubmission JSON body
	var ss SolveSubmission
	err = json.NewDecoder(r.Body).Decode(&ss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// solve quiz
	subResponse, solve_err := uh.userStore.SolveQuiz(userId, ss)
	if solve_err != nil {
		http.Error(w, solve_err.Error(), http.StatusBadRequest)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subResponse)
}
