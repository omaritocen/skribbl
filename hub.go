package main

import "log"

type Hub struct {
	// Registered Clients
	clients map[*Client]bool

	// Register Rooms
	rooms map[string]*Room

	// Register/Unregister requests from client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		}
	}

}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) createRoom(roomName string) *Room {
	room := NewRoom(roomName)
	h.rooms[room.id] = room

	go room.Run()

	log.Printf("Created new room, Name: [%s], id: [%s]\n", room.name, room.id)

	return room
}

func (h *Hub) findRoomById(id string) *Room {
	return h.rooms[id]
}

func (h *Hub) findRoomByName(name string) *Room {
	for _, room := range h.rooms {
		if room.name == name {
			return room
		}
	}

	return nil
}
