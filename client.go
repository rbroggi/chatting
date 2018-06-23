package main

import (
	"github.com/gorilla/websocket"
)

//represets a single chat user
type client struct {
	//the websocket for this client
	socket *websocket.Conn
	//channel on which messages are sent
	send chan []byte
	//the room this client is chatting on
	room *room
}

//reads messages from the socket and publish them to the room
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

//write client messages to the public room
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
