package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// LogEntry Log
type LogEntry struct {
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	Body       string            `json:"body,omitempty"`
	Response   string            `json:"response,omitempty"`
	HttpStatus int               `json:"http_status,omitempty"`
	Header     map[string]string `json:"header,omitempty"`

	Start int64  `json:"start,omitempty"`
	End   int64  `json:"end,omitempty"`
	Cost  int64  `json:"cost,omitempty"`
	Error string `json:"error,omitempty"`
}

// NewLogEntryFromRequest ...
//
//	@param request
//	@return *RequestLog
func NewLogEntryFromRequest(request *http.Request) *LogEntry {
	if request == nil {
		return &LogEntry{}
	}
	entry := &LogEntry{
		Method: request.Method,
		URL:    request.URL.String(),
	}
	header := map[string]string{}
	for k, v := range request.Header {
		header[k] = strings.Join(v, "")
	}
	if request.Body != nil {

	}
	entry.Header = header
	return entry
}

// Start
//
//	@receiver entry
func (entry *LogEntry) SetStart() {
	entry.Start = time.Now().UnixMilli()
}

// SetEnd EndWithRespone SetEnd
//
//	@receiver entry
func (entry *LogEntry) SetEnd() {
	entry.End = time.Now().UnixMilli()
	entry.Cost = entry.End - entry.Start
}

func (entry *LogEntry) AddRespone(bytes []byte, httpStatus int64) {
	entry.Response = string(bytes)
	entry.HttpStatus = http.DefaultMaxHeaderBytes
}

// EndWithError
//
//	@receiver entry
//	@param err
func (entry *LogEntry) AddError(err string) {
	entry.Error = err
}

func (entry *LogEntry) JSON() ([]byte, error) {
	return json.Marshal(entry)
}

// Map
//
//	@receiver entry
//	@return map
func (entry *LogEntry) Map() map[string]interface{} {
	return map[string]interface{}{
		"method":   entry.Method,
		"url":      entry.URL,
		"body":     entry.Body,
		"start":    entry.Start,
		"end":      entry.End,
		"cost":     entry.Cost,
		"response": entry.Response,
		"error":    entry.Error,
	}
}
