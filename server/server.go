package server

import (
	"jabber/server/handlers"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)


type Server struct {
	Router *chi.Mux
}

/*
	Router for join
	it handles all the path with "/join" path
*/
func joinRouter() http.Handler {
	Rt := chi.NewRouter();
	Rt.HandleFunc("/", handlers.JoinHandler)
	return Rt
}

// set all the routers for server
func SetRouters(router *chi.Mux){
	router.Mount("/join", joinRouter())
	router.Get("/stats", handlers.StatsHandler)
}

// return new main/global router
func NewRouter() *chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	SetRouters(r)
	return r
}

// return new server
func NewServer() *Server {
	return &Server{
		Router: NewRouter(),
	}
}



