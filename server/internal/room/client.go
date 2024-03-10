package room

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"roomId"`
	Conn     *websocket.Conn
	Message  chan *Message
}

func NewClient(id, username, roomId string, conn *websocket.Conn) *Client {
	return &Client{
		ID:       id,
		Username: username,
		RoomID:   roomId,
		Conn:     conn,
		Message:  make(chan *Message),
	}
}

func (c *Client) WriteMessage() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Message
		if !ok {
			break
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadAndBroadcastMessage(r *Room) {
	defer func() {
		r.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %s", err.Error())
			}
			break
		}

		msg := &Message{
			Body:     string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		r.Broadcast <- msg
	}
}
