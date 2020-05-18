package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// InternalErrorCode

	// UnknownInternalErrorCode represents an error caused by an unexpected state.
	UnknownInternalErrorCode = iota + 1

	// DBErrorCode

	// UnknownDBErrorCode represents an error caused by the db being in an unexpected state.
	UnknownDBErrorCode
	// RecordNotFoundDBErrorCode represents a record not found error
	RecordNotFoundDBErrorCode

	// BadRequestErrorCode

	// InvalidTimeFormatBadRequestErrorCode represents a time parameter being in invalid format.
	InvalidTimeFormatBadRequestErrorCode
	// InvalidTimeTooEarlyBadRequestErrorCode represents a time parameter being to early in the current context.
	InvalidTimeTooEarlyBadRequestErrorCode
	// InvalidTimeTooLateBadRequestErrorCode represents a time parameter being to late in the current context.
	InvalidTimeTooLateBadRequestErrorCode

	// CryptoErrorCode

	// UnknownCryptoErrorCode represents an error caused by the crypto computation resulting in an unexpected state.
	UnknownCryptoErrorCode
)

// ErrorResponse represents an error response from the api
type ErrorResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause,omitempty"`
}

// Error represents an api error, the json marshall is meant to be sent to client
type Error struct {
	error
	HTTPStatusCode int
	ErrorCode      int
	ClientMessage  string
	Cause          error
}

func (e *Error) Error() string {
	return e.Cause.Error()
}

// MarshalJSON custom json marshal meant to be sent to client
func (e *Error) MarshalJSON() ([]byte, error) {
	cause := e.Error()
	if e.ErrorCode == UnknownInternalErrorCode && !gin.IsDebugging() && gin.Mode() != gin.TestMode {
		cause = ""
	}

	return json.Marshal(&ErrorResponse{
		ErrorCode: e.ErrorCode,
		Message:   e.ClientMessage,
		Cause:     cause,
	})
}

// NewDBError returns a DB error
func NewDBError(httpStatusCode int, code int, cause error, clientMessage string) *Error {
	return &Error{
		HTTPStatusCode: httpStatusCode,
		ClientMessage:  "DB Error: " + clientMessage,
		ErrorCode:      code,
		Cause:          cause,
	}
}

// NewRecordNotFoundDBError returns a DB error when a record is not found (with record information)
func NewRecordNotFoundDBError(cause error, recordInfo string) *Error {
	return NewDBError(http.StatusNotFound, RecordNotFoundDBErrorCode, cause, "Could not find the specified record "+recordInfo)
}

// NewBadRequestError returns a Bad Request error with bad parameter info
func NewBadRequestError(code int, cause error, badParameterInfo string) *Error {
	return &Error{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorCode:      code,
		ClientMessage:  "Bad request Invalid parameter: " + badParameterInfo,
		Cause:          cause,
	}
}

// NewUnknownDBError returns an unknown DB error with default message
func NewUnknownDBError(cause error) *Error {
	return NewUnknownInternalError(cause, "Database")
}

// NewUnknownDataFeedError returns an unknown Datafeed error with default message
func NewUnknownDataFeedError(cause error) *Error {
	return NewUnknownInternalError(cause, "Datafeed")
}

// NewUnknownCryptoServiceError returns an unknown CryptoService error with default message
func NewUnknownCryptoServiceError(cause error) *Error {
	return NewUnknownInternalError(cause, "CryptoService")
}

// NewUnknownInternalError returns a default Internal Error
func NewUnknownInternalError(cause error, relatedCause string) *Error {
	return &Error{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorCode:      UnknownInternalErrorCode,
		ClientMessage:  "Internal server error: Unexpected error occurred related to " + relatedCause,
		Cause:          cause,
	}
}
