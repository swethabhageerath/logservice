package repositories

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/swethabhageerath/logging/lib/constants"
	"github.com/swethabhageerath/logging/lib/logger"
	"github.com/swethabhageerath/logservice/internal/dto"
	"github.com/swethabhageerath/logservice/internal/exceptions"
	"github.com/swethabhageerath/logservice/internal/models"
)

type DataAccessor interface {
	ExecuteContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type LogRepository struct {
	da     DataAccessor
	logger logger.ILogger
}

func NewLogRepository(da DataAccessor, logger logger.ILogger) LogRepository {
	return LogRepository{
		da:     da,
		logger: logger,
	}
}

func (l LogRepository) Create(ctx context.Context, le dto.LogEntry) (dto.LogEntry, error) {
	query := "INSERT INTO \"public\".\"Log\"(AppName, RequestId, User, LogLevel,Message, Frames, Params, Details) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	result, err := l.da.ExecuteContext(ctx, query, le.AppName, le.RequestId, le.User, le.LogLevel, le.Message, le.Frames, le.Params, le.Details)
	if err != nil {
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: "Create"}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return dto.LogEntry{}, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}

	lastRowInsertedId, err := result.LastInsertId()
	if err != nil {
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: "Create"}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return dto.LogEntry{}, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}

	if lastRowInsertedId == 0 {
		err := errors.New("An error occured while creating a log record")
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: "Create"}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return dto.LogEntry{}, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}

	return l.GetById(ctx, int(lastRowInsertedId))
}

func (l LogRepository) GetById(ctx context.Context, id int) (dto.LogEntry, error) {
	query := "SELECT Id, AppName, RequestId, User, LogLevel, Message, Frames, Params, Details WHERE \"ID\"=$1"
	row := l.da.QueryRowContext(ctx, query, id)
	var le dto.LogEntry
	err := row.Scan(&le.Id, &le.AppName, &le.RequestId, &le.User, &le.LogLevel, &le.Message, &le.Frames, &le.Params, &le.Details)
	if err != nil {
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: "GetById"}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return dto.LogEntry{}, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}
	return le, nil
}

func (l LogRepository) GetByAppName(ctx context.Context, appName string) ([]dto.LogEntry, error) {
	query := "SELECT Id, AppName, RequestId, User, LogLevel, Message, Frames, Params, Details WHERE \"AppName\"=$1"
	rows, err := l.da.QueryContext(ctx, query, appName)
	return l.get(rows, ctx, "GetByAppName", err)
}

func (l LogRepository) GetByUser(ctx context.Context, user string) ([]dto.LogEntry, error) {
	query := "SELECT Id, AppName, RequestId, User, LogLevel, Message, Frames, Params, Details WHERE \"User\"=$1"
	rows, err := l.da.QueryContext(ctx, query, user)
	return l.get(rows, ctx, "GetByUser", err)
}

func (l LogRepository) GetByLogLevel(ctx context.Context, logLevel string) ([]dto.LogEntry, error) {
	query := "SELECT Id, AppName, RequestId, User, LogLevel, Message, Frames, Params, Details WHERE \"LogLevel\"=$1"
	rows, err := l.da.QueryContext(ctx, query, logLevel)
	return l.get(rows, ctx, "GetByLogLevel", err)
}

func (l LogRepository) get(row *sql.Rows, ctx context.Context, operation string, err error) ([]dto.LogEntry, error) {
	if err != nil {
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: operation}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return nil, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}

	logEntries := make([]dto.LogEntry, 0)

	for row.Next() {
		var le dto.LogEntry
		err := row.Scan(&le.Id, &le.AppName, &le.RequestId, &le.User, &le.LogLevel, &le.Message, &le.Frames, &le.Params, &le.Details)
		if err != nil {
			l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: operation}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
			return nil, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
		}
		logEntries = append(logEntries, le)
	}

	err = row.Err()
	if err != nil {
		l.logger.AddAppName("logservice").AddContext(ctx).AddDetails(models.LogDetail{PackageName: "Repositories", FileName: "LogRepository.go", Operation: operation}.String()).AddFrames(constants.ALL).AddMessage(err.Error()).Log()
		return nil, exceptions.Exception{Err: errors.Wrap(err, err.Error()), StatusCode: 500}
	}

	return logEntries, nil
}
