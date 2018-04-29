package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudarshan-reddy/benjerry/httputils"
)

type fakeAuthHandler struct {
	AuthHandler
	err *httputils.HandlerError
}

func (f *fakeAuthHandler) Authenticate(r *http.Request) (*http.Request, *httputils.HandlerError) {
	if f.err != nil {
		return nil, f.err
	}
	return r, nil
}

func Test_Authenticate(t *testing.T) {

	var testCases = []struct {
		desc         string
		authHandlers []AuthHandler
		expectedResp string
	}{
		{
			"basic case",
			[]AuthHandler{&fakeAuthHandler{err: nil}},
			"",
		},
		{
			"when the first authentication is success dont call the second one",
			[]AuthHandler{
				&fakeAuthHandler{err: nil},
				&fakeAuthHandler{err: httputils.NewHandlerError(http.StatusUnauthorized, nil)},
			},
			"",
		},
		{
			"when the first authentication is a failure call the second one",
			[]AuthHandler{
				&fakeAuthHandler{err: httputils.NewHandlerError(http.StatusUnauthorized, nil)},
				&fakeAuthHandler{err: nil},
			},
			"",
		},
		{
			"when the all cases are a failure throw error",
			[]AuthHandler{
				&fakeAuthHandler{err: httputils.NewHandlerError(http.StatusUnauthorized, nil)},
				&fakeAuthHandler{err: httputils.NewHandlerError(http.StatusUnauthorized, nil)},
			},
			`{"httpStatus":401,"httpCode":"unauthorized","requestId":"","errors":[]}` + "\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			assert := assert.New(t)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			authenticator := NewAuthenticator(testCase.authHandlers...)
			authHandler := authenticator.Authenticate(handler)

			req, err := http.NewRequest("GET", "/url", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			authHandler.ServeHTTP(rr, req)

			assert.Equal(testCase.expectedResp, rr.Body.String())
		})
	}

}
