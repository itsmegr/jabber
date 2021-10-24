package handlers

import (
	"jabber/server/service"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func findGroup(name string) (*service.Group, bool) {
	allGroups := service.GlobalHub.Groups
	for group, ok := range allGroups {
		if(group.Name==name&&ok){
			return group, true
		}
	}
	return nil, false
}

func manageGroup(name string) *service.Group {
	//finding the group
	group, isThere := findGroup(name)

	//no group
	if !isThere {
		group = service.NewGroup(name)
		service.GlobalHub.Register <- group
		go group.Run()
	}
	return group
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func JoinHandler(w http.ResponseWriter, r *http.Request){
	queryParams := r.URL.Query();
	groupName := queryParams.Get("group")
	clientName := queryParams.Get("name")


	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &service.Client{
		Name : clientName,
		Group: manageGroup(groupName),
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Group.Register <- client
	go client.WritePump()
	go client.ReadPump()
	// // fmt.Println(groupName, clientName)
	// w.Write([]byte(fmt.Sprintf("%v, %v", groupName, clientName)))
}





