package main

import (
	"encoding/json"
	"log"
)

const (
	JoinRoomAction    string = "join-room"
	LeaveRoomAction          = "leave-room"
	TextMessageAction        = "text-message"
)

type Message struct {
	sender *Client
	body   []byte
	target string
	action string
}

func (m *Message) encode() []byte {
	jsonMessage, err := json.Marshal(*m)
	if err != nil {
		log.Fatal("Failed to marshall json message")
	}
	return jsonMessage
}

func decodeMessage(jsonMessage []byte) (message Message) {
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Fatal(err)
		return
	}

	return message
}
