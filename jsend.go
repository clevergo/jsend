package jsend

import (
	"encoding/json"
	"net/http"
)

// Status constants
const (
	StatusError   = "error"
	StatusFail    = "fail"
	StatusSuccess = "success"
)

// Body contains
type Body struct {
	// The status indicates the execution result of request,
	// it can be one of "success", "fail" and "error".
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// Error writes error body with the given message.
func Error(w http.ResponseWriter, message string, statuses ...int) error {
	body := Body{
		Status:  StatusError,
		Message: message,
	}
	return Write(w, body, statuses...)
}

// ErrorCode writes error body with the given message and code.
func ErrorCode(w http.ResponseWriter, message string, code int, statuses ...int) error {
	body := Body{
		Status:  StatusError,
		Message: message,
		Code:    code,
	}
	return Write(w, body, statuses...)
}

// ErrorCodeData writes error body with the given message, code and data.
func ErrorCodeData(w http.ResponseWriter, message string, code int, data interface{}, statuses ...int) error {
	body := Body{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Data:    data,
	}
	return Write(w, body, statuses...)
}

// Fail writes failed body with the given data.
func Fail(w http.ResponseWriter, data interface{}, statuses ...int) error {
	body := Body{
		Status: StatusFail,
		Data:   data,
	}
	return Write(w, body, statuses...)
}

// Success writes successful body with the given data.
func Success(w http.ResponseWriter, data interface{}, statuses ...int) error {
	body := Body{
		Status: StatusSuccess,
		Data:   data,
	}
	return Write(w, body, statuses...)
}

// Write writes the body to http.ResponseWriter.
//
// If necessary, the status code can be specified through the third parameter.
func Write(w http.ResponseWriter, body Body, statuses ...int) error {
	w.Header().Set("Content-Type", "application/json")

	if len(statuses) > 0 {
		w.WriteHeader(statuses[0])
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
