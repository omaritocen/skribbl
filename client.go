package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {

	// unique id for the client
	id string

	// Name of the client
	name string

	// Reference of the hub in the client
	hub *Hub

	// Websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte
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

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {

}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {

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

	go client.readPump()
	go client.writePump()
}
