package server

import "github.com/gorilla/mux"

type Router struct {
}

func SetRoutes(mux *mux.Router) {
	router := mux.PathPrefix("/log").Subrouter()

}
