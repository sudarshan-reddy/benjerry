package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/sudarshan-reddy/benjerry/httputils"
)

//AuthHandler dictates the interface that can be used to inject
//authentication
type AuthHandler interface {
	Authenticate(r *http.Request) (*http.Request, *httputils.HandlerError)
}

//authenticator is the primary structure holding all the data required
//to process scope authentication
type authenticator struct {
	authHandlers []AuthHandler
}

//Authenticator provides the implementation functions
//needed for scope authentication
type Authenticator interface {
	Authenticate(h http.Handler) http.Handler
}

//NewAuthenticator initiates a new instance of Authenticator
func NewAuthenticator(authHandlers ...AuthHandler) Authenticator {
	return &authenticator{
		authHandlers: authHandlers,
	}
}

//Authenticate implements an authentication scheme that
//can be strung together with registerScope or used on it's own with
//global registration
func (s *authenticator) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handlerError *httputils.HandlerError
		var authHandlerReq *http.Request
		//the idea here is to try to get the first positive result.
		//the authfuncs are run in order of assignment during
		//creation if the previous func yields no result (yields error)
		for _, authHandler := range s.authHandlers {
			authHandlerReq, handlerError = authHandler.Authenticate(r)
			if handlerError == nil {
				break
			}
		}

		if handlerError != nil {
			httputils.WriteHandlerError(handlerError, r, w)
			return
		}

		h.ServeHTTP(w, authHandlerReq)
	})
}

type statictokenauthenticator struct {
	//tokenScopes is provided as a mapping of allowed scopes by
	//static tokens.
	tokenScopes map[string][]string
}

//NewStaticTokenAuthenticator instantiates a new instance of authenticator.ScopeGetter
func NewStaticTokenAuthenticator(tokenScopes map[string][]string) AuthHandler {
	return &statictokenauthenticator{
		tokenScopes: tokenScopes,
	}
}

func (s *statictokenauthenticator) Authenticate(r *http.Request) (*http.Request, *httputils.HandlerError) {
	authToken := r.Header.Get("Authorization")

	if !strings.HasPrefix(authToken, "Bearer ") {
		subError := httputils.NewSubError(httputils.InvalidScope, "message", "Authorization Type 'Bearer ' is missing")
		return nil, httputils.NewHandlerError(http.StatusUnauthorized, subError)
	}

	authToken = strings.Replace(authToken, "Bearer ", "", 1)

	if scopes, ok := s.tokenScopes[authToken]; ok {
		scopeContext := context.WithValue(r.Context(), ContextKeyScopes, scopes)
		authContext := context.WithValue(scopeContext, ContextKeyAuthToken, authToken)
		r = r.WithContext(authContext)
		return r, nil
	}

	return nil, httputils.NewHandlerError(http.StatusUnauthorized, nil)
}
