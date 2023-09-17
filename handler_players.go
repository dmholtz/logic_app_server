package server

import (
	"encoding/json"
	"net/http"
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
