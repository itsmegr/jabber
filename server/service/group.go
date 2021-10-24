package service

// Group maintains the set of active clients and broadcasts messages to the
type Group struct {
	Name string
	Hub *Hub
	Clients map[*Client]bool
	Broadcast chan []byte
	Register chan *Client
	Unregister chan *Client
}

func newGroup(name string) *Group {
	return &Group{
		Name: name,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Group) run() {
	defer func ()  {
		//delete the group from groupList in Hub
		h.Hub.Unregister <- h
	}()
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			//check here if there is any client or not, if not then delete the group
			if len(h.Clients)==0 {
				//delete the group from groupList in Hub
				h.Hub.Unregister <- h
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