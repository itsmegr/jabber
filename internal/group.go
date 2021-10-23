

package internal

// Group maintains the set of active clients and broadcasts messages to the
type Group struct {
	hub *Hub
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newGroup() *Group {
	return &Group{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Group) run() {
	defer func ()  {
		//delete the group from groupList in Hub
		h.hub.Unregister <- h
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			//check here if there is any client or not, if not then delete the group
			if len(h.clients)==0 {
				//delete the group from groupList in Hub
				h.hub.Unregister <- h
				return
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}