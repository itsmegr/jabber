package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type Client struct {
	Name string
	Group *Group
	Conn *websocket.Conn
	// channel to receive msg to client, broacasted from group
	Send chan []byte
}
/*
	readPump pumps messages from the websocket connection to the group.
	The application runs readPump in a per-connection goroutine. The application
	ensures that there is at most one reader on a connection by executing all
	reads from this goroutine.
*/
func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
		c.Group.Unregister <- c
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("Client Disconected: Closing connection", err)
			// }
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		message = []byte(fmt.Sprintf("%v sent : %v", c.Name, string(message)))
		c.Group.Broadcast <- message
	}
}

/*
	writePump pumps messages from the group to the websocket connection.
	A goroutine running writePump is started for each connection. The
	application ensures that there is at most one writer to a connection by
	executing all writes from this goroutine.
*/
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		c.Group.Unregister <- c
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The group closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			/*
				Add queued chat messages to the current websocket message.
				for better performance
				this will be automatically handled in another iterations
			*/
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
/*
	PING and PONG messages are described in the RFC. In summary, peers (including the browser) 
	automatically respond to a PING message with a PONG message.
	The best practice for detecting a dead client is to read with a deadline. 
	If the client application does not send messages frequently enough for the deadline you want, 
	then send PING messages to induce the client to send a PONG. Update the deadline in the pong handler and after reading a message.
	The chat, command and filewatch examples in this repository show how to use PING/PONG.
	The code you posted has a data race on lastResponse.
*/