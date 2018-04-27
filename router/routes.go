package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sudarshan-reddy/benjerry/handlers"
	"github.com/sudarshan-reddy/benjerry/models"
)

const (
	apiVersion1 = "/api/v1"
)

//Router holds all the api based routes
type Router struct {
	*chi.Mux
	authenticator Authenticator
	Config
}

//Config holds the config values required for router to work
type Config struct {
	IceCreamStore models.IceCreamStore
}

//NewRouter returns a new instance of Router
func NewRouter(staticTokens map[string][]string, cfg Config) *Router {
	staticTokenAuth := NewStaticTokenAuthenticator(staticTokens)
	return &Router{
		authenticator: NewAuthenticator(staticTokenAuth),
		Mux:           chi.NewRouter(),
		Config:        cfg,
	}
}

//AddRoutes adds all the routes to the router
//Scoping and middleware should also be done here
func (router *Router) AddRoutes() {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	iceCreamHandler := handlers.NewIceCreamHandler(router.Config.IceCreamStore)

	router.Group(func(r chi.Router) {
		r.Use(router.authenticator.Authenticate)

		r.With(AnyScope([]string{"*", "post.icecream"})).
			Post(apiVersion1+"/create", iceCreamHandler.PostIceCreamData)

		r.With(AnyScope([]string{"*", "read.icecream"})).
			Get(apiVersion1+"/read/{ice-cream-name}", iceCreamHandler.GetIceCreamData)

		r.With(AnyScope([]string{"*", "post.icecream"})).
			Put(apiVersion1+"/update", iceCreamHandler.UpdateIceCreamData)

		r.With(AnyScope([]string{"*", "delete.icecream"})).
			Delete(apiVersion1+"/delete/{ice-cream-name}", iceCreamHandler.DeleteIceCreamData)
	})
}
