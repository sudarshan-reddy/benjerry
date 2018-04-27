package router

import (
	"net/http"

	"github.com/sudarshan-reddy/benjerry/httputils"
)

//ScopesType is used to indicate tokens that are
//static/provided by application
type ScopesType int

//AuthTokenType is used to indicate authtoken in context
type AuthTokenType int

const (
	//ContextKeyScopes is a constant that can be used to pick up the
	//scopes injected into context
	ContextKeyScopes ScopesType = 0
	//ContextKeyAuthToken is a constant that can be used to pick up
	//the authToken injected into context
	ContextKeyAuthToken AuthTokenType = 1
)

//ScopeMiddleware is a type cast of the http.Handler signature for
//ease of use
type ScopeMiddleware func(http.Handler) http.Handler

//AnyScope checks if the method is authorized to any particular scope
//This is done by looking for the `Scope` type within the request context
//Any implementation that wants to use IsAuthorized will have to use
//the underlying `middleware.Scopes` type to set context value.
func AnyScope(scopes []string) ScopeMiddleware {
	var scopeMap = make(map[string]struct{}, 0)
	for _, scopes := range scopes {
		scopeMap[scopes] = struct{}{}
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestScopes, ok := r.Context().Value(ContextKeyScopes).([]string)
			if !ok {
				subError := httputils.NewSubError(httputils.InvalidScope, "message", "context does not have scope set")
				httputils.WriteHandlerError(httputils.NewHandlerError(http.StatusUnauthorized, subError), r, w)
				return
			}

			for _, requestScope := range requestScopes {
				if _, ok := scopeMap[requestScope]; ok {
					h.ServeHTTP(w, r)
					return
				}
			}
			httputils.WriteHandlerError(httputils.NewHandlerError(http.StatusUnauthorized, nil), r, w)
		})
	}
}

//AllScopes checks if the method is authorized to access all the scopes
func AllScopes(scopes []string) ScopeMiddleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestScopes, ok := r.Context().Value(ContextKeyScopes).([]string)
			if !ok {
				subError := httputils.NewSubError(httputils.InvalidScope, "message", "context does not have scope set")
				httputils.WriteHandlerError(httputils.NewHandlerError(http.StatusUnauthorized, subError), r, w)
				return
			}

			for _, scope := range scopes {
				var foundScope bool
				for _, requestScope := range requestScopes {
					if scope == requestScope {
						foundScope = true
						break
					}
				}
				if !foundScope {
					httputils.WriteHandlerError(httputils.NewHandlerError(http.StatusUnauthorized, nil), r, w)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}

}
