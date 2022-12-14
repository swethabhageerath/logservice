package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swethabhageerath/logservice/api/v1/server"
)

func run() {
	router := mux.NewRouter()
	server.NewRouters().SetRoutes(router)
	http.ListenAndServe(":80", router)
}
