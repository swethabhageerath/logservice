package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/swethabhageerath/logging/lib/logger"
	"github.com/swethabhageerath/logservice/internal/models"
)

type LogServicer interface {
	Create(ctx context.Context, le models.LogEntryRequest) (models.LogEntryResponse, error)
	GetById(ctx context.Context, id int) (models.LogEntryResponse, error)
	GetByAppName(ctx context.Context, appName string) ([]models.LogEntryResponse, error)
	GetByUser(ctx context.Context, user string) ([]models.LogEntryResponse, error)
	GetByLogLevel(ctx context.Context, logLevel string) ([]models.LogEntryResponse, error)
}

type handlers struct {
	logServicer LogServicer
	logger      logger.ILogger
}

func NewHandlers(logServicer LogServicer, logger logger.ILogger) handlers {
	return handlers{
		logServicer: logServicer,
		logger:      logger,
	}
}

func (h handlers) createHandler(w http.ResponseWriter, r *http.Request) models.StandardResponse {
	fmt.Println("My Context:", r.Header.Get("RequestId"))
	var request models.LogEntryRequest
	json.NewDecoder(r.Body).Decode(&request)
	if request == (models.LogEntryRequest{}) {
		return models.StandardResponse{Data: nil, Err: errors.New("request is invalid")}
	} else {
		r, err := h.logServicer.Create(r.Context(), request)
		return models.StandardResponse{Data: r, Err: err}
	}
}

func (h handlers) getByIdHandler(w http.ResponseWriter, r *http.Request) models.StandardResponse {
	params := mux.Vars(r)
	idParam := params["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return models.StandardResponse{Data: nil, Err: errors.New("id is not a valid integer")}
	} else {
		r, err := h.logServicer.GetById(r.Context(), id)
		return models.StandardResponse{Data: r, Err: err}
	}
}

func (h handlers) getByAppNameHandler(w http.ResponseWriter, r *http.Request) models.StandardResponse {
	params := mux.Vars(r)
	appName := params["app"]

	if appName == "" {
		return models.StandardResponse{Data: nil, Err: errors.New("Appname is not valid")}
	} else {
		r, err := h.logServicer.GetByAppName(r.Context(), appName)
		return models.StandardResponse{Data: r, Err: err}
	}
}

func (h handlers) getByUserHandler(w http.ResponseWriter, r *http.Request) models.StandardResponse {
	params := mux.Vars(r)
	user := params["user"]

	if user == "" {
		return models.StandardResponse{Data: nil, Err: errors.New("User is not valid")}
	} else {
		r, err := h.logServicer.GetByUser(r.Context(), user)
		return models.StandardResponse{Data: r, Err: err}
	}
}

func (h handlers) getByLogLevelHandler(w http.ResponseWriter, r *http.Request) models.StandardResponse {
	params := mux.Vars(r)
	loglevel := params["level"]

	if loglevel == "" {
		return models.StandardResponse{Data: nil, Err: errors.New("LogLevel is not valid")}
	} else {
		r, err := h.logServicer.GetByLogLevel(r.Context(), loglevel)
		return models.StandardResponse{Data: r, Err: err}
	}
}
