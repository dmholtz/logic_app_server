package main

import (
	"log"
	"net/http"

	las "github.com/dmholtz/logic_app_server"
)

func main() {
	server := las.NewLogicAppServer()
	err := http.ListenAndServeTLS(":443", "cert/dev.cert.pem", "cert/dev.key.pem", server)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
