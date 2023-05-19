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
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// TODO: Check if we really need the ping-pong handlers
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleNewMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	// TODO: Check if we really need the ticker
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send

		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}

	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	message := decodeMessage(jsonMessage)
	message.Sender = c.name

	switch message.Action {
	case JoinRoomAction:
		c.handleJoinRoomMessage(message)
	case LeaveRoomAction:
		c.handleLeaveRoomMessage(message)
	case TextMessageAction:
		c.handleTextMessage(message)

	}
}

func (c *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Body
	room := c.hub.findRoomByName(roomName)

	if room == nil {
		room = c.hub.createRoom(roomName)
	}

	// We need a way to send room id for the user

	room.register <- c
}

func (c *Client) handleLeaveRoomMessage(message Message) {
	roomName := message.Body
	room := c.hub.findRoomByName(roomName)

	if room == nil {
		return
	}

	room.unregister <- c
}

func (c *Client) handleTextMessage(message Message) {
	roomId := message.Target
	room := c.hub.findRoomById(roomId)

	if room == nil {
		return
	}

	room.broadcast <- &message
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

	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "NewUser"
	}

	client := NewClient(hub, conn, clientName)
	log.Printf("New user joined the system, Name: [%s], ID: [%s]\n", client.name, client.id)

	hub.register <- client

	go client.readPump()
	go client.writePump()
}
