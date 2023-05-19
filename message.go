package main

import (
	"encoding/json"
	"log"
)

const (
	/*
			JoinRoomMessage: used to join/create a room
		{
			Sender: 0x0031f (pointer to client),
			Action: "join-room",
			Body: "room-name-1",
		}
	*/
	JoinRoomAction string = "join-room"

	/*
			LeaveRoomMessage: used to leave a room
		{
			Sender: "client-name-1",
			Action: "leave-room",
			Body: "room-name-1",
		}
	*/
	LeaveRoomAction = "leave-room"

	/*
			TextMessage: used to send messages in chat
		{
			Sender: "client-name-1",
			Action: "text-message",
			Body: "Hello this is text body",
		}
	*/
	TextMessageAction = "text-message"

	/*
			JoinRoomMessage: used to confirm user joining the room and send room info
		{
			Sender: "client-name-1",
			Action: "join-room",
			Body: "room-name-1",
		}
	*/
	UserJoinedRoomAction string = "user-joined-room"
)

type Message struct {
	Sender string `json:"sender"`
	Body   string `json:"body"`
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
