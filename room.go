package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	//forward is a channel that handles all incoming messages for
	//this room and send them to all the subscribed clients
	forward chan []byte

	//clients whishing to join the room - for sync purpose
	join chan *client

	//clients whishing to leave the room - for sync purpose
	leave chan *client

	//Set of all clients in the room
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {

	for {
		//select statement will only run one block at time ensuring concurrent access to the map
		select {
		//new joining room:
		case client := <-r.join:
			r.clients[client] = true
		//client leaving the room:
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		//incoming message
		case msg := <-r.forward:
			//forward message to all clients in the room
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

//Turning room into an http.Handler
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte),
		room:   r,
	}

	r.join <- client
	defer func() { r.leave <- client }()

	go client.write()
	client.read()
}
