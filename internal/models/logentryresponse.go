package models

import (
	"time"
)

type LogEntryResponse struct {
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
