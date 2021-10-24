package service

import "log"

// Group maintains the set of active clients and broadcasts messages to the
type Group struct {
	Name string
	Hub *Hub
	Clients map[*Client]bool
	Broadcast chan []byte
	Register chan *Client
	Unregister chan *Client
}

func NewGroup(name string) *Group {
	return &Group{
		Name: name,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Group) Run() {
	defer func ()  {
		GlobalHub.Unregister <- h
	}()
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client: %v, Joined Group: %v", client.Name, h.Name)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client: %v, Disconnected From Group: %v", client.Name, h.Name)
			}
			//check here if there is any client or not, if not then delete the group
			if len(h.Clients)==0 {
				//delete the group from groupList in Hub
				GlobalHub.Unregister <- h
				return
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}