package main

type Hub struct {
	// Registered Clients
	clients map[*Client]bool

	// Inbound messages from client
	broadcast chan []byte

	// Register/Unregister requests from client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {

}
