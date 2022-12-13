package models

type LogEntryRequest struct {
	AppName   string
	RequestId string
	User      string
	LogLevel  string
	Frames    string
	Message   string
	Params    string
	Details   string
}
