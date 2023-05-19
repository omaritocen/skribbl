package main

import (
	"encoding/json"
	"log"
)

const (
	/*
			JoinRoomMessage
		{
			Author: 0x0031f (pointer to client),
			Action: "join-room",
			Body: "room-name-1",
			Target: nil
		}
	*/
	JoinRoomAction string = "join-room"

	/*
			LeaveRoomMessage
		{
			Author: 0x0031f (pointer to client),
			Action: "leave-room",
			Body: "room-name-1",
			Target: nil
		}
	*/
	LeaveRoomAction = "leave-room"

	/*
			TextMessage
		{
			Author: 0x0031f (pointer to client),
			Action: "text-message",
			Body: "Hello this is text body",
			Target: "room-id-1"
		}
	*/
	TextMessageAction = "text-message"
)

type Message struct {
	Sender string `json:"sender"`
	Body   string `json:"body"`
	Target string `json:"target"`
	Action string `json:"action"`
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
