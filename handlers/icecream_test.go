package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/sudarshan-reddy/benjerry/models"
)

type fakeIceCreamStore struct {
	models.IceCreamStore
	iceCream        *models.IceCream
	serializedStore string
	err             error
}

func (i *fakeIceCreamStore) StoreContext(ctx context.Context, iceCreamInput models.IceCream) error {
	bdy, err := json.Marshal(iceCreamInput)
	if err != nil {
		return err
	}
	i.serializedStore += string(bdy)
	return i.err
}

func (i *fakeIceCreamStore) Get(name string) (*models.IceCream, error) {
	return i.iceCream, i.err
}

func (i *fakeIceCreamStore) Update(iceCreamInput models.IceCream) error {
	bdy, err := json.Marshal(iceCreamInput)
	if err != nil {
		return err
	}
	i.serializedStore += string(bdy)
	return i.err
}

func (i *fakeIceCreamStore) Delete(name string) error {
	i.serializedStore += name
	return nil
}

func Test_PostIceCreamData(t *testing.T) {
	var tests = []struct {
		desc               string
		reqBody            io.Reader
		expectedResponse   string
		dbError            error
		expectedStatusCode int
		expectedDataStore  string
	}{
		{
			desc: "successful create should return a 201",
			reqBody: bytes.NewReader([]byte(`
				{
				    "name": "Chocobar",
				    "story": "cheap and best",
				    "description": "Some new stuff",
				    "sourcing_values": [
				        "Responsibly Sourced Packaging",
				        "Caring Dairy"
				    ],
				    "ingredients": [
				        "cream",
				        "skim milk",
				        "cocoa (processed with alkali)",
				        "natural flavors",
				        "cocoa",
				        "guar gum",
				        "butteroil",
				        "milk protein concentrate",
				        "corn starch",
				        "salt",
				        "soy lecithin",
				        "tapioca starch",
				        "pectin",
				        "caramelized sugar syrup",
				        "baking soda",
				        "molasses",
				        "honey",
				        "carrageenan",
				        "vanilla extract"
				    ],
				    "allergy_info": "contains milk, eggs, wheat and soy",
				    "dietary_certification": "",
				    "product_id": "1111"
				}
			`)),
			expectedResponse:   "\"\"\n",
			expectedStatusCode: 201,
			expectedDataStore: "{\"name\":\"Chocobar\",\"image_open\":\"\"," +
				"\"image_closed\":\"\",\"story\":\"cheap and best\"," +
				"\"description\":\"Some new stuff\",\"sourcing_values\":" +
				"[\"Responsibly Sourced Packaging\",\"Caring Dairy\"]," +
				"\"ingredients\":" +
				"[\"cream\",\"skim milk\",\"cocoa (processed with alkali)\"," +
				"\"natural flavors\",\"cocoa\",\"guar gum\",\"butteroil\"," +
				"\"milk protein concentrate\",\"corn starch\",\"salt\"," +
				"\"soy lecithin\",\"tapioca starch\",\"pectin\"," +
				"\"caramelized sugar syrup\",\"baking soda\",\"molasses\"," +
				"\"honey\",\"carrageenan\",\"vanilla extract\"]," +
				"\"allergy_info\":\"contains milk, eggs, wheat and soy\"," +
				"\"dietary_certification\":\"\",\"product_id\":\"1111\"}",
		},
		{
			desc:               "if database returns error, throw unexpected error response",
			reqBody:            bytes.NewReader([]byte("{}")),
			dbError:            errors.New("pg error : error in db"),
			expectedStatusCode: 500,
			expectedResponse: "{\"httpStatus\":500," +
				"\"httpCode\":\"internal_server_error\"," +
				"\"requestId\":\"\",\"errors\":[]}\n",
			expectedDataStore: "{\"name\":\"\",\"image_open\":\"\"," +
				"\"image_closed\":\"\",\"story\":\"\",\"description\":\"\"," +
				"\"sourcing_values\":null,\"ingredients\":null," +
				"\"allergy_info\":\"\",\"dietary_certification\":\"\"," +
				"\"product_id\":\"\"}",
		},
		{
			desc:               "invalid format returns a 400 error",
			reqBody:            bytes.NewReader([]byte("")),
			expectedStatusCode: 400,
			expectedResponse: "{\"httpStatus\":400,\"httpCode\":\"bad_request\"," +
				"\"requestId\":\"\",\"errors\":[{\"code\":\"format_error\"," +
				"\"message\":\"invalid input format. error: EOF\"}]}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert := assert.New(t)
			iceCreamStore := &fakeIceCreamStore{err: test.dbError}
			ich := NewIceCreamHandler(iceCreamStore)

			req, err := http.NewRequest("POST", "/url", test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ich.PostIceCreamData)
			handler.ServeHTTP(rr, req)
			assert.Equal(test.expectedResponse, rr.Body.String())
			assert.Equal(test.expectedStatusCode, rr.Code)
			assert.Equal(test.expectedDataStore, iceCreamStore.serializedStore)
		})
	}
}

func Test_GetIceCreamData(t *testing.T) {
	var tests = []struct {
		desc               string
		dbResponse         *models.IceCream
		dbError            error
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			desc: "if the database returns a successful response, show that in response to get",
			dbResponse: &models.IceCream{
				Name:        "chocobar",
				ImageOpen:   "imageOpen",
				ImageClosed: "imageClose",
				Story:       "some story",
				Description: "some description",
			},
			expectedResponse: "{\"name\":\"chocobar\",\"image_open\":\"imageOpen\"," +
				"\"image_closed\":\"imageClose\",\"story\":\"some story\"," +
				"\"description\":\"some description\",\"sourcing_values\":null," +
				"\"ingredients\":null,\"allergy_info\":\"\"," +
				"\"dietary_certification\":\"\",\"product_id\":\"\"}\n",
			expectedStatusCode: 200,
		},
		{
			desc:    "if database returns ErrNoRows, return 404 in response",
			dbError: models.ErrNoRows,
			expectedResponse: "{\"httpStatus\":404,\"httpCode\":\"not_found\"," +
				"\"requestId\":\"\",\"errors\":[{\"code\":\"not_found\"," +
				"\"message\":\"Icecream: chocobar Not Found\"}]}\n",
			expectedStatusCode: 404,
		},
		{
			desc:               "if database returns a different error, return 500 in response",
			dbError:            errors.New("pg error: error connecting to db"),
			expectedStatusCode: 500,
			expectedResponse: "{\"httpStatus\":500," +
				"\"httpCode\":\"internal_server_error\"," +
				"\"requestId\":\"\",\"errors\":[]}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert := assert.New(t)
			iceCreamStore := &fakeIceCreamStore{iceCream: test.dbResponse, err: test.dbError}
			ich := NewIceCreamHandler(iceCreamStore)

			req, err := http.NewRequest("GET", "/url", nil)
			if err != nil {
				t.Fatal(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ice-cream-name", "chocobar")
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ich.GetIceCreamData)
			handler.ServeHTTP(rr, req.WithContext(ctx))
			assert.Equal(test.expectedResponse, rr.Body.String())
			assert.Equal(test.expectedStatusCode, rr.Code)
		})
	}

}
