package dto

import (
	"time"
)

type LogEntry struct {
	Id        int
	AppName   string
	RequestId string
	User      string
	LogLevel  string
	Frames    string
	Message   string
	Params    string
	Details   string
	CreatedAt time.Time
}
