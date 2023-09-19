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
	quizHandler    *QuizHandler
}

func NewLogicAppServer(db *sql.DB) *LogicAppServer {
	las := new(LogicAppServer)
	las.db = db
	userStore := NewMyUserStore(db)
	las.userHandler = NewUserHandler(userStore)
	las.playersHandler = NewPlayersHandler(userStore)
	las.quizHandler = NewQuizHandler(userStore)

	router := http.NewServeMux()
	router.Handle("/players/", http.StripPrefix("/players", las.playersHandler))
	router.Handle("/quiz/", http.StripPrefix("/quiz", las.quizHandler))
	router.Handle("/user/", http.StripPrefix("/user", las.userHandler))

	las.Handler = router

	return las
}
