package server

import (
	"io"

	"github.com/gorilla/mux"
	c "github.com/swethabhageerath/dataaccess/lib/client"
	"github.com/swethabhageerath/dataaccess/lib/client/connection"
	"github.com/swethabhageerath/dataaccess/lib/client/connection/config"
	"github.com/swethabhageerath/logging/lib/logger"
	"github.com/swethabhageerath/logging/lib/writers"
	"github.com/swethabhageerath/logservice/internal/middlewares"
	"github.com/swethabhageerath/logservice/internal/repositories"
	"github.com/swethabhageerath/logservice/internal/services"
	"github.com/swethabhageerath/utilities/lib/utilities/helpers"
)

type Router struct {
	handlers handlers
}

func NewRouters() Router {
	config := config.PgConfig{}
	env := helpers.EnvironmentHelper{}
	ctx := helpers.ContextHelper{}
	mar := helpers.MarshallingHelper{}
	file := helpers.FileHelper{}
	writers := []io.Writer{writers.New(env, file)}
	logger := logger.NewDefaultLogger(writers, env, ctx, mar)
	connection := connection.NewPostgresConnection(config, logger, env)

	db, err := connection.Connect()
	if err != nil {
		panic(err)
	}
	dataAccessor := c.NewExecute(db)
	repository := repositories.NewLogRepository(dataAccessor, logger)
	service := services.NewLogService(logger, repository)
	return Router{
		handlers: NewHandlers(service, logger),
	}
}

func (r Router) SetRoutes(mux *mux.Router) {
	router := mux.PathPrefix("/log").Subrouter()

	router.Handle("/create", middlewares.RequestHandler(r.handlers.createHandler)).Methods("POST")
	router.Handle("id/{id}", middlewares.RequestHandler(r.handlers.getByIdHandler)).Methods("GET")
	router.Handle("app/{app}", middlewares.RequestHandler(r.handlers.getByAppNameHandler)).Methods("GET")
	router.Handle("user/{user}", middlewares.RequestHandler(r.handlers.getByUserHandler)).Methods("GET")
	router.Handle("level/{level}", middlewares.RequestHandler(r.handlers.getByLogLevelHandler)).Methods("GET")
}
