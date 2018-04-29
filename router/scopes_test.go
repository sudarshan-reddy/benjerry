package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AnyScope(t *testing.T) {

	var testCases = []struct {
		desc             string
		contextScopes    []string
		registeredScopes []string
		expectedResponse string
	}{
		{
			"if the request context has a valid scope, dont throw an error",
			[]string{"first", "second", "third"},
			[]string{"zeroth", "third", "fifth"},
			"",
		},
		{
			"if the request context does not has a valid scope, dont throw an error",
			[]string{"first", "second", "third"},
			[]string{"zeroth", "sixth", "fifth"},
			`{"httpStatus":401,"httpCode":"unauthorized","requestId":"","errors":[]}` + "\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			assert := assert.New(t)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			scopedHandlerFunc := AnyScope(testCase.registeredScopes)

			req, err := http.NewRequest("GET", "/url", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), ContextKeyScopes, testCase.contextScopes)

			rr := httptest.NewRecorder()

			scopedHandlerFunc(handler).ServeHTTP(rr, req.WithContext(ctx))

			assert.Equal(testCase.expectedResponse, rr.Body.String())
		})

	}
}

func Test_AllScopes(t *testing.T) {
	var testCases = []struct {
		desc             string
		contextScopes    []string
		registeredScopes []string
		expectedResponse string
	}{
		{
			"if the request context has all valid scopes, dont throw an error",
			[]string{"first", "second", "third"},
			[]string{"first", "second", "third"},
			"",
		},
		{
			"if the request context does not match with all valid scope throw an error",
			[]string{"first", "second", "third"},
			[]string{"first", "second", "fifth"},
			`{"httpStatus":401,"httpCode":"unauthorized","requestId":"","errors":[]}` + "\n",
		},
		{
			"if the request context has all valid scopes and they are jumbled, dont throw an error",
			[]string{"first", "second", "third"},
			[]string{"first", "third", "second"},
			"",
		},
		{
			"if the request context has all less scopes and they are jumbled, throw an error",
			[]string{"first", "second"},
			[]string{"first", "third", "second"},
			`{"httpStatus":401,"httpCode":"unauthorized","requestId":"","errors":[]}` + "\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			assert := assert.New(t)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			scopedHandlerFunc := AllScopes(testCase.registeredScopes)

			req, err := http.NewRequest("GET", "/url", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), ContextKeyScopes, testCase.contextScopes)

			rr := httptest.NewRecorder()

			scopedHandlerFunc(handler).ServeHTTP(rr, req.WithContext(ctx))

			assert.Equal(testCase.expectedResponse, rr.Body.String())
		})

	}

}
