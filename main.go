package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func NumberOfGoroutines(){
    ticker := time.NewTicker(10 * time.Second)
	defer func ()  {
		ticker.Stop()
	}()
    done := make(chan bool)
	for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                fmt.Println("Tick at", t, runtime.NumGoroutine())
            }
    }
}

func handleClient(conn *websocket.Conn){
	defer func ()  {
		conn.Close()
		fmt.Println("function returned from handleClient function")	
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msgStr := string(p);
		msgStr = fmt.Sprintf("Same msg is returned from server : %v", msgStr);
		if err := conn.WriteMessage(messageType, []byte(msgStr)); err != nil {
			log.Println(err)
			return
		}
	}
}

func groupHandler(resW http.ResponseWriter, reqR *http.Request){
	defer func ()  {
		fmt.Println("function returned from groupHandler")	
	}()
	upgrader.CheckOrigin= func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(resW, reqR, nil)
	if err !=nil {
		log.Println(err)
		return
	}
	go handleClient(conn);
}

func main() {
	go NumberOfGoroutines()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.HandleFunc("/group", groupHandler)
	http.ListenAndServe(":8080", r)
}