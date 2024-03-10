package room

import (
	"fmt"
)

type Room struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Done       chan struct{}
}

func NewRoom(id, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (r *Room) Run() {
	for {
		select {
		case c := <-r.Register:
			if _, ok := r.Clients[c.ID]; !ok {
				r.Clients[c.ID] = c
			}
		case c := <-r.Unregister:
			delete(r.Clients, c.ID)
			close(c.Message)
			c.Conn.Close()

			if len(r.Clients) != 0 {
				msg := &Message{
					Body:     fmt.Sprintf("user %s left the chat", c.Username),
					RoomID:   c.RoomID,
					Username: c.Username,
				}
				r.Broadcast <- msg
			}

			// TODO: stop gorutine on last client exit
		case m := <-r.Broadcast:
			for _, c := range r.Clients {
				c.Message <- m
			}
		case <-r.Done:
			return
		}
	}
}

func (r *Room) Stop() {
	close(r.Done)
}

type RoomService interface {
	IsRoomExist(id string) bool
	GetRooms() (*GetRoomsResponse, error)
	CreateRoom(req *CreateRoomRequest) (*CreateRoomResponse, error)
	DeleteRoom(id string) error
	JoinRoom(id string, client *Client) error
	GetClients(id string) (*GetClientsResponse, error)
}

type Message struct {
	Body     string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type CreateRoomResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type GetRoomsResponse struct {
	Rooms []RoomData
}

type RoomData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetClientsResponse struct {
	Clients []ClientData
}

type ClientData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type DeleteRoomResponse struct {
	ID string `json:"id"`
}
