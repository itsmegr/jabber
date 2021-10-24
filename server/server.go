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


func joinRouter() http.Handler {
	Rt := chi.NewRouter();
	Rt.HandleFunc("/", handlers.JoinHandler)
	return Rt
}


func SetRouters(router *chi.Mux){
	router.Mount("/join", joinRouter())
}


func NewRouter() *chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	SetRouters(r)
	return r
}




func NewServer() *Server {

	return &Server{
		Router: NewRouter(),
	}
}



