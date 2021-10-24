package main

import (
	"jabber/server"
	"jabber/server/service"
	"net/http"
)


func main() {
	service.StartGlobalHub()
	server := server.NewServer()
	http.ListenAndServe(":8080", server.Router)
}