package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	hub *Hub

	// Websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Name of the client
	name string

	// unique id for the client
	id string
}

func NewClient(hub *Hub, conn *websocket.Conn, name string) *Client {
	return &Client{
		id:   uuid.NewString(),
		name: name,
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientName := r.URL.Query().Get("Name")
	if clientName == "" {
		clientName = "NewUser"
	}

	client := NewClient(hub, conn, clientName)
	log.Printf("New user joined the system, Name: [%s], ID: [%s]\n", client.name, client.id)

	hub.register <- client

	// TODO: Handle Client Read/Write
}
