package service

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
		case newGroup := <-h.Unregister:
			if _, ok := h.Groups[newGroup]; ok {
				delete(h.Groups, newGroup)
			}
		}
	}
}