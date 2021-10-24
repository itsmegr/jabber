package service

import "log"


/*
	One Goroutine for Hub running all the time
	GlobalHub represents complete application
	Hub contains all the groups
*/
type Hub struct {
	Groups map[*Group]bool
	Register chan *Group
	Unregister chan *Group
}
var GlobalHub *Hub

func newHub() *Hub{
	return &Hub{
		Groups: make(map[*Group]bool),
		Register: make(chan *Group),
		Unregister: make(chan *Group),
	}
}


func StartGlobalHub(){
	GlobalHub = newHub()
	go GlobalHub.Run()
}

func (h *Hub) Run(){
	for {
		select {
		case newGroup := <-h.Register:
			h.Groups[newGroup] = true
			log.Printf("Group : %v, Registered!!!", newGroup.Name)
		case newGroup := <-h.Unregister:
			if _, ok := h.Groups[newGroup]; ok {
				groupName := newGroup.Name
				delete(h.Groups, newGroup)
				log.Printf("No Client in Group : %v, hence Unregistered!!!", groupName)
			}
		}
	}
}