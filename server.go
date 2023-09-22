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
	adminHandler   *AdminHandler
}

func NewLogicAppServer(db *sql.DB) *LogicAppServer {
	las := new(LogicAppServer)
	las.db = db
	userStore := NewMyUserStore(db)
	las.userHandler = NewUserHandler(userStore)
	las.playersHandler = NewPlayersHandler(userStore)
	las.quizHandler = NewQuizHandler(userStore)
	las.adminHandler = NewAdminHandler(userStore)

	router := http.NewServeMux()
	// handles of the REST API
	router.Handle("/players/", http.StripPrefix("/players", las.playersHandler))
	router.Handle("/quiz/", http.StripPrefix("/quiz", las.quizHandler))
	router.Handle("/user/", http.StripPrefix("/user", las.userHandler))
	// handles of the HTML admin interface
	router.Handle("/admin/", http.StripPrefix("/admin", las.adminHandler))

	// static fileserver for providing stylesheets
	router.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	las.Handler = router

	return las
}
