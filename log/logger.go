package log

import (
	"time"
)

type Logger interface {
	Log(Entry)
}

type Entry struct {
	LogLevel    Level
	Message     string
	ThreadName  string
	LoggerName  string
	Version     string
	Component   string
	ExtraFields map[string]interface{}
}

type AccessEntry struct {
	Entry

	RequestID    string
	Time         time.Time
	Duration     time.Duration
	RequestURI   string
	RemoteAddr   string
	Status       int
	Proto        string
	Method       string
	RequestBody  string
	ResponseBody string
}

type SqlEntry struct {
	Entry

	FileLineNum string
	Text        string
	Sql         string
	Duration    int64
	Error       interface{}
}
