package main

import (
	"jabber/internal"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)




func joinHandler(resW http.ResponseWriter, reqR *http.Request){
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//starting new goroutine for GlobalHub, always up
	GlobalHub := internal.NewHub()
	go GlobalHub.Run()

	r.HandleFunc("/join/", joinHandler)
	http.ListenAndServe(":8080", r)
}