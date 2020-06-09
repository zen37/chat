package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {

	//forward is channel that holds incoming messges that should be forwarded to the other clients
	forward chan []byte

	//channel for clients wishing to join the room
	join chan *client

	//channel for clients wishing to leave the room
	leave chan *client

	//holds all current clients in this room
	clients map[*client]struct{}
}

func (r *room) run() {
	log.Println("I am in run function")

	var empty struct{}

	for {
		select {
		//if a message is received in any of the room channels
		//select will run the code for the respective channel
		case client := <-r.join:
			r.clients[client] = empty
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("I am in ServeHTTP")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	log.Println("client: ", client)

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

// newRoom makes a new room.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]struct{}),
	}
}
