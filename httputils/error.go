package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//ErrorCodeStrings defines some possible error codes
var ErrorCodeStrings = map[ErrorCode]string{
	FormatError:      "format_error",
	NotFound:         "not_found",
	BadRequest:       "bad_request",
	InvalidScope:     "invalid_scope",
	UnexpectedError:  "unexpected_server_error",
	NotImplemented:   "not_implemented",
	InvalidOperation: "invalid_operation",
	InvalidParameter: "invalid_parameter",
	Deprecated:       "deprecated",
}

//ErrorCode int typecast for enum below
type ErrorCode int

const (
	_ ErrorCode = iota
	Custom
	NotFound
	FormatError
	BadRequest
	InvalidScope
	UnexpectedError
	NotImplemented
	InvalidOperation
	InvalidParameter
	Deprecated
)

//ErrorDetails is useful to parse error details
type ErrorDetails map[string]interface{}

//SubError is part of the http error structure
type SubError struct {
	Code    ErrorCode
	Details ErrorDetails
}

//MarshalJSON marshals error jsons
func (subError *SubError) MarshalJSON() ([]byte, error) {
	if subError.Code != Custom {
		subError.Details["code"] = ErrorCodeStrings[subError.Code]
	}
	return json.Marshal(subError.Details)
}

func (subError *SubError) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Code: %s", ErrorCodeStrings[subError.Code])
	for key, value := range subError.Details {
		fmt.Fprintf(&buf, "\n%s: %s", key, value)
	}
	return buf.String()
}

//NewSubError generates a new instance of SubError
func NewSubError(code ErrorCode, key string, value interface{}) *SubError {
	return &SubError{
		Code:    code,
		Details: ErrorDetails{key: value},
	}
}

//HandlerError defines HandlerError sructure
type HandlerError struct {
	HTTPStatusCode int
	SubErrors      []*SubError
}

func (handlerErr *HandlerError) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTPStatusCode: %d", handlerErr.HTTPStatusCode)
	for _, subError := range handlerErr.SubErrors {
		fmt.Fprintf(&buf, "\n%s", subError)
	}
	return buf.String()
}

//NewHandlerError returns a new instance of HandlerError
func NewHandlerError(statusCode int, subErrors ...*SubError) *HandlerError {
	return &HandlerError{
		HTTPStatusCode: statusCode,
		SubErrors:      subErrors,
	}
}

//NewNotFoundError ...
func NewNotFoundError(message string) *HandlerError {
	subError := NewSubError(NotFound, "message", message)
	return NewHandlerError(http.StatusNotFound, subError)
}

//NewFormatError ...
func NewFormatError(message string) *HandlerError {
	subError := NewSubError(FormatError, "message", message)
	return NewHandlerError(http.StatusBadRequest, subError)
}

//NewInvalidOperation ...
func NewInvalidOperation(message string) *HandlerError {
	subError := NewSubError(InvalidOperation, "message", message)
	return NewHandlerError(http.StatusConflict, subError)
}

//NewInvalidParameterError ...
func NewInvalidParameterError(message string) *HandlerError {
	subError := NewSubError(InvalidParameter, "message", message)
	return NewHandlerError(http.StatusBadRequest, subError)
}

//NewUnexpectedError ...
func NewUnexpectedError(err error) *HandlerError {
	subError := NewSubError(UnexpectedError, "error", err.Error())
	return NewHandlerError(http.StatusInternalServerError, subError)
}

//NewDeprecatedError ...
func NewDeprecatedError(message string) *HandlerError {
	subError := NewSubError(Deprecated, "message", message)
	return NewHandlerError(http.StatusGone, subError)
}

//NewCustomError ...
func NewCustomError(httpStatus int, code, message string) *HandlerError {
	subError := NewSubError(Custom, "code", code)
	subError.Details["message"] = message
	return NewHandlerError(httpStatus, subError)
}

//HTTPError ...
type HTTPError struct {
	Status    int         `json:"httpStatus"`
	Code      string      `json:"httpCode"`
	RequestID string      `json:"requestId"`
	Errors    []*SubError `json:"errors"`
}
