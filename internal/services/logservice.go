package services

import (
	"context"

	"github.com/pkg/errors"
	"github.com/swethabhageerath/logging/lib/logger"
	"github.com/swethabhageerath/logservice/internal/dto"
	"github.com/swethabhageerath/logservice/internal/models"
)

type LogRepositor interface {
	Create(ctx context.Context, le dto.LogEntry) (dto.LogEntry, error)
	GetById(ctx context.Context, id int) (dto.LogEntry, error)
	GetByAppName(ctx context.Context, appName string) ([]dto.LogEntry, error)
	GetByUser(ctx context.Context, user string) ([]dto.LogEntry, error)
	GetByLogLevel(ctx context.Context, logLevel string) ([]dto.LogEntry, error)
}

type LogService struct {
	logger       logger.ILogger
	logRepositor LogRepositor
}

func NewLogService(logger logger.ILogger, logRepositor LogRepositor) LogService {
	return LogService{
		logger:       logger,
		logRepositor: logRepositor,
	}
}

func (l LogService) Create(ctx context.Context, le models.LogEntryRequest) (models.LogEntryResponse, error) {
	logEntry := dto.LogEntry{
		AppName:   le.AppName,
		RequestId: le.RequestId,
		User:      le.User,
		LogLevel:  le.LogLevel,
		Frames:    le.Frames,
		Message:   le.Message,
		Params:    le.Params,
		Details:   le.Details,
	}

	r, err := l.logRepositor.Create(ctx, logEntry)
	if err != nil {
		return models.LogEntryResponse{}, err
	}

	logEntryResponse := models.LogEntryResponse{
		Id:        r.Id,
		AppName:   r.AppName,
		RequestId: r.RequestId,
		User:      r.User,
		LogLevel:  r.LogLevel,
		Frames:    r.Frames,
		Message:   r.Message,
		Params:    r.Params,
		Details:   r.Details,
		CreatedAt: r.CreatedAt,
	}

	return logEntryResponse, nil
}

func (l LogService) GetById(ctx context.Context, id int) (models.LogEntryResponse, error) {
	r, err := l.logRepositor.GetById(ctx, id)
	if err != nil {
		return models.LogEntryResponse{}, errors.Wrap(err, "LogService.GetById()")
	}

	logEntryResponse := models.LogEntryResponse{
		Id:        r.Id,
		AppName:   r.AppName,
		RequestId: r.RequestId,
		User:      r.User,
		LogLevel:  r.LogLevel,
		Frames:    r.Frames,
		Message:   r.Message,
		Params:    r.Params,
		Details:   r.Details,
		CreatedAt: r.CreatedAt,
	}
	return logEntryResponse, nil
}

func (l LogService) GetByAppName(ctx context.Context, appName string) ([]models.LogEntryResponse, error) {
	r, err := l.logRepositor.GetByAppName(ctx, appName)
	return l.get(ctx, r, err)
}

func (l LogService) GetByUser(ctx context.Context, user string) ([]models.LogEntryResponse, error) {
	r, err := l.logRepositor.GetByUser(ctx, user)
	return l.get(ctx, r, err)
}

func (l LogService) GetByLogLevel(ctx context.Context, logLevel string) ([]models.LogEntryResponse, error) {
	r, err := l.logRepositor.GetByLogLevel(ctx, logLevel)
	return l.get(ctx, r, err)
}

func (l LogService) get(ctx context.Context, le []dto.LogEntry, err error) ([]models.LogEntryResponse, error) {
	if err != nil {
		return nil, err
	}

	logEntryResponses := make([]models.LogEntryResponse, len(le))

	for _, i := range le {
		logEntryResponse := models.LogEntryResponse{
			Id:        i.Id,
			AppName:   i.AppName,
			RequestId: i.RequestId,
			User:      i.User,
			LogLevel:  i.LogLevel,
			Frames:    i.Frames,
			Message:   i.Message,
			Params:    i.Params,
			Details:   i.Details,
			CreatedAt: i.CreatedAt,
		}

		logEntryResponses = append(logEntryResponses, logEntryResponse)
	}

	return logEntryResponses, nil
}
