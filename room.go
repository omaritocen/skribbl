package main

import (
	"fmt"
	"github.com/google/uuid"
)

type Room struct {
	id         string
	name       string
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewRoom(name string) *Room {
	return &Room{
		id:         uuid.NewString(),
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func (r *Room) Run() {
	for {
		select {
		case c := <-r.register:
			r.registerClient(c)
		case c := <-r.unregister:
			r.unregisterClient(c)
		case message := <-r.broadcast:
			r.broadcastToClients(message.encode())
		}
	}
}

func (r *Room) registerClient(client *Client) {
	r.notifyNewClientJoin(client)
	r.clients[client] = true
}

func (r *Room) unregisterClient(client *Client) {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
	}

	// If no more members in room delete room
	if len(r.clients) == 0 {
		// TODO: Handle closing the room
		// close all channels and free maps

	} else {
		r.notifyNewClientLeave(client)
	}
}

func (r *Room) broadcastToClients(message []byte) {
	for client := range r.clients {
		client.send <- message
	}
}

func (r *Room) notifyNewClientJoin(client *Client) {
	message := &Message{
		Sender: client.name,
		Target: r.id,
		Action: TextMessageAction,
		Body:   fmt.Sprintf("User %s has joined the room!", client.name),
	}

	r.broadcastToClients(message.encode())
}

func (r *Room) notifyNewClientLeave(client *Client) {
	message := &Message{
		Sender: client.name,
		Target: r.id,
		Action: TextMessageAction,
		Body:   fmt.Sprintf("User %s has left the room!", client.name),
	}

	r.broadcastToClients(message.encode())
}
