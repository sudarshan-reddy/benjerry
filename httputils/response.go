//Package httputils is a package that has
//helper functions to print to responseWriter
//these can be extrapolated into a separate util library as well
package httputils

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

//ContextRequestIDKey specifies the request id for scope context set
var ContextRequestIDKey interface{} = "requestId"

var httpStatusCodes = map[int]string{
	http.StatusInternalServerError: "internal_server_error",
	http.StatusConflict:            "conflict",
	http.StatusNotFound:            "not_found",
	http.StatusBadRequest:          "bad_request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
}

//AbbreAuthToken helps abbreviate the auth token to prevent showing
//all of it in the log messages
func AbbreAuthToken(authToken string) string {
	charsToReveal := 4
	if len(authToken) < charsToReveal {
		charsToReveal = len(authToken)
	}
	return authToken[:charsToReveal] + "..."
}

//GetAbbreAuthToken is a helper to get AbbreAuthToken
func GetAbbreAuthToken(r *http.Request) string {
	authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	return AbbreAuthToken(authToken)
}

//WriteHandlerError parses error to ResponseWriter and logs it
func WriteHandlerError(handlerErr *HandlerError, r *http.Request, w http.ResponseWriter) {
	requestID := ""
	if id, ok := r.Context().Value(ContextRequestIDKey).(string); ok {
		requestID = id
	}
	httpError := &HTTPError{
		Status:    handlerErr.HTTPStatusCode,
		RequestID: requestID,
		Errors:    handlerErr.SubErrors,
		Code:      httpStatusCodes[handlerErr.HTTPStatusCode],
	}

	logFields := map[string]interface{}{
		"error":      handlerErr,
		"requestURI": r.RequestURI,
		"method":     r.Method,
		"requestId":  requestID,
		"authToken":  GetAbbreAuthToken(r),
	}
	log.WithFields(logFields).Error("request failed")

	if handlerErr.HTTPStatusCode == http.StatusInternalServerError ||
		handlerErr.HTTPStatusCode == http.StatusForbidden ||
		handlerErr.HTTPStatusCode == http.StatusUnauthorized {
		httpError.Errors = make([]*SubError, 0)
	}

	err := WriteJSON(handlerErr.HTTPStatusCode, httpError, w)
	if err != nil {
		log.Fatal("serializing http error failed: ", err)
	}
}

//WriteJSON writes to w
func WriteJSON(status int, response interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}
