package server

import (
	"database/sql"
	"net/http"
)

type LogicAppServer struct {
	db *sql.DB
	http.Handler
	userHandler    *UserHandler
	playersHandler *PlayersHandler
}

func NewLogicAppServer(db *sql.DB) *LogicAppServer {
	las := new(LogicAppServer)
	las.db = db
	userStore := NewMyUserStore(db)
	las.userHandler = NewUserHandler(userStore)
	las.playersHandler = NewPlayersHandler(userStore)

	router := http.NewServeMux()
	router.Handle("/user/", http.StripPrefix("/user", las.userHandler))
	router.Handle("/players/", http.StripPrefix("/players", las.playersHandler))

	las.Handler = router

	return las
}
