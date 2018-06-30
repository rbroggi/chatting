package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//represets a single chat user
type client struct {
	//the websocket for this client
	socket *websocket.Conn
	//channel on which messages are sent
	send chan *message
	//the room this client is chatting on
	room *room
	//userData holds info about the user
	userData map[string]interface{}
}

//reads messages from the socket and publish them to the room
func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error while parsing JSON data incoming from socket: %v", err)
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- msg
	}
}

//write client messages to the public room
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			log.Fatal("Error while serializing msg data to socket", err)
			return
		}
	}
}
