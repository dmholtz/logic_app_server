package server

import (
	"net/http"
)

type LogicAppServer struct {
	http.Handler
	userHandler *UserHandler
}

func NewLogicAppServer() *LogicAppServer {
	las := new(LogicAppServer)
	las.userHandler = NewUserHandler(&MyUserStore{})

	router := http.NewServeMux()
	router.Handle("/user/", http.StripPrefix("/user", las.userHandler))

	las.Handler = router

	return las
}
