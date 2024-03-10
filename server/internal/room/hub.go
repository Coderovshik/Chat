package room

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrRoomExist    = errors.New("room already exists")
	ErrRoomNotExist = errors.New("room does not exist")
)

type Hub struct {
	Rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) IsRoomExist(id string) bool {
	_, ok := h.Rooms[id]
	return ok
}

func (h *Hub) CreateRoom(req *CreateRoomRequest) (*CreateRoomResponse, error) {
	r := NewRoom(uuid.NewString(), req.Name)

	if _, ok := h.Rooms[r.ID]; !ok {
		go r.Run()
		h.Rooms[r.ID] = r
		return &CreateRoomResponse{
			Name: r.Name,
			ID:   r.ID,
		}, nil
	}

	return nil, ErrRoomExist
}

func (h *Hub) DeleteRoom(id string) error {
	_, ok := h.Rooms[id]
	if !ok {
		return ErrRoomNotExist
	}

	h.Rooms[id].Stop()
	delete(h.Rooms, id)

	return nil
}

func (h *Hub) JoinRoom(id string, client *Client) error {
	room := h.Rooms[id]

	room.Register <- client

	msg := &Message{
		Body:     fmt.Sprintf("user %s has joined the room", client.Username),
		RoomID:   id,
		Username: client.Username,
	}

	room.Broadcast <- msg

	go client.WriteMessage()
	go client.ReadAndBroadcastMessage(room)

	return nil
}

func (h *Hub) GetRooms() (*GetRoomsResponse, error) {
	res := &GetRoomsResponse{
		Rooms: make([]RoomData, 0, len(h.Rooms)),
	}

	for _, v := range h.Rooms {
		res.Rooms = append(res.Rooms, RoomData{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return res, nil
}

func (h *Hub) GetClients(id string) (*GetClientsResponse, error) {
	room, ok := h.Rooms[id]
	if !ok {
		return nil, ErrRoomNotExist
	}

	res := &GetClientsResponse{
		Clients: make([]ClientData, 0, len(room.Clients)),
	}

	for _, v := range room.Clients {
		res.Clients = append(res.Clients, ClientData{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return res, nil
}
