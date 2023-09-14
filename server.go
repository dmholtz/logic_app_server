package server

import (
	"database/sql"
	"net/http"
)

type LogicAppServer struct {
	db *sql.DB
	http.Handler
	userHandler *UserHandler
}

func NewLogicAppServer(db *sql.DB) *LogicAppServer {
	las := new(LogicAppServer)
	las.db = db
	las.userHandler = NewUserHandler(&MyUserStore{db: db})

	router := http.NewServeMux()
	router.Handle("/user/", http.StripPrefix("/user", las.userHandler))

	las.Handler = router

	return las
}
