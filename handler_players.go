package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PlayersHandler struct {
	http.Handler
	userStore UserStore
}

func NewPlayersHandler(store UserStore) *PlayersHandler {
	uh := new(PlayersHandler)
	uh.userStore = store

	router := http.NewServeMux()
	router.HandleFunc("/", http.HandlerFunc(uh.LeaderboardHandler))
	router.HandleFunc("/player/achievements", http.HandlerFunc(uh.AchievementsHandler))

	uh.Handler = router
	return uh
}

func (uh *PlayersHandler) LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// query leaderboard
	leaderboard, err := uh.userStore.Leaderboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return leaderboard
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderboard)
}

func (uh *PlayersHandler) AchievementsHandler(w http.ResponseWriter, r *http.Request) {
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

	// query achievements
	achievements, err := uh.userStore.Achievements(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return achievements
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(achievements)
}
